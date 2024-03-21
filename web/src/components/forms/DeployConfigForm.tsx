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

const createDeploymentEmptyState = (): CreateDeployDto => {
  return {
    serverId: "",
    name: "",
    enableTls: false,
    email: null,
    pathToSource: "",
    envs: [{ name: "", secret: "" }],
    deployOnCommit: false,
  };
};

export function DeployConfigForm() {
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [newDeploy, setNewDeploy] = useState<CreateDeployDto>(
    createDeploymentEmptyState()
  );
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const navigate = useNavigate();

  async function fetchServerList() {
    const serverList = await getServersListApi();
    // TODO: check error
    setServerList(serverList);
    if (serverList.length > 0) {
      setNewDeploy({ ...newDeploy, serverId: serverList[0].id });
    }
  }

  useEffect(() => {
    fetchServerList();
  }, []);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!event.target) return;

    setConnectButtonState(ButtonStateEnum.PENDING);

    await createDeployApi(newDeploy);
    setConnectButtonState(ButtonStateEnum.SUCESS);
    navigate("/");
  };

  const addNewEnv = () => {
    setNewDeploy({
      ...newDeploy,
      envs: [...newDeploy.envs, { name: "", secret: "" }],
    });
  };

  const removeEnv = (idx: number) => {
    setNewDeploy({
      ...newDeploy,
      envs: newDeploy.envs.filter((_, index) => index !== idx),
    });
  };

  return (
    <div className="flex justify-center mt-16">
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
                    console.log(value);
                    setNewDeploy({ ...newDeploy, serverId: value });
                  }}
                  defaultValue={newDeploy.serverId || ""}
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
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Name of your project"
                  value={newDeploy.name}
                  onChange={(e) =>
                    setNewDeploy({
                      ...newDeploy,
                      name: e.target.value,
                    })
                  }
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="projectPath">Project foder</Label>
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
                <Label>Env Variables</Label>
                {newDeploy.envs.map((env, idx) => (
                  <div className="flex gap-4">
                    <Input
                      id="envName"
                      name="envName"
                      type="envName"
                      placeholder="Env Name"
                      value={env.name}
                      onChange={(e) => {
                        const updatedEnvs = [...newDeploy.envs];
                        updatedEnvs[idx] = {
                          ...updatedEnvs[idx],
                          name: e.target.value,
                        };
                        setNewDeploy({
                          ...newDeploy,
                          envs: updatedEnvs,
                        });
                      }}
                    />
                    <Input
                      id="envSecret"
                      name="envSecret"
                      type="envSecret"
                      placeholder="Env Secret"
                      value={env.secret}
                      onChange={(e) => {
                        const updatedEnvs = [...newDeploy.envs];
                        updatedEnvs[idx] = {
                          ...updatedEnvs[idx],
                          secret: e.target.value,
                        };
                        setNewDeploy({
                          ...newDeploy,
                          envs: updatedEnvs,
                        });
                      }}
                    />
                    <Button
                      onClick={(e) => {
                        e.stopPropagation();
                        e.preventDefault();
                        if (idx === 0) {
                          addNewEnv();
                        } else {
                          removeEnv(idx);
                        }
                      }}
                    >
                      {idx === 0 ? "+" : "-"}
                    </Button>
                  </div>
                ))}
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
    </div>
  );
}
