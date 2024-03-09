import CheckIcon from "@/assets/checkIcon";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { Button } from "./ui/button";
import { removeApplicationApi } from "@/services/removeApplication";
import { GetDeployConfigResponse } from "@/services/getDeployConfig";
import { ButtonStateEnum } from "@/lib/utils";
import { useState } from "react";
import SpinnerIcon from "@/assets/SpinnerIcon";
import FileIcon from "@/assets/fileIcon";

type DeploySuccessProps = {
  deployConfig: GetDeployConfigResponse;
  fetchCurrentConfigData: () => void;
};

export default function DeploySuccess({
  deployConfig,
  fetchCurrentConfigData,
}: DeploySuccessProps) {
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );

  async function removeApplication() {
    setConnectButtonState(ButtonStateEnum.PENDING);
    try {
      await removeApplicationApi();
      setConnectButtonState(ButtonStateEnum.SUCESS);
      fetchCurrentConfigData();
    } catch (e) {
      console.error(e);
    }
  }
  if (deployConfig.appConfig === null) return null;

  return (
    <Alert className="w-2/3">
      <CheckIcon></CheckIcon>
      <AlertTitle>Your App is deploy</AlertTitle>
      <AlertDescription>
        <div>
          Your app <strong>{deployConfig.appConfig.name}</strong> is deploy you
          car check on{" "}
          <a href={deployConfig.url} target="_blank" className="underline">
            {deployConfig.url}
          </a>
          <div className="flex items-center gap-2">
            <FileIcon />
            <span className="text-sm text-gray-500 dark:text-gray-400">
              {deployConfig.pathToProject}
            </span>
          </div>
        </div>
        <Button className="mt-2" onClick={removeApplication}>
          {connectButtonState === ButtonStateEnum.PENDING ? (
            <SpinnerIcon color="text-white" />
          ) : (
            "Stop Application"
          )}
        </Button>
      </AlertDescription>
    </Alert>
  );
}
