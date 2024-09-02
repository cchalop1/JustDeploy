import { Service } from "@/services/getServicesByDeployId";
import { useEffect, useState } from "react";
import { Card } from "../ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import ModalDeleteConfirmation from "../modals/ModalDeleteConfirmation";
import { deleteServiceApi } from "@/services/deleteServiceApi";
import DatabaseServiceCard from "./DatabaseServiceCard";
import DatabaseServiceCardLoading from "./DatabaseServiceCardLoading";

type ServiceListDeployProps = {
  deployId?: string;
  services: Service[];
  loadingNewService: boolean;
  fetchServiceList: (deployId?: string) => Promise<void>;
};

export default function ServiceListDeploy({
  deployId,
  services,
  loadingNewService,
  fetchServiceList,
}: ServiceListDeployProps) {
  const [serviceToDelete, setServiceToDelete] = useState<Service | null>(null);

  useEffect(() => {
    fetchServiceList(deployId);
  }, [deployId]);

  async function deleteService() {
    const serviceId = serviceToDelete?.id;

    if (!serviceId) {
      return;
    }
    try {
      await deleteServiceApi(serviceId, deployId);
      await fetchServiceList(deployId);
      // TODO: send a toast to the user
    } catch (e) {
      console.error(e);
    }
    setServiceToDelete(null);
  }
  return (
    <>
      <ModalDeleteConfirmation
        open={serviceToDelete !== null}
        content="This action is irreversible. Deleting this service will remove all data associated with it."
        title="Are you sure you want to delete this service?"
        onConfirm={deleteService}
        setOpen={() => {}}
        textConfirm="Delete"
      />
      <div className="flex flex-col">
        {services.map((s) => (
          <DatabaseServiceCard
            key={s.id}
            service={s}
            setServiceToDelete={setServiceToDelete}
          />
        ))}
        {loadingNewService && <DatabaseServiceCardLoading />}
      </div>
    </>
  );
}
