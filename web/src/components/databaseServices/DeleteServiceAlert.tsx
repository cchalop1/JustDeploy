import { deleteServiceApi } from "@/services/deleteServiceApi";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@radix-ui/react-dialog";
import { DialogFooter, DialogHeader } from "../ui/dialog";
import { Button } from "../ui/button";
import { Service } from "@/services/getServicesByDeployId";

type DeleteServiceAlertProps = {
  deployId: string;
  serviceToDelete: Service | null;
  fetchServiceList: () => Promise<void>;
  setServiceToDelete: (service: Service | null) => void;
};

export default function DeleteServiceAlert({
  deployId,
  serviceToDelete,
  fetchServiceList,
  setServiceToDelete,
}: DeleteServiceAlertProps) {
  async function deleteService() {
    const serviceId = serviceToDelete?.id;

    if (!serviceId) {
      return;
    }
    try {
      await deleteServiceApi(deployId, serviceId);
      await fetchServiceList();
    } catch (e) {
      console.error(e);
    }
    setServiceToDelete(null);
  }

  return (
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
  );
}
