import { Check, Terminal } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { DeployConfigDto } from "@/services/getDeployConfig";

type ConfigDeployInfosProps = {
  config: DeployConfigDto;
};

export default function ConfigDeployInfos({ config }: ConfigDeployInfosProps) {
  return (
    <Alert className="w-[500px]">
      <Terminal />
      <AlertTitle>
        JustDeploy load informations from your local folder
      </AlertTitle>
      <AlertDescription>
        <div className="flex flex-col">
          {config.dockerFileFound && (
            <div className="flex">
              <Check /> Dockerfile
            </div>
          )}
          {config.composeFileFound && (
            <div className="flex">
              <Check /> Docker compose file
            </div>
          )}
          {config.envFileFound && (
            <div className="flex">
              <Check /> Envs load from <code>.env</code>
            </div>
          )}
        </div>
      </AlertDescription>
    </Alert>
  );
}
