import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "../ui/checkbox";
import { useEffect, useState } from "react";
import {
  CreateDeployDto,
  Env,
  createDeployApi,
} from "../../services/postFormDetails";
import { ButtonStateEnum } from "../../lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";
import {
  SelectValue,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from "../ui/select";
import { ServerDto, getServersListApi } from "@/services/getServerListApi";
import { useNavigate } from "react-router-dom";
import { DeployConfigDto, getDeployConfig } from "@/services/getDeployConfig";
import EnvsManagements from "./EnvsManagements";
import ConfigDeployInfos from "../ConfigDeployInfos";
import AlertDestructive from "../alerts/AlertDestructive";

const createDeploymentEmptyState = (): CreateDeployDto => {
  return {
    serverId: "",
    enableTls: false,
    email: null,
    pathToSource: "",
    envs: [{ name: "", value: "" }],
    deployOnCommit: false,
  };
};

export function CreateDeployForm() {
  const navigate = useNavigate();
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const [config, setConfig] = useState<DeployConfigDto | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [newDeploy, setNewDeploy] = useState<CreateDeployDto>(
    createDeploymentEmptyState()
  );

  async function fetchServerList() {
    const resServerList = await getServersListApi();
    // TODO: check error
    setServerList(resServerList);
    if (resServerList.length > 0) {
      setNewDeploy((d) => ({ ...d, serverId: resServerList[0].id }));
    }
  }

  async function fetchDeployConfig(path: string | null = null) {
    const deployConfig = await getDeployConfig(path);
    const envs = deployConfig.envs;
    setNewDeploy((d) => ({
      ...d,
      pathToSource: deployConfig.pathToSource,
      envs: envs,
    }));
    setConfig(deployConfig);
    envs.push({ name: "", value: "" });
  }

  useEffect(() => {
    fetchDeployConfig();
    fetchServerList();
  }, []);

  useEffect(() => {
    fetchDeployConfig(newDeploy.pathToSource);
  }, [newDeploy.pathToSource]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!event.target) return;

    setConnectButtonState(ButtonStateEnum.PENDING);

    try {
      const res = await createDeployApi(newDeploy);
      setConnectButtonState(ButtonStateEnum.SUCESS);
      navigate(`/deploy/${res.id}/installation`);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (e: any) {
      setError(e.message);
      setTimeout(() => {
        setError(null);
      }, 5000);
      setConnectButtonState(ButtonStateEnum.INIT);
    }
  };

  function setEnvs(envs: Array<Env>) {
    setNewDeploy({ ...newDeploy, envs: envs });
  }

  return (
    <div className="flex flex-col gap-4 mt-16 items-center">
      {error && <AlertDestructive message={error} />}
      {config && <ConfigDeployInfos config={config} />}
      <Card className="w-[500px]">
        <CardHeader>
          <CardTitle>Deploy project</CardTitle>
          <CardDescription>
            Deploy your new project in one-click.
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent>
            <div className="grid w-full items-center gap-4">
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="auth-methode">Select server for your app</Label>
                <Select
                  onValueChange={(value) => {
                    setNewDeploy({ ...newDeploy, serverId: value });
                  }}
                  defaultValue={newDeploy.serverId}
                >
                  <SelectTrigger>
                    <SelectValue
                      placeholder={
                        serverList.length > 0 ? serverList[0].name : "Server"
                      }
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {serverList.map((s) => (
                      <SelectItem key={s.id} value={s.id}>
                        {s.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="projectPath">
                  Path to your project folder or Dockerfile
                </Label>
                <Input
                  id="projectPath"
                  name="projectPath"
                  // type="file"
                  placeholder="/path/to/your/source"
                  value={newDeploy.pathToSource}
                  onChange={(e) =>
                    setNewDeploy({
                      ...newDeploy,
                      pathToSource: e.target.value,
                    })
                  }
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="enable-tls"
                  name="enable-tls"
                  value={newDeploy.enableTls.toString()}
                  onCheckedChange={(isChecked) => {
                    if (isChecked === "indeterminate") return;
                    setNewDeploy({
                      ...newDeploy,
                      enableTls: isChecked,
                      email: isChecked ? "" : null,
                    });
                  }}
                />
                <Label
                  htmlFor="enable-tls"
                  className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Setup HTTPS for this deploy
                </Label>
              </div>
              {newDeploy.enableTls && newDeploy.email !== null && (
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    name="email"
                    type="email"
                    placeholder="Email for tls setup"
                    value={newDeploy.email}
                    onChange={(e) =>
                      setNewDeploy({
                        ...newDeploy,
                        email: e.target.value,
                      })
                    }
                  />
                </div>
              )}
              <div className="space-y-2">
                <EnvsManagements envs={newDeploy.envs} setEnvs={setEnvs} />
              </div>
            </div>
          </CardContent>
          <CardFooter className="flex justify-between">
            <Button type="submit" className="w-full">
              {connectButtonState === ButtonStateEnum.PENDING ? (
                <SpinnerIcon color="text-white" />
              ) : (
                "Click to deploy your application"
              )}
            </Button>
          </CardFooter>
        </form>
      </Card>
      {config && config.composeFileFound && <Card>fjlkze</Card>}
    </div>
  );
}
