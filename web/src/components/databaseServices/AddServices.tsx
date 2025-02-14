import {
  ServiceDto,
  getPreConfiguredServiceListApi,
} from "@/services/getServicesApi";
import { useEffect, useState } from "react";

import CommandModal from "./CommandModal";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { getServerInfoApi } from "@/services/getServerInfoApi";
import { GithubRepo, getGithubRepos } from "@/services/getGithubRepos";
import { createRepoApi } from "@/services/createRepoApi";

type AddServiceProps = {
  setLoading: (loading: boolean) => void;
  fetchServices: () => void;
};

export default function AddService({
  setLoading,
  fetchServices,
}: AddServiceProps) {
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const [isGithubConnected, setIsGithubConnected] = useState(false);
  const [githubRepos, setGithubRepos] = useState<Array<GithubRepo>>([]);
  const [serverIp, setServerIp] = useState<string>("");
  const text =
    "Connect a github repos or create a new service. You can also press";
  const [openCommandModal, setOpenCommandModal] = useState(false);

  async function fetchServerInfo() {
    const serverInfo = await getServerInfoApi();
    setServerIp(serverInfo.ip);
  }

  async function getServices() {
    const res = await getPreConfiguredServiceListApi();
    setPreConfiguredServices(res);
  }

  async function fetchIsGithubConnected() {
    const { isConnected } = await githubIsConnectedApi();
    setIsGithubConnected(isConnected);
    if (isConnected) {
      const repos = await getGithubRepos();
      setGithubRepos(repos);
    }
  }

  async function createRepoToDeploy(repoUrl: string) {
    try {
      setLoading(true);
      setOpenCommandModal(false);

      await createRepoApi({
        repoUrl,
      });

      await fetchServices();
    } catch (error) {
      console.error("Failed to create repo:", error);
      // You might want to add error handling/notification here
    } finally {
      setLoading(false);
    }
  }

  async function createDatabaseToDeploy(databaseName: string) {
    setLoading(true);
    setOpenCommandModal(false);
    console.log(databaseName);
    // await createService({
    //   path: createServiceParams.path,
    //   serviceName: createServiceParams.serviceName,
    // });
    setLoading(false);
    await getServices();
  }

  useEffect(() => {
    getServices();
    fetchIsGithubConnected();

    fetchServerInfo();
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpenCommandModal((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  return (
    <>
      <div
        onClick={() => setOpenCommandModal(true)}
        className="hover:shadow-md text-gray-600 hover:text-black bg-white cursor-pointer pt-5 pb-6 pl-5 pr-5 w-80 h-36 rounded-lg border border-dashed border-gray-500 hover:border-black"
      >
        <div className="font-bold text-xl mb-1">Click here</div>
        <div>
          {text}{" "}
          <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
            <span className="text-xs">⌘</span>K
          </kbd>
        </div>
      </div>
      <CommandModal
        open={openCommandModal}
        setOpen={setOpenCommandModal}
        preConfiguredServices={preConfiguredServices}
        isGithubConnected={isGithubConnected}
        createRepoToDeploy={createRepoToDeploy}
        createDatabaseToDeploy={createDatabaseToDeploy}
        githubRepos={githubRepos}
        serverIp={serverIp}
      />
    </>
  );
}
