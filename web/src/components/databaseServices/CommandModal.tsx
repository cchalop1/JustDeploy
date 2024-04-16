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
import { ResponseServiceFromDockerComposeDto } from "@/services/getServicesFromDockerCompose";

type CommandModalProps = {
  open: boolean;
  preConfiguredServices: ServiceDto[];
  serviceFromDockerCompose: ResponseServiceFromDockerComposeDto;
  setOpen: (open: boolean) => void;
  createService: (serviceId: string, fromDockerCompose: boolean) => void;
};

export default function CommandModal({
  open,
  setOpen,
  preConfiguredServices,
  serviceFromDockerCompose,
  createService,
}: CommandModalProps) {
  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Type to search and lauch a service in the list..." />
      <CommandList onSelect={() => setOpen(false)}>
        <CommandEmpty>No results found.</CommandEmpty>
        {serviceFromDockerCompose && (
          <CommandGroup heading="Service from your docker compose file">
            {serviceFromDockerCompose.map((s) => (
              <CommandItem
                className="flex gap-3"
                onSelect={() => createService(s, true)}
              >
                <FileSlidersIcon className="w-5"></FileSlidersIcon>
                <span className="h-4">{s}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        )}
        <CommandSeparator />
        <CommandGroup heading="Other sercices">
          {preConfiguredServices.map((s) => (
            <NewServiceItem
              key={s.name}
              service={s}
              onSelect={(serviceName) => createService(serviceName, false)}
            />
          ))}
          {/* TODO: command item to add new service link to github issue template */}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
