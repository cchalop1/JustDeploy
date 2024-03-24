import { DeployDto } from "@/services/getDeployListApi";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { ButtonStateEnum } from "@/lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { addServiceDatabasesForDeployApi } from "@/services/addServiceDatabasesForDeployApi";
import { useState } from "react";

type DeployDatabaseServicesProps = {
  deploy: DeployDto;
};
export default function DeployDatabaseServices({
  deploy,
}: DeployDatabaseServicesProps) {
  // TODO: get from backend
  const databasesServices = ["redis"];

  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );

  async function addServiceDatabasesForDeploy() {
    setConnectButtonState(ButtonStateEnum.PENDING);
    try {
      await addServiceDatabasesForDeployApi(deploy.id);
      setConnectButtonState(ButtonStateEnum.SUCESS);
    } catch (e) {
      console.error(e);
    }
  }

  return (
    <div className="flex gap-3">
      {databasesServices.map((d) => (
        <Card className="w-1/2">
          <CardHeader>
            <CardTitle className="flex m-2">
              <img
                className="h-10 w-10"
                src="https://logowik.com/content/uploads/images/redis.jpg"
              />
              <div>{d}</div>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Button onClick={addServiceDatabasesForDeploy}>
              {connectButtonState === ButtonStateEnum.PENDING ? (
                <SpinnerIcon color="text-white" />
              ) : (
                "Add"
              )}
            </Button>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
