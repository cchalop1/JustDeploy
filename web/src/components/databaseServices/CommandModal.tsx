import { FileSlidersIcon, Folder, Plus } from "lucide-react";
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

export type CreateServiceFunc = (parms: {
  serviceName?: string;
  path?: string;
  fromDockerCompose?: boolean;
}) => void;

type CommandModalProps = {
  open: boolean;
  preConfiguredServices: ServiceDto[];
  serviceFromDockerCompose: ResponseServiceFromDockerComposeDto;
  setOpen: (open: boolean) => void;
  create: CreateServiceFunc;
  currentPath?: string;
};

export default function CommandModal({
  open,
  setOpen,
  preConfiguredServices,
  serviceFromDockerCompose,
  create,
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
              onSelect={() => create({ path: currentPath })}
            >
              <Folder className="w-5" />
              <span className="h-4">{currentPath}</span>
            </CommandItem>
          </CommandGroup>
        )}
        {serviceFromDockerCompose && serviceFromDockerCompose.length > 0 && (
          <CommandGroup heading="Service from your docker compose file">
            {serviceFromDockerCompose.map((s) => (
              <CommandItem
                className="flex gap-3"
                onSelect={() =>
                  create({ serviceName: s, fromDockerCompose: true })
                }
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
              onSelect={(serviceName) => create({ serviceName })}
            />
          ))}

          <CommandItem
            onSelect={() => {
              window.open("https://github.com/cchalop1/JustDeploy/issues/new");
            }}
            className="flex gap-3"
          >
            <Plus className="w-5" />
            <span className="h-4">Add new databases</span>
          </CommandItem>
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
