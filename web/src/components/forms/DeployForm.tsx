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
import { useState } from "react";
import {
  PostCreateDeploymentDto,
  postFormDetails,
} from "../../services/postFormDetails";
import { ButtonStateEnum } from "../../lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { GetDeployConfigResponse } from "@/services/getDeployConfig";

type FromToDeployProps = {
  deployConfig: GetDeployConfigResponse;
  fetchCurrentConfigData: () => void;
};

export function AppConfigForm({
  fetchCurrentConfigData,
  deployConfig,
}: FromToDeployProps) {
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [createDeployement, setCreateDeployement] =
    useState<PostCreateDeploymentDto>({
      name: "",
      enableTls: false,
      email: null,
      pathToSource: deployConfig.appConfig?.pathToSource || "",
      envs: [{ name: "", secret: "" }],
    });

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!event.target) return;

    setConnectButtonState(ButtonStateEnum.PENDING);

    await postFormDetails(createDeployement);
    setConnectButtonState(ButtonStateEnum.SUCESS);
    fetchCurrentConfigData();
  };

  const addNewEnv = () => {
    setCreateDeployement({
      ...createDeployement,
      envs: [...createDeployement.envs, { name: "", secret: "" }],
    });
  };

  const removeEnv = (idx: number) => {
    setCreateDeployement({
      ...createDeployement,
      envs: createDeployement.envs.filter((_, index) => index !== idx),
    });
  };

  return (
    <>
      <Card className="w-[500px] m-10">
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
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Name of your project"
                  value={createDeployement.name}
                  onChange={(e) =>
                    setCreateDeployement({
                      ...createDeployement,
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
                  value={createDeployement.pathToSource}
                  onChange={(e) =>
                    setCreateDeployement({
                      ...createDeployement,
                      pathToSource: e.target.value,
                    })
                  }
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="enable-tls"
                  name="enable-tls"
                  value={createDeployement.enableTls.toString()}
                  onCheckedChange={(isChecked) => {
                    if (isChecked === "indeterminate") return;
                    setCreateDeployement({
                      ...createDeployement,
                      enableTls: isChecked,
                      email: isChecked ? "" : null,
                    });
                  }}
                />
                <Label
                  htmlFor="enable-tls"
                  className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Setup this application with tls connection
                </Label>
              </div>
              {createDeployement.enableTls &&
                createDeployement.email !== null && (
                  <div className="flex flex-col space-y-1.5">
                    <Label htmlFor="email">Email</Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      placeholder="Email for tls setup"
                      value={createDeployement.email}
                      onChange={(e) =>
                        setCreateDeployement({
                          ...createDeployement,
                          email: e.target.value,
                        })
                      }
                    />
                  </div>
                )}
              <div className="space-y-2">
                <Label>Env Variables</Label>
                {createDeployement.envs.map((env, idx) => (
                  <div className="flex gap-4">
                    <Input
                      id="envName"
                      name="envName"
                      type="envName"
                      placeholder="Env Name"
                      value={env.name}
                      onChange={(e) => {
                        const updatedEnvs = [...createDeployement.envs];
                        updatedEnvs[idx] = {
                          ...updatedEnvs[idx],
                          name: e.target.value,
                        };
                        setCreateDeployement({
                          ...createDeployement,
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
                        const updatedEnvs = [...createDeployement.envs];
                        updatedEnvs[idx] = {
                          ...updatedEnvs[idx],
                          secret: e.target.value,
                        };
                        setCreateDeployement({
                          ...createDeployement,
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
    </>
  );
}
