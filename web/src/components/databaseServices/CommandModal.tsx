import { Plus } from "lucide-react";
import {
  CommandDialog,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "../ui/command";
import NewServiceItem from "./NewServiceItem";
import { ServiceDto } from "@/services/getServicesApi";

export type CreateServiceFunc = (parms: {
  serviceName?: string;
  path?: string;
  fromDockerCompose?: boolean;
}) => void;

type CommandModalProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
  preConfiguredServices: ServiceDto[];
  create: CreateServiceFunc;
};

export default function CommandModal({
  open,
  setOpen,
  create,
  preConfiguredServices,
}: CommandModalProps) {
  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Search a github repos or a services to deploy on your server ..." />
      <CommandList onSelect={() => setOpen(false)}>
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
