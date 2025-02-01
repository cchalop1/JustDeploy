import { GithubIcon, Plus } from "lucide-react";
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
import { redirectToGithubAppRegistration } from "@/services/createGithubApp";
import { GithubRepo } from "@/services/getGithubRepos";

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
  isGithubConnected: boolean;
  githubRepos: Array<GithubRepo>;
  serverIp: string;
};

export default function CommandModal({
  open,
  setOpen,
  create,
  preConfiguredServices,
  isGithubConnected,
  serverIp,
  githubRepos,
}: CommandModalProps) {
  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <CommandInput placeholder="Search a github repos or a services to deploy on your server ..." />
      <CommandList onSelect={() => setOpen(false)}>
        <CommandGroup heading="Github repositories">
          {!isGithubConnected && (
            <CommandItem
              onSelect={() => redirectToGithubAppRegistration(serverIp)}
              className="flex gap-3"
            >
              <GithubIcon className="mr-2 h-5 w-5" />
              <span className="h-4">Connect GitHub</span>
            </CommandItem>
          )}
          {isGithubConnected &&
            githubRepos.map((repo) => (
              <CommandItem
                key={repo.id}
                onSelect={() =>
                  create({ serviceName: repo.name, path: repo.full_name })
                }
                className="flex gap-3"
              >
                <GithubIcon className="mr-2 h-5 w-5" />
                <span className="h-4">{repo.name}</span>
              </CommandItem>
            ))}
        </CommandGroup>
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
