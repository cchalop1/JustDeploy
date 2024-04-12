import { Service } from "@/services/getServicesByDeployId";
import { useEffect, useState } from "react";
import { Card } from "./ui/card";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
} from "./ui/dialog";
import { deleteServiceApi } from "@/services/deleteServiceApi";
import { DialogTitle } from "@radix-ui/react-dialog";
import { Skeleton } from "@/components/ui/skeleton";

type ServiceListDeployProps = {
  deployId: string;
  services: Service[];
  loadingNewService: boolean;
  fetchServiceList: (deployId: string) => Promise<void>;
};

export default function ServiceListDeploy({
  deployId,
  services,
  loadingNewService,
  fetchServiceList,
}: ServiceListDeployProps) {
  const [serviceToDelete, setServiceToDelete] = useState<Service | null>(null);

  async function deleteService() {
    if (serviceToDelete === null) return;
    try {
      await deleteServiceApi(deployId, serviceToDelete.id);
      await fetchServiceList(deployId);
    } catch (e) {
      console.error(e);
    }
    setServiceToDelete(null);
  }

  useEffect(() => {
    fetchServiceList(deployId);
  }, [deployId]);

  return (
    <div className="flex flex-col">
      <Dialog
        open={serviceToDelete !== null}
        onOpenChange={(open) => {
          if (open) {
            return;
          } else {
            setServiceToDelete(null);
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              Are you sure you want to delete this service?
            </DialogTitle>
          </DialogHeader>
          <DialogDescription>
            <div>
              This action is irreversible. Deleting this service will remove all
              data associated with it.
            </div>
          </DialogDescription>
          <DialogFooter>
            <Button onClick={() => setServiceToDelete(null)}>Cancel</Button>
            <Button variant="destructive" onClick={deleteService}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      {services.map((s) => (
        <Card className="flex justify-between p-3 mb-3 h-24">
          <div className="flex flex-col gap-3">
            <div className="flex items-center gap-5">
              <div className="text-xl font-bold">{s.name}</div>
              <div>{s.status}</div>
            </div>
            <div>{s.imageName}</div>
          </div>
          <div>
            <Button variant="destructive" onClick={() => setServiceToDelete(s)}>
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
  );
}
