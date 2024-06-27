import { Service } from "@/services/getServicesByDeployId";
import { useEffect, useState } from "react";
import { Card } from "../ui/card";
import { Button } from "../ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import ModalDeleteConfirmation from "../modals/ModalDeleteConfirmation";
import { deleteServiceApi } from "@/services/deleteServiceApi";

type ServiceListDeployProps = {
  deployId: string;
  services: Service[];
  loadingNewService: boolean;
  fetchServiceList: () => Promise<void>;
};

export default function ServiceListDeploy({
  deployId,
  services,
  loadingNewService,
  fetchServiceList,
}: ServiceListDeployProps) {
  const [serviceToDelete, setServiceToDelete] = useState<Service | null>(null);

  useEffect(() => {
    fetchServiceList();
  }, [deployId]);

  async function deleteService() {
    const serviceId = serviceToDelete?.id;

    if (!serviceId) {
      return;
    }
    try {
      await deleteServiceApi(deployId, serviceId);
      await fetchServiceList();
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
          <Card key={s.name} className="flex justify-between p-3 mb-3 h-24">
            <div className="flex flex-col gap-3">
              <div className="flex items-center gap-5">
                <div className="text-xl font-bold">{s.name}</div>
                <div>{s.status}</div>
              </div>
              <div>{s.imageName}</div>
            </div>
            <div>
              <Button
                variant="destructive"
                onClick={() => setServiceToDelete(s)}
              >
                Delete
              </Button>
            </div>
          </Card>
        ))}
        {loadingNewService && (
          <Card className="flex justify-between p-3 h-24">
            <div className="flex flex-col gap-3 w-full">
              <div className="flex justify-between">
                <Skeleton className="w-44 h-6" />
                <Skeleton className="w-24 h-10" />
              </div>
              <Skeleton className="w-20 h-3" />
            </div>
          </Card>
        )}
      </div>
    </>
  );
}
