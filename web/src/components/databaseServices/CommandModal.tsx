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
import { useState } from "react";
import { baseUrl } from "@/services/api";

type CommandModalProps = {
  category?: "all" | "application" | "database";
  open: boolean;
  setOpen: (open: boolean) => void;
  preConfiguredServices: Array<ServiceDto>;
  isGithubConnected: boolean;
  createRepoToDeploy: (repoUrl: string) => Promise<void>;
  createDatabaseToDeploy: (databaseName: string) => Promise<void>;
  githubRepos: Array<GithubRepo>;
};

export default function CommandModal({
  category = "all",
  open,
  setOpen,
  preConfiguredServices,
  isGithubConnected,
  githubRepos,
  createDatabaseToDeploy,
  createRepoToDeploy,
}: CommandModalProps) {
  const isHttps = new URL(baseUrl || window.location.origin).protocol === "https:";
  const llmServices = preConfiguredServices.filter((s) => s.type === "llm");
  const databaseServices = preConfiguredServices.filter(
    (s) => s.type === "database"
  );
  const analyticsServices = preConfiguredServices.filter(
    (s) => s.type === "analytics"
  );
  const [connectionError, setConnectionError] = useState("");

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      {connectionError && <p role="alert" className="p-3 text-red-600">{connectionError}</p>}
      <CommandInput placeholder="Search repositories or service templates…" />
      {category !== "database" && !isHttps && (
        <div className="p-4 text-sm space-y-2 border-b">
          <p className="font-semibold">Set up public HTTPS before connecting GitHub</p>
          <p>Point an administration domain to your VPS, configure an HTTPS reverse proxy to 127.0.0.1:5915, then reopen JustDeploy at that HTTPS address.</p>
          <a href="https://github.com/cchalop1/JustDeploy/blob/main/SECURITY-UPGRADE.md#administration-and-first-setup" target="_blank" rel="noreferrer" className="underline">Read the HTTPS setup instructions</a>
          <p>You can deploy a private database now without this step.</p>
        </div>
      )}
      <CommandList>
        {category !== "database" && <CommandGroup heading="GitHub repositories">
          {!isGithubConnected && isHttps && (
            <CommandItem
              onSelect={() =>
                redirectToGithubAppRegistration().catch((e: unknown) => setConnectionError(e instanceof Error ? e.message : "GitHub connection failed"))
              }
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
                onSelect={() => createRepoToDeploy(repo.full_name)}
                className="flex gap-3"
              >
                <GithubIcon className="mr-2 h-5 w-5" />
                <span className="h-4">{repo.name}</span>
              </CommandItem>
            ))}
        </CommandGroup>}
        <CommandSeparator />

        {/* LLM Models Section */}
        {category === "all" && <CommandGroup heading="LLM Models">
          {llmServices.map((s) => (
            <NewServiceItem
              key={s.name}
              service={s}
              onSelect={(serviceName) => createDatabaseToDeploy(serviceName)}
            />
          ))}
        </CommandGroup>}
        <CommandSeparator />

        {category !== "application" && <CommandGroup heading="Databases">
          {databaseServices.map((s) => (
            <NewServiceItem
              key={s.name}
              service={s}
              onSelect={(serviceName) => createDatabaseToDeploy(serviceName)}
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
        </CommandGroup>}
        <CommandSeparator />

        {category === "all" && <CommandGroup heading="Analytics">
          {analyticsServices.map((s) => (
            <NewServiceItem
              key={s.name}
              service={s}
              onSelect={(serviceName) => createDatabaseToDeploy(serviceName)}
            />
          ))}
        </CommandGroup>}
      </CommandList>
    </CommandDialog>
  );
}
