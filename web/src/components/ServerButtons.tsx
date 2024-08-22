import SpinnerIcon from "@/assets/SpinnerIcon";
import { ButtonStateEnum } from "@/lib/utils";
import { Button } from "./ui/button";
import { useNavigate } from "react-router-dom";
import { Suspense, useState } from "react";
import { removeServerByIdApi } from "@/services/removeServerByIdApi";
import ModalAddDnsSettings from "./modals/ModalAddDnsSettings";
import { ServerDto } from "@/services/getServerListApi";
import ModalLogs from "./modals/ModalLogs";
import { Logs } from "@/services/getApplicationLogs";
import { getServerProxyLogs } from "@/services/getServerProxyLogs";

type ServerButtonsProps = {
  server: ServerDto;
  fetchServerById: (id: string) => void;
};

export default function ServerButtons({
  server,
  fetchServerById,
}: ServerButtonsProps) {
  const navigate = useNavigate();
  const [removeServerButtonState, setRemoveServerButtonState] =
    useState<ButtonStateEnum>(ButtonStateEnum.INIT);
  const [openAddDomainModal, setOpenAddDomainModal] = useState<boolean>(false);
  const [openProxyLogsModal, setProxyLogsModal] = useState<boolean>(false);

  async function removeServerById() {
    setRemoveServerButtonState(ButtonStateEnum.PENDING);
    try {
      await removeServerByIdApi(server.id);
      setRemoveServerButtonState(ButtonStateEnum.SUCESS);
      navigate("/");
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchServerProxyLogsById(): Promise<Array<Logs>> {
    const logs = await getServerProxyLogs(server.id);
    return logs;
  }

  return (
    <div className="flex gap-2">
      <ModalAddDnsSettings
        open={openAddDomainModal}
        serverId={server.id}
        fetchServerById={fetchServerById}
        onOpenChange={(o) => setOpenAddDomainModal(o)}
      />
      <Suspense fallback={<></>}>
        <ModalLogs
          open={openProxyLogsModal}
          fetchServerProxyLogs={fetchServerProxyLogsById}
          onOpenChange={(o) => setProxyLogsModal(o)}
        />
      </Suspense>
      <Button variant="outline" onClick={() => setProxyLogsModal(true)}>
        Proxy Logs
      </Button>
      {!server.domain && (
        <Button variant="outline" onClick={() => setOpenAddDomainModal(true)}>
          Add Domain
        </Button>
      )}
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
