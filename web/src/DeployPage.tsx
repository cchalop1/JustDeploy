import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { DeployDto } from "./services/getDeployListApi";
import { getDeployByIdApi } from "./services/getDeployById";
import DeployButtons from "./components/DeployButtons";
import Status from "./components/ServerStatus";
import LinkIcon from "./assets/linkIcon";
import FolderIcon from "./assets/FolderIcon";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./components/ui/tabs";
import DeployLogs from "./components/DeployLogs";
import DeploySettings from "./components/DeploySettings";
import { ServerDto } from "./services/getServerListApi";
import { getServerByIdApi } from "./services/getServerById";
import DeployDatabaseServices from "./components/DeployDatabaseServices";

export default function DeployPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [deploy, setDeploy] = useState<DeployDto | null>(null);
  const [server, setServer] = useState<ServerDto | null>(null);

  async function fetchDeployById(id: string) {
    const res = await getDeployByIdApi(id);
    setDeploy(res);
    return res;
  }

  async function fetchServer(serverId: string) {
    const res = await getServerByIdApi(serverId);
    setServer(res);
  }

  useEffect(() => {
    if (id) {
      fetchDeployById(id).then((res) => fetchServer(res.serverId));
    } else {
      navigate("/");
    }
  }, [id, navigate]);

  if (deploy === null || server === null) {
    return null;
  }

  return (
    <div className="mt-40">
      <div className="flex justify-between">
        <div className="text-xl font-bold mb-2">{deploy.name}</div>
        <DeployButtons deploy={deploy} fetchDeployById={fetchDeployById} />
      </div>
      <Status status={deploy.status} />
      <div className="flex items-center mt-2 gap-2">
        <LinkIcon />
        <a href={deploy.url} target="_blank" className="underline">
          {deploy.url}
        </a>
      </div>
      <div className="flex items-center mt-2 gap-2">
        <FolderIcon />
        {deploy.pathToSource}
      </div>
      <Tabs defaultValue="database-service" className="mt-20">
        <TabsList className="w-full justify-between pl-5 pr-5">
          <TabsTrigger value="database-service">Database Service</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>
        <TabsContent value="database-service">
          <DeployDatabaseServices deploy={deploy} />
        </TabsContent>
        <TabsContent value="logs">
          <DeployLogs id={deploy.id} />
        </TabsContent>
        <TabsContent value="settings">
          <DeploySettings
            serverDomain={server.domain}
            deploy={deploy}
            fetchDeployById={fetchDeployById}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
