import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { Service } from "@/services/getServicesByDeployId";
import { CircleX } from "lucide-react";
import { Button } from "../ui/button";
import { motion } from "framer-motion";

type ModalServiceSettingsProps = {
  service: Service;
  setServiceSelected: (service: Service | null) => void;
  getProjectById: () => Promise<void>;
};
export default function ModalServiceSettings({
  service,
  setServiceSelected,
  getProjectById,
}: ModalServiceSettingsProps) {
  const isDevContainer = service.isDevContainer;

  const deleteSelectedService = async () => {
    await deleteServiceByIdApi(service.id);
    setServiceSelected(null);
    await getProjectById();
  };

  return (
    <motion.div
      initial={{
        opacity: 0,
        x: 100,
        height: "95%",
        position: "absolute",
        right: 5,
        top: 28,
        width: "100%",
      }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
    >
      <div className="absolute right-5 top-28 w-1/3 border border-spacing-3 rounded-lg h-3/4 bg-slate-100 shadow-lg">
        <div className="flex justify-between p-3 bg-white">
          <div className="flex items-center gap-4">
            {isDevContainer ? (
              <img src="/icons/folder.png" className="w-8" />
            ) : (
              <img src="/icons/service.png" className="w-8" />
            )}
            <div className="font-bold">{service.name}</div>
          </div>
          <div className="flex items-center">
            <div>tag</div>
            <CircleX
              className="w-8 cursor-pointer"
              onClick={() => setServiceSelected(null)}
            />
          </div>
        </div>
        <div className="bg-gray-100 p-3 border-t">
          <Button variant="destructive" onClick={() => deleteSelectedService()}>
            Delete
          </Button>
          <div></div>
        </div>
      </div>
    </motion.div>
  );
}
