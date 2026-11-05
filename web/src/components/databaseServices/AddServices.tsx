import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";

import { useNotification } from "@/hooks/useNotifications";
import CommandModal from "./CommandModal";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { GithubRepo, getGithubRepos } from "@/services/getGithubRepos";
import { createRepoApi } from "@/services/createRepoApi";
import { Service } from "@/services/getServicesByDeployId";
import { createDatabaseApi } from "@/services/createDatabaseApi";

type AddServiceProps = {
  setLoading: (loading: boolean) => void;
  fetchServices: () => Promise<void>;
  firstService?: boolean;
  onCreated?: (service: Service) => void;
};

export default function AddService({
  setLoading,
  fetchServices,
  firstService = false,
  onCreated,
}: AddServiceProps) {
  const notif = useNotification();
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const [isGithubConnected, setIsGithubConnected] = useState(false);
  const [githubRepos, setGithubRepos] = useState<Array<GithubRepo>>([]);
  const [category, setCategory] = useState<"all" | "application" | "database">("all");
  const [openCommandModal, setOpenCommandModal] = useState(false);

  async function getServices() {
    const res = await getPreConfiguredServiceListApi();
    setPreConfiguredServices(res);
  }

  async function fetchIsGithubConnected() {
    const { isConnected } = await githubIsConnectedApi();
    setIsGithubConnected(isConnected);
    if (isConnected) {
      const repos = await getGithubRepos();
      setGithubRepos(repos || []);
    }
  }

  async function createRepoToDeploy(repoUrl: string) {
    try {
      setLoading(true);
      setOpenCommandModal(false);

      const created = await createRepoApi({
        repoUrl,
      });

      await fetchServices();
      onCreated?.(created);
    } catch (error) {
      console.error("Failed to create repo:", error);
      notif.error({title:"Creation failed",content:error instanceof Error ? error.message : "Creation failed"});
    } finally {
      setLoading(false);
    }
  }

  async function createDatabaseToDeploy(databaseName: string) {
    try {
      setLoading(true);
      setOpenCommandModal(false);

      const created = await createDatabaseApi({
        databaseName,
      });

      await fetchServices();
      onCreated?.(created);
    } catch (error) {
      console.error("Failed to create database:", error);
      notif.error({title:"Creation failed",content:error instanceof Error ? error.message : "Creation failed"});
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    getServices().catch((error: unknown) => console.error("Service templates unavailable", error));
    fetchIsGithubConnected().catch((error: unknown) => console.error("GitHub repositories unavailable", error));

    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setCategory("all");
        setOpenCommandModal((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      {firstService ? (
        <div className="flex flex-wrap justify-center gap-4">
          <button type="button" onClick={() => { setCategory("application"); setOpenCommandModal(true); }} className="text-left bg-white rounded-xl border p-6 w-72 hover:shadow-md">
            <span className="block text-lg font-semibold">Deploy my application</span>
            <span className="text-sm text-gray-600">Connect GitHub and choose a repository. Public HTTPS access is required.</span>
          </button>
          <button type="button" onClick={() => { setCategory("database"); setOpenCommandModal(true); }} className="text-left bg-white rounded-xl border p-6 w-72 hover:shadow-md">
            <span className="block text-lg font-semibold">Create a database</span>
            <span className="text-sm text-gray-600">Start with a private database. No domain or GitHub account required.</span>
          </button>
        </div>
      ) : (
        <button type="button" onClick={() => { setCategory("all"); setOpenCommandModal(true); }} className="text-left hover:shadow-md bg-white p-5 w-80 min-h-36 rounded-xl border border-dashed border-gray-500">
          <span className="block font-bold text-xl mb-1">Add a service</span>
          <span>Choose a GitHub repository or a service template. Shortcut: ⌘ / Ctrl + K.</span>
        </button>
      )}
      <CommandModal
        category={category}
        open={openCommandModal}
        setOpen={setOpenCommandModal}
        preConfiguredServices={preConfiguredServices}
        isGithubConnected={isGithubConnected}
        createRepoToDeploy={createRepoToDeploy}
        createDatabaseToDeploy={createDatabaseToDeploy}
        githubRepos={githubRepos}
      />
    </>
  );
}
