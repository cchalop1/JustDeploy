import { motion } from "framer-motion";
import { Suspense, useState } from "react";

import Modal from "@/components/modals/Modal";
import { Service } from "@/services/getServicesByDeployId";
import ServiceDeploySettings from "../ServiceDeploySettings";
import { CardIcon } from "../CardIcon";
import { Button } from "../ui/button";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { useNotification } from "@/hooks/useNotifications";
import { Badge } from "../ui/badge";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ServiceLogs from "../ServiceLogs";
import ServiceCommitInfo from "../ServiceCommitInfo";

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
  const [isDeleting, setIsDeleting] = useState(false);

  async function deleteServiceById() {
    setIsDeleting(true);
    try {
      const res = await deleteServiceByIdApi(service.id);
      if (res) {
        onClose();
        notif.success({
          title: "Service is deleted",
          content: `${service.name} is deleted`,
        });
        await fetchServices();
      }
    } catch (error) {
      notif.error({
        title: "Error deleting service",
        content: "An error occurred while deleting the service.",
      });
    } finally {
      setIsDeleting(false);
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
        width: "600px",
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
        className="w-[600px] max-w-full"
      >
        <div className="p-3 border-t mt-2 w-full max-w-full overflow-hidden">
          <div className="flex flex-col justify-between gap-2 w-full">
            <div>
              <Badge variant="outline">{service.status}</Badge>
            </div>
            <Tabs defaultValue="settings" className="w-full">
              <TabsList className="w-full">
                <TabsTrigger value="settings" className="flex-1">
                  Settings
                </TabsTrigger>
                <TabsTrigger value="logs" className="flex-1">
                  Logs
                </TabsTrigger>
              </TabsList>
              <TabsContent value="settings">
                <div className="space-y-4">
                  <ServiceCommitInfo serviceId={service.id} />
                  <ServiceDeploySettings
                    service={service}
                    fetchServices={fetchServices}
                  />
                </div>
              </TabsContent>
              <TabsContent value="logs">
                <Suspense fallback={<div>Loading...</div>}>
                  <ServiceLogs serviceId={service.id} />
                </Suspense>
              </TabsContent>
            </Tabs>
          </div>
        </div>
        <div className="flex justify-end mt-4">
          <Button
            variant="destructive"
            onClick={deleteServiceById}
            disabled={isDeleting}
          >
            {isDeleting ? (
              <>
                <SpinnerIcon color="text-white" />
                Deleting...
              </>
            ) : (
              "Delete"
            )}
          </Button>
        </div>
      </Modal>
    </motion.div>
  );
}
