import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { Service } from "@/services/getServicesByDeployId";
import { Button } from "../ui/button";
import { motion } from "framer-motion";
import EnvsManagements from "../forms/EnvsManagements";
import { useNotification } from "@/hooks/useNotifications";
import Modal from "./Modal";

type ModalServiceSettingsProps = {
  service: Service;
  onClose: () => void;
  getProjectById: () => Promise<void>;
};
export default function ModalServiceSettings({
  service,
  onClose,
  getProjectById,
}: ModalServiceSettingsProps) {
  const notif = useNotification();
  const isDevContainer = service.isDevContainer;
  const url = `http://localhost:${service.exposePort}`;

  const deleteSelectedService = async () => {
    await deleteServiceByIdApi(service.id);
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
            <div className="font-bold">{service.name}</div>
          </div>
        }
      >
        <div className="bg-gray-100 p-3 border-t">
          <Button variant="destructive" onClick={() => deleteSelectedService()}>
            Delete
          </Button>
          <Button onClick={copyEnv}>Copy Env</Button>
          <div>
            <div>
              <a href={url} target="_blank" className="underline">
                {url}
              </a>
            </div>
            <div>
              <div className="font-bold">Environment variables</div>
              <EnvsManagements envs={service.envs} setEnvs={() => {}} />
            </div>
          </div>
        </div>
      </Modal>
    </motion.div>
  );
}
