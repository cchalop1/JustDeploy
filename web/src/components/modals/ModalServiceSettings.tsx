import { motion } from "framer-motion";

import Modal from "@/components/modals/Modal";
import { Service } from "@/services/getServicesByDeployId";
import ServiceDeploySettings from "../ServiceDeploySettings";
import { CardIcon } from "../CardIcon";
import { Button } from "../ui/button";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { useNotification } from "@/hooks/useNotifications";
import { Badge } from "../ui/badge";

type ModalServiceSettingsProps = {
  service: Service;
  onClose: () => void;
  fetchServices: () => Promise<void>;
};
export default function ModalServiceSettings({
  service,
  onClose,
  fetchServices,
}: ModalServiceSettingsProps) {
  const notif = useNotification();

  async function deleteServiceById() {
    const res = await deleteServiceByIdApi(service.id);
    if (res) {
      onClose();
      notif.success({
        title: "Service is deleted",
        content: `${service.name} is deleted`,
      });
      await fetchServices();
    }
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
            <CardIcon service={service} />
            <div className="font-bold">{service.name}</div>
          </div>
        }
      >
        <div className="p-3 border-t mt-2">
          <div className="flex flex-col justify-between gap-2">
            <div>
              <Badge variant="outline">{service.status}</Badge>
            </div>
            <ServiceDeploySettings
              service={service}
              fetchService={fetchServices}
            />
          </div>
        </div>
        <div className="flex justify-end mt-4">
          <Button variant="destructive" onClick={deleteServiceById}>
            Delete
          </Button>
        </div>
      </Modal>
    </motion.div>
  );
}
