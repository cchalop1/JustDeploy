import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "../ui/button";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { useState } from "react";

type ModalConfirmationProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
  onConfirm: () => Promise<void>;
  title: string;
  content: string;
  textConfirm: string;
};

export default function ModalDeleteConfirmation({
  open,
  setOpen,
  onConfirm,
  title,
  content,
  textConfirm,
}: ModalConfirmationProps) {
  const [isLoading, setIsLoading] = useState(false);
  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogDescription>{content}</AlertDialogDescription>
        <AlertDialogFooter>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button
            variant="destructive"
            onClick={async () => {
              setIsLoading(true);
              await onConfirm();
              setIsLoading(false);
            }}
          >
            {isLoading ? <SpinnerIcon color="text-white" /> : textConfirm}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
