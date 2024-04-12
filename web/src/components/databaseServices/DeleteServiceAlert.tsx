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
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "../ui/alert-dialog";

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
    <AlertDialog
      open={serviceToDelete !== null}
      onOpenChange={(open) => {
        if (open) {
          return;
        } else {
          setServiceToDelete(null);
        }
      }}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Are you sure you want to delete this service?
          </AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogDescription>
          <div>
            This action is irreversible. Deleting this service will remove all
            data associated with it.
          </div>
        </AlertDialogDescription>
        <AlertDialogFooter>
          <Button onClick={() => setServiceToDelete(null)}>Cancel</Button>
          <Button variant="destructive" onClick={deleteService}>
            Delete
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
