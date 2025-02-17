import { Service } from "@/services/getServicesByDeployId";
import { X } from "lucide-react";
import { Button } from "../ui/button";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";

type ServiceSideBarProps = {
  serviceSelected: Service;
  setServiceSelected: (service: Service | null) => void;
  getProjectById: () => Promise<void>;
};

export default function ServiceSideBar({
  serviceSelected,
  setServiceSelected,
  getProjectById,
}: ServiceSideBarProps) {
  const deleteSelectedService = async () => {
    await deleteServiceByIdApi(serviceSelected.id);
    setServiceSelected(null);
    await getProjectById();
  };

  return (
    <div
      hidden={!serviceSelected}
      className="absolute right-5 top-28 w-1/4 rounded-lg border shadow-lg h-full bg-white p-8"
    >
      <div className="flex justify-between">
        <div className="font-bold text-2xl">{serviceSelected?.url}</div>
        <X
          className="h-6 cursor-pointer"
          onClick={() => setServiceSelected(null)}
        />
      </div>
      <Button variant="destructive" onClick={() => deleteSelectedService()}>
        Delete
      </Button>
    </div>
  );
}
