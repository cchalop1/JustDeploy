import { FileSlidersIcon } from "lucide-react";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "../ui/command";
import NewServiceItem from "./NewServiceItem";
import { ServiceDto } from "@/services/getServicesApi";

type CommandModalProps = {
  open: boolean;
  services: ServiceDto[];
  setOpen: (open: boolean) => void;
  createService: (serviceId: string) => void;
  createServiceFromCompose: () => void;
  composeFileFound?: boolean;
};

export default function CommandModal({
  open,
  setOpen,
  services,
  createService,
  createServiceFromCompose,
  composeFileFound,
}: CommandModalProps) {
  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Type to search and lauch a service in the list..." />
      <CommandList onSelect={() => setOpen(false)}>
        <CommandEmpty>No results found.</CommandEmpty>
        {composeFileFound && (
          <CommandGroup heading="Local Config">
            <CommandItem
              className="flex gap-3"
              value="compose"
              onSelect={createServiceFromCompose}
            >
              <FileSlidersIcon className="w-5"></FileSlidersIcon>
              <span className="h-4">From your docker-compose file</span>
            </CommandItem>
          </CommandGroup>
        )}
        <CommandSeparator />
        <CommandGroup heading="Other sercices">
          {services.map((s) => (
            <NewServiceItem key={s.name} service={s} onSelect={createService} />
          ))}
          {/* TODO: command item to add new service link to github issue template */}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
