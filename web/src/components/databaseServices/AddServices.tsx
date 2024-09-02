import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";

import { CreateServiceApi } from "@/services/createServiceApi";
import CommandModal from "./CommandModal";
import {
  ResponseServiceFromDockerComposeDto,
  getServicesFromDockerComposeApi,
} from "@/services/getServicesFromDockerCompose";

type AddServiceProps = {
  deployId?: string;
  setLoading: (loading: boolean) => void;
  createService: (serviceParams: CreateServiceApi) => Promise<void>;
  fetchServiceList: (deployId?: string) => Promise<void>;
};

export default function AddService({
  deployId,
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

  async function getServicesFromDockerCompose() {
    if (!deployId) return;
    const res = await getServicesFromDockerComposeApi(deployId);
    setServiceFromDockerCompose(res);
  }

  useEffect(() => {
    getServices();
    getServicesFromDockerCompose();
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
      <Button
        variant="outline"
        className="w-full mb-3"
        onClick={() => setOpen(true)}
      >
        {text}{" "}
        <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandModal
        open={open}
        setOpen={setOpen}
        preConfiguredServices={preConfiguredServices}
        serviceFromDockerCompose={serviceFromDockerCompose}
        createService={async (serviceName, fromDockerCompose) => {
          setLoading(true);
          setOpen(false);
          await createService({
            serviceName,
            fromDockerCompose,
            deployId,
          });
          await fetchServiceList(deployId);
          setLoading(false);
        }}
      />
    </>
  );
}
