import { DeployDto } from "@/services/getDeployListApi";
import { Button } from "./ui/button";
import { ButtonStateEnum } from "@/lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { useState } from "react";
import { stopApplicationByIdApi } from "@/services/stopApplication";
import { startApplicationApi } from "@/services/startApplication";
import { reDeployAppApi } from "@/services/reDeployApp";
import { useNavigate } from "react-router-dom";
import { removeApplicationApi } from "@/services/removeApplication";

type DeployButtonsProps = {
  deploy: DeployDto;
};

export default function DeployButtons({ deploy }: DeployButtonsProps) {
  const navigate = useNavigate();
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [redeployButtonState, setReDeployButtonState] =
    useState<ButtonStateEnum>(ButtonStateEnum.INIT);
  const [stopStartButtonState, setStopStartButtonState] =
    useState<ButtonStateEnum>(ButtonStateEnum.INIT);

  async function removeApplication() {
    setConnectButtonState(ButtonStateEnum.PENDING);
    try {
      await removeApplicationApi(deploy.id);
      setConnectButtonState(ButtonStateEnum.SUCESS);
      navigate("/");
    } catch (e) {
      console.error(e);
    }
  }

  async function reDeployApplication() {
    setReDeployButtonState(ButtonStateEnum.PENDING);
    try {
      await reDeployAppApi(deploy.id);
      setReDeployButtonState(ButtonStateEnum.SUCESS);
    } catch (e) {
      console.error(e);
    }
  }

  async function startStopApplication() {
    setStopStartButtonState(ButtonStateEnum.PENDING);
    try {
      {
        deploy.status === "Runing"
          ? await stopApplicationByIdApi(deploy.id)
          : await startApplicationApi(deploy.id);
      }
      setStopStartButtonState(ButtonStateEnum.SUCESS);
    } catch (e) {
      console.error(e);
    }
  }

  return (
    <div className="flex gap-2">
      <Button variant="destructive" onClick={removeApplication}>
        {connectButtonState === ButtonStateEnum.PENDING ? (
          <SpinnerIcon color="text-white" />
        ) : (
          "Delete"
        )}
      </Button>
      <Button onClick={startStopApplication}>
        {stopStartButtonState === ButtonStateEnum.PENDING ? (
          <SpinnerIcon color="text-white" />
        ) : deploy.status === "Runing" ? (
          "Stop"
        ) : (
          "Start"
        )}
      </Button>
      <Button variant="secondary" onClick={() => reDeployApplication()}>
        {redeployButtonState === ButtonStateEnum.PENDING ? (
          <SpinnerIcon color="text-black" />
        ) : (
          "Redeploy"
        )}
      </Button>
    </div>
  );
}
