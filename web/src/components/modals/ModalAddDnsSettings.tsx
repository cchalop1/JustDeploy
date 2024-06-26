import { useActionState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogFooter,
  DialogHeader,
} from "../ui/dialog";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { addDomainToServerApi } from "@/services/addDomainToServerApi";

type ModalAddDnsSettingsProps = {
  open: boolean;
  onOpenChange: (bool: boolean) => void;
  serverId: string;
  fetchServerById: (id: string) => void;
};

export default function ModalAddDnsSettings({
  open,
  onOpenChange,
  serverId,
  fetchServerById,
}: ModalAddDnsSettingsProps) {
  const [error, formAction] = useActionState(
    async (prev: string | null, formData: FormData) => {
      const domainRegex = new RegExp(
        "^([a-z0-9]+(-[a-z0-9]+)*\\.)+[a-z]{2,}$",
        "i"
      );
      const newDomain = formData.get("server-domain");

      if (!newDomain) {
        return "Domain name is required";
      }

      if (!domainRegex.test(newDomain.toString())) {
        return "Invalid domain name";
      }
      await addDomainToServerApi(serverId, { domain: newDomain.toString() });
      fetchServerById(serverId);

      onOpenChange(false);
      return null;
    },
    null
  );

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form action={formAction}>
          <DialogHeader>
            <DialogTitle>Connect a domain to this server</DialogTitle>
            <DialogDescription>
              <div>
                To have a deployement ender a domain, you need to connect a
                domain to this server.
              </div>
              <div>
                It is use to generate the ssl certificate and to have a better.
              </div>
            </DialogDescription>
            <Label>Domain name</Label>
            <Input
              id="server-domain"
              name="server-domain"
              placeholder="domain.com"
            />
            <div className="text-red-500">{error}</div>
          </DialogHeader>
          <DialogFooter>
            <Button type="submit">Add domain</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
