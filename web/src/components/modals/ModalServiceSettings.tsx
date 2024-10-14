import { motion } from "framer-motion";

import Modal from "@/components/modals/Modal";
import ServiceLocalSettings from "@/components/ServiceLocalSettings";
import { Service } from "@/services/getServicesByDeployId";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ServiceDeploySettings from "../ServiceDeploySettings";
import { ProjectDto } from "@/services/getProjectById";
import { Suspense } from "react";
import SpinnerIcon from "@/assets/SpinnerIcon";

type ModalServiceSettingsProps = {
  project: ProjectDto;
  service: Service;
  onClose: () => void;
  getProjectById: () => Promise<void>;
};
export default function ModalServiceSettings({
  project,
  service,
  onClose,
  getProjectById,
}: ModalServiceSettingsProps) {
  const isDevContainer = service.isDevContainer;

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
        <Tabs defaultValue="local">
          <TabsList>
            <TabsTrigger value="local">Local Settings</TabsTrigger>
            {isDevContainer && (
              <TabsTrigger value="deploy">Deploy Settings</TabsTrigger>
            )}
          </TabsList>
          <div className="p-3 border-t mt-2">
            <TabsContent value="local">
              <div className="flex flex-col justify-between gap-2">
                <ServiceLocalSettings
                  project={project}
                  service={service}
                  getProjectById={getProjectById}
                  onClose={onClose}
                />
              </div>
            </TabsContent>
            <TabsContent value="deploy">
              <div className="flex flex-col justify-between gap-2">
                <ServiceDeploySettings
                  project={project}
                  service={service}
                  getProjectById={getProjectById}
                  onClose={onClose}
                />
              </div>
            </TabsContent>
          </div>
        </Tabs>
      </Modal>
    </motion.div>
  );
}
