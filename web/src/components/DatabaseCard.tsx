import { ServiceDto } from "@/services/getServicesApi";
import { Button } from "./ui/button";
import { Card } from "./ui/card";
import { createServiceApi } from "@/services/createServiceApi";
import { useState } from "react";
import { ButtonStateEnum } from "@/lib/utils";
import SpinnerIcon from "@/assets/SpinnerIcon";

type DatabaseCardProps = {
  service: ServiceDto;
  deployId: string;
};

export default function DatabaseCard({ service, deployId }: DatabaseCardProps) {
  const [createBtnState, setCreateBtnState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  async function createService() {
    setCreateBtnState(ButtonStateEnum.PENDING);
    try {
      await createServiceApi(service.name, deployId);
      setCreateBtnState(ButtonStateEnum.SUCESS);
    } catch (e) {
      console.error(e);
    }
  }

  return (
    <Card className="flex justify-between p-3 w-1/3">
      <div className="flex items-center gap-3">
        <img className="w-10" src={service.icon} />
        <div className="text-xl font-bold">{service.name}</div>
      </div>
      <div>
        <Button
          className={
            ButtonStateEnum.SUCESS === createBtnState ? "bg-green-600" : ""
          }
          onClick={createService}
        >
          {createBtnState === ButtonStateEnum.PENDING ? (
            <SpinnerIcon color="text-white" />
          ) : (
            "Create"
          )}
        </Button>
      </div>
    </Card>
  );
}
