import { Copy, SquareArrowOutUpRight, Trash } from "lucide-react";
import EnvsManagements from "./forms/EnvsManagements";
import { Button } from "./ui/button";
import { useNotification } from "@/hooks/useNotifications";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { copyToClipboard } from "@/lib/utils";
import { Service } from "@/services/getServicesByDeployId";

type ServiceLocalSettingsProps = {
  service: Service;
  onClose: () => void;
  getProjectById: () => Promise<void>;
};

export default function ServiceLocalSettings({
  service,
  onClose,
  getProjectById,
}: ServiceLocalSettingsProps) {
  const notif = useNotification();
  const url = `http://localhost:${service.exposePort}`;

  function copyEnv() {
    const env = service.envs
      .map((e: { name: string; value: string }) => `${e.name}=${e.value}`)
      .join("\n");
    copyToClipboard(env, ({ title }) =>
      notif.success({
        title: title,
        content: "Environment variables copied to clipboard",
      })
    );
  }

  const deleteSelectedService = async () => {
    await deleteServiceByIdApi(service.id);
    onClose();
    await getProjectById();
  };

  return (
    <>
      <div className="flex justify-between w-[25vw]">
        <div className="font-bold">Expose URL: </div>
        <a href={url} target="_blank" className="underline flex items-center">
          {url}
          <SquareArrowOutUpRight className="h-4" />
        </a>
      </div>
      <div>
        <div className="flex justify-between items-center mb-2">
          <div className="font-bold">Environment variables</div>
          <Button variant="outline" onClick={copyEnv}>
            <Copy className="h-4" />
            Copy env
          </Button>
        </div>
        <EnvsManagements envs={service.envs} setEnvs={() => {}} />
      </div>
      <Button
        className="mt-2 mb-2 w-full"
        variant="destructive"
        onClick={() => deleteSelectedService()}
      >
        <Trash className="h-4 font-bold" />
        Delete
      </Button>
    </>
  );
}
