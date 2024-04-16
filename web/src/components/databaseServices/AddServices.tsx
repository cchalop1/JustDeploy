import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";

import { createServiceApi } from "@/services/createServiceApi";
import CommandModal from "./CommandModal";
import {
  ResponseServiceFromDockerComposeDto,
  getServicesFromDockerComposeApi,
} from "@/services/getServicesFromDockerCompose";

type AddServiceProps = {
  deployId: string;
  setLoading: (loading: boolean) => void;
  fetchServiceList: () => Promise<void>;
};

export default function AddService({
  deployId,
  fetchServiceList,
  setLoading,
}: AddServiceProps) {
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const [serviceFromDockerCompose, setServiceFromDockerCompose] =
    useState<ResponseServiceFromDockerComposeDto>(null);

  const [open, setOpen] = useState(false);

  async function getServices() {
    const res = await getPreConfiguredServiceListApi();
    setPreConfiguredServices(res);
  }

  async function getServicesFromDockerCompose() {
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

  async function createService(
    serviceName: string,
    fromDockerCompose: boolean,
  ) {
    try {
      setLoading(true);
      setOpen(false);
      await createServiceApi(deployId, {
        serviceName,
        fromDockerCompose,
      });
      await fetchServiceList();
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  }

  return (
    <>
      <Button
        variant="outline"
        className="w-full mb-3"
        onClick={() => setOpen(true)}
      >
        Click here for create and connect new database or press{" "}
        <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandModal
        open={open}
        setOpen={setOpen}
        preConfiguredServices={preConfiguredServices}
        serviceFromDockerCompose={serviceFromDockerCompose}
        createService={createService}
      />
    </>
  );
}
