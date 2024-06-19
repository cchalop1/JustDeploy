import SpinnerIcon from "@/assets/SpinnerIcon";
import { ButtonStateEnum } from "@/lib/utils";
import { Button } from "./ui/button";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { removeServerByIdApi } from "@/services/removeServerByIdApi";
import ModalAddDnsSettings from "./modals/ModalAddDnsSettings";

type ServerButtonsProps = {
  serverId: string;
};

export default function ServerButtons({ serverId }: ServerButtonsProps) {
  const navigate = useNavigate();
  const [removeServerButtonState, setRemoveServerButtonState] =
    useState<ButtonStateEnum>(ButtonStateEnum.INIT);
  const [openAddDomainModal, setOpenAddDomainModal] = useState<boolean>(false);

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
      <ModalAddDnsSettings
        open={openAddDomainModal}
        serverId={serverId}
        onOpenChange={(o) => setOpenAddDomainModal(o)}
      />
      <Button variant="outline" onClick={() => setOpenAddDomainModal(true)}>
        Add Domain
      </Button>
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
