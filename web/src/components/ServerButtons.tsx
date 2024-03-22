import SpinnerIcon from "@/assets/SpinnerIcon";
import { ButtonStateEnum } from "@/lib/utils";
import { Button } from "./ui/button";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { removeServerByIdApi } from "@/services/removeServerByIdApi";

type ServerButtonsProps = {
  serverId: string;
};

export default function ServerButtons({ serverId }: ServerButtonsProps) {
  const navigate = useNavigate();
  const [removeServerButtonState, setRemoveServerButtonState] =
    useState<ButtonStateEnum>(ButtonStateEnum.INIT);

  async function removeServerById() {
    setRemoveServerButtonState(ButtonStateEnum.PENDING);
    try {
      await removeServerByIdApi(serverId);
      setRemoveServerButtonState(ButtonStateEnum.SUCESS);
      navigate("/");
    } catch (e) {
      console.error(e);
    }
  }

  return (
    <div className="flex gap-2">
      <Button variant="destructive" onClick={removeServerById}>
        {removeServerButtonState === ButtonStateEnum.PENDING ? (
          <SpinnerIcon color="text-white" />
        ) : (
          "Delete"
        )}
      </Button>
    </div>
  );
}
