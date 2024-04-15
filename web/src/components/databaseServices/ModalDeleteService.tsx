import { deleteServiceApi } from "@/services/deleteServiceApi";
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
import { useState } from "react";
import { ButtonStateEnum } from "@/lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";

type DeleteServiceAlertProps = {
  deployId: string;
  serviceToDelete: Service | null;
  fetchServiceList: () => Promise<void>;
  setServiceToDelete: (service: Service | null) => void;
};

export default function ModalDeleteService({
  deployId,
  serviceToDelete,
  fetchServiceList,
  setServiceToDelete,
}: DeleteServiceAlertProps) {
  const [btnIsLoading, setBtnIsLoading] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  async function deleteService() {
    const serviceId = serviceToDelete?.id;

    if (!serviceId) {
      return;
    }
    setBtnIsLoading(ButtonStateEnum.PENDING);
    try {
      await deleteServiceApi(deployId, serviceId);
      await fetchServiceList();
      setBtnIsLoading(ButtonStateEnum.SUCESS);
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
          This action is irreversible. Deleting this service will remove all
          data associated with it.
        </AlertDialogDescription>
        <AlertDialogFooter>
          <Button onClick={() => setServiceToDelete(null)}>Cancel</Button>
          <Button variant="destructive" onClick={deleteService}>
            {btnIsLoading === ButtonStateEnum.PENDING ? (
              <SpinnerIcon color="text-white" />
            ) : (
              "Delete"
            )}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
