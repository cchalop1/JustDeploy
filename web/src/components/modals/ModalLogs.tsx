import { use } from "react";
import { Dialog, DialogContent } from "../ui/dialog";
import DisplayLogs from "../DisplayLogs";
import { Logs } from "@/services/getApplicationLogs";

type ModalLogsProps = {
  open: boolean;
  onOpenChange: (bool: boolean) => void;
  fetchServerProxyLogs: () => Promise<Array<Logs>>;
};

export default function ModalLogs({
  open,
  onOpenChange,
  fetchServerProxyLogs,
}: ModalLogsProps) {
  const logs = use(fetchServerProxyLogs());
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="h-64 overflow-scroll">
          <DisplayLogs logs={logs} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
