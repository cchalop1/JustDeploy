import SpinnerIcon from "@/assets/SpinnerIcon";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogContent,
  DialogDescription,
} from "@/components/ui/dialog";
import { getApplicationLogs } from "@/services/getApplicationLogs";
import { useEffect, useState } from "react";

type ModalApplicationLogs = {
  open: boolean;
  onOpenChange: (bool: boolean) => void;
  appName: string;
};

export default function ModalApplicationLogs({
  open,
  onOpenChange,
  appName,
}: ModalApplicationLogs) {
  const [logs, setLogs] = useState<null | string[]>(null);

  useEffect(() => {
    getApplicationLogs(appName).then(setLogs);
  }, []);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Logs of {appName} </DialogTitle>
        </DialogHeader>
        <DialogDescription className="flex flex-col">
          {logs === null ? (
            <SpinnerIcon color="text-black" />
          ) : (
            logs.map((log) => <code className="w-full mb-2">{log}</code>)
          )}
        </DialogDescription>
      </DialogContent>
    </Dialog>
  );
}
