import { Button } from "./ui/button";
import { removeApplicationApi } from "@/services/removeApplication";
import { GetDeployConfigResponse } from "@/services/getDeployConfig";
import { ButtonStateEnum } from "@/lib/utils";
import { useState } from "react";
import SpinnerIcon from "@/assets/SpinnerIcon";
import FileIcon from "@/assets/fileIcon";
import { Card } from "./ui/card";
import { Badge } from "./ui/badge";
import LinkIcon from "@/assets/linkIcon";
import ModalApplicationLogs from "./modals/ModalLogs";

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
  const [openLogs, setOpenLogs] = useState(false);

  async function removeApplication() {
    if (deployConfig.appConfig === null) return null;

    setConnectButtonState(ButtonStateEnum.PENDING);
    try {
      await removeApplicationApi(deployConfig.appConfig.name);
      setConnectButtonState(ButtonStateEnum.SUCESS);
      fetchCurrentConfigData();
    } catch (e) {
      console.error(e);
    }
  }
  if (deployConfig.appConfig === null) return null;

  return (
    <>
      <ModalApplicationLogs
        appName={deployConfig.appConfig.name}
        open={openLogs}
        onOpenChange={setOpenLogs}
      />
      <Card className="w-1/2 p-4">
        <div className="flex justify-between">
          <div className="font-bold">{deployConfig.appConfig.name}</div>
          <div className="flex gap-2">
            <Button variant="destructive" onClick={removeApplication}>
              {connectButtonState === ButtonStateEnum.PENDING ? (
                <SpinnerIcon color="text-white" />
              ) : (
                "Delete"
              )}
            </Button>
            <Button>Stop</Button>
            <Button variant="secondary">Redeploy</Button>
          </div>
        </div>
        <Badge className="bg-green-600">Runing</Badge>
        <div className="flex items-center gap-2 mt-4">
          <LinkIcon />
          <a href={deployConfig.url} target="_blank" className="underline">
            {deployConfig.url}
          </a>
        </div>
        <div className="flex items-center gap-2 mt-4">
          <FileIcon />
          <span className="text-sm text-gray-500 dark:text-gray-400">
            {deployConfig.pathToProject}
          </span>
        </div>
        <Button className="mt-4" onClick={() => setOpenLogs(!openLogs)}>
          Logs
        </Button>
      </Card>
    </>
  );
}
