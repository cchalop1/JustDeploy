import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";

import { CreateServiceApi } from "@/services/createServiceApi";
import CommandModal from "./CommandModal";

type AddServiceProps = {
  deployId?: string;
  projectId?: string;
  setLoading: (loading: boolean) => void;
  createService: (serviceParams: CreateServiceApi) => Promise<void>;
  fetchServiceList: (deployId?: string) => Promise<void>;
};

export default function AddService({
  projectId,
  setLoading,
  createService,
  fetchServiceList,
}: AddServiceProps) {
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const text =
    "Click here to connect a github repos or create a new service. You can also press";
  const [openCommandModal, setOpenCommandModal] = useState(false);

  async function getServices() {
    const res = await getPreConfiguredServiceListApi(projectId);
    setPreConfiguredServices(res);
  }
  useEffect(() => {
    getServices();
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpenCommandModal((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <div
        onClick={() => setOpenCommandModal(true)}
        className="hover:shadow-md hover:bg-slate-100 cursor-pointer pt-3 pb-6 pl-5 pr-5 flex w-80 h-36 rounded-lg border border-dashed border-gray-600 justify-center items-center"
      >
        <div>
          {text}{" "}
          <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
            <span className="text-xs">⌘</span>K
          </kbd>
        </div>
      </div>
      <CommandModal
        open={openCommandModal}
        setOpen={setOpenCommandModal}
        preConfiguredServices={preConfiguredServices}
        create={async (createServiceParams) => {
          setLoading(true);
          setOpenCommandModal(false);
          await createService({
            path: createServiceParams.path,
            serviceName: createServiceParams.serviceName,
          });
          setLoading(false);
          await getServices();
        }}
      />
    </>
  );
}
