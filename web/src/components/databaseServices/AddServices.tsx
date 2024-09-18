import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";

import { CreateServiceApi } from "@/services/createServiceApi";
import CommandModal from "./CommandModal";
import { ResponseServiceFromDockerComposeDto } from "@/services/getServicesFromDockerCompose";
import { Card } from "../ui/card";

type AddServiceProps = {
  deployId?: string;
  projectId?: string;
  setLoading: (loading: boolean) => void;
  createService: (serviceParams: CreateServiceApi) => Promise<void>;
  fetchServiceList: (deployId?: string) => Promise<void>;
};

export default function AddService({
  setLoading,
  createService,
  fetchServiceList,
}: AddServiceProps) {
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const text = "Click here for create and connect new database or press";
  const [serviceFromDockerCompose, setServiceFromDockerCompose] =
    useState<ResponseServiceFromDockerComposeDto>(null);

  const [open, setOpen] = useState(false);

  async function getServices() {
    const res = await getPreConfiguredServiceListApi();
    setPreConfiguredServices(res);
  }

  // async function getServicesFromDockerCompose() {
  //   if (!deployId) return;
  //   const res = await getServicesFromDockerComposeApi(deployId);
  //   setServiceFromDockerCompose(res);
  // }

  useEffect(() => {
    getServices();
    // getServicesFromDockerCompose();
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <Card
        onClick={() => setOpen(true)}
        className="hover:shadow-md hover:bg-slate-200 cursor-pointer pt-3 pb-6 pl-5 pr-5 flex gap-6 w-80 h-32 align-top"
      >
        {text}{" "}
        <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Card>
      <CommandModal
        open={open}
        setOpen={setOpen}
        preConfiguredServices={preConfiguredServices}
        serviceFromDockerCompose={serviceFromDockerCompose}
        create={async (createServiceParams) => {
          setLoading(true);
          setOpen(false);
          await createService({
            ...createServiceParams,
          });
          await fetchServiceList(deployId);
          setLoading(false);
        }}
      />
    </>
  );
}
