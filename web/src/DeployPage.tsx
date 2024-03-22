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
import { Checkbox } from "./components/ui/checkbox";
import { Label } from "./components/ui/label";
import { editDeployementApi } from "./services/editDeploymentApi";

export default function DeployPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [deploy, setDeploy] = useState<DeployDto | null>(null);

  async function fetchDeployById(id: string) {
    const res = await getDeployByIdApi(id);
    setDeploy(res);
  }

  async function onCheckDeployOnCommit(deployOnCommit: boolean, id: string) {
    try {
      await editDeployementApi({
        deployOnCommit,
        id,
      });
      fetchDeployById(id);
    } catch (e) {
      console.error(e);
    }
  }

  useEffect(() => {
    if (id) {
      fetchDeployById(id);
    } else {
      navigate("/");
    }
  }, [id, navigate]);

  if (deploy === null) {
    return null;
  }

  return (
    <div className="mt-40">
      <div className="flex justify-between">
        <div className="text-xl font-bold mb-2">{deploy.name}</div>
        <DeployButtons deploy={deploy} />
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
      <div className="mt-4 flex items-center space-x-2">
        <Checkbox
          id="deploy-on-commit"
          name="deploy-on-commit"
          checked={deploy.deployOnCommit}
          onCheckedChange={(state) => {
            if (typeof state === "boolean" && id) {
              onCheckDeployOnCommit(state, id);
            }
          }}
        />
        <Label
          htmlFor="deploy-on-commit"
          className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
        >
          Deploy on every commit on <code>main</code> branch
        </Label>
      </div>
      <Tabs defaultValue="database-service" className="mt-20">
        <TabsList className="w-full justify-between pl-10 pr-10">
          <TabsTrigger value="database-service">Database Service</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
          <TabsTrigger value="settings">Settings</TabsTrigger>
        </TabsList>
        <TabsContent value="database-service">
          <div>databases</div>
        </TabsContent>
        <TabsContent value="logs">
          <DeployLogs id={deploy.id} />
        </TabsContent>
        <TabsContent value="settings">
          <div>Settings</div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
