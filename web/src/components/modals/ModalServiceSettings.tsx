import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { Service } from "@/services/getServicesByDeployId";
import { Button } from "../ui/button";
import { motion } from "framer-motion";
import EnvsManagements from "../forms/EnvsManagements";
import { useNotification } from "@/hooks/useNotifications";
import Modal from "./Modal";
import { Copy, SquareArrowOutUpRight, Trash } from "lucide-react";

type ModalServiceSettingsProps = {
  projectId: string;
  service: Service;
  onClose: () => void;
  getProjectById: () => Promise<void>;
};
export default function ModalServiceSettings({
  projectId,
  service,
  onClose,
  getProjectById,
}: ModalServiceSettingsProps) {
  const notif = useNotification();
  const isDevContainer = service.isDevContainer;
  const url = `http://localhost:${service.exposePort}`;

  const deleteSelectedService = async () => {
    await deleteServiceByIdApi(projectId, service.id);
    onClose();
    await getProjectById();
  };

  function copyEnv() {
    const env = service.envs.map((e) => `${e.name}=${e.value}`).join("\n");
    navigator.clipboard.writeText(env);
    notif.success({
      title: "Copied",
      content: "Environment variables copied to clipboard",
    });
  }

  return (
    <motion.div
      initial={{
        opacity: 0,
        x: 100,
        height: "95%",
        position: "absolute",
        right: 5,
        top: 28,
        width: "35%",
        zIndex: 20,
      }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
    >
      <Modal
        onClose={onClose}
        headerNode={
          <div className="flex items-center gap-4">
            {isDevContainer ? (
              <img src="/icons/folder.png" className="w-8" />
            ) : (
              <img src="/icons/service.png" className="w-8" />
            )}
            <div className="font-bold">{service.hostName}</div>
          </div>
        }
      >
        <div className="p-3 border-t">
          <div className="flex flex-col gap-2">
            <div className="flex justify-between">
              <div className="font-bold">Expose URL: </div>
              <div className="">
                <a
                  href={url}
                  target="_blank"
                  className="underline flex items-center"
                >
                  {url}
                  <SquareArrowOutUpRight className="h-4" />
                </a>
              </div>
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
          </div>

          <Button
            className="mt-2 mb-2 w-full"
            variant="destructive"
            onClick={() => deleteSelectedService()}
          >
            <Trash className="h-4 font-bold" />
            Delete
          </Button>
        </div>
      </Modal>
    </motion.div>
  );
}
