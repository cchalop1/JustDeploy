import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";
import { Button } from "./ui/button";

type ModalDnsSettings = {
  domain: string;
  open: boolean;
  onOpenChange: (bool: boolean) => void;
  onClick: (event: React.FormEvent) => void;
};

export default function ModalDnsSettings({
  domain,
  open,
  onClick,
  onOpenChange,
}: ModalDnsSettings) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Server correcly setup ?</DialogTitle>
          <DialogDescription>
            <div>
              Before connect your server make sure your dns is correcly setup.
            </div>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Type</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Value</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow>
                  <TableCell>A</TableCell>
                  <TableCell>{domain}</TableCell>
                  <TableCell>{"{your server ip}"}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button onClick={onClick}>Connect and setup server</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
