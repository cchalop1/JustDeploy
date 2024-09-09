import { FileSlidersIcon, Folder } from "lucide-react";
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
  createService: (serviceName: string, fromDockerCompose: boolean) => void;
  currentPath?: string;
};

export default function CommandModal({
  open,
  setOpen,
  preConfiguredServices,
  serviceFromDockerCompose,
  createService,
  currentPath,
}: CommandModalProps) {
  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Type to search and lauch a service in the list..." />
      <CommandList onSelect={() => setOpen(false)}>
        <CommandEmpty>No results found.</CommandEmpty>
        {currentPath && (
          <CommandGroup heading="Load your current folder">
            <CommandItem
              className="flex gap-3"
              // TODO: edit to load the current folder as a service
              onSelect={() => createService(currentPath, false)}
            >
              <Folder className="w-5" />
              <span className="h-4">{currentPath}</span>
            </CommandItem>
          </CommandGroup>
        )}
        {serviceFromDockerCompose && (
          <CommandGroup heading="Service from your docker compose file">
            {serviceFromDockerCompose.map((s) => (
              <CommandItem
                className="flex gap-3"
                onSelect={() => createService(s, true)}
                key={s}
              >
                <FileSlidersIcon className="w-5"></FileSlidersIcon>
                <span className="h-4">{s}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        )}
        <CommandSeparator />
        <CommandGroup heading="Other databases">
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
