import { Globe, Server } from "lucide-react";
import { useEffect, useState } from "react";
import SelectServerList from "../forms/SelectServerList";
import { ServerDto, getServersListApi } from "@/services/getServerListApi";
import { Button } from "../ui/button";
import { Checkbox } from "../ui/checkbox";
import { Label } from "../ui/label";
import { DeployProjectDto } from "@/services/deployProjectApi";
import SpinnerIcon from "@/assets/SpinnerIcon";
import Modal from "./Modal";

type ModalDeployProjectProps = {
  projectId: string;
  onClose: () => void;
  onDeployProject: (deployProjectDto: DeployProjectDto) => Promise<void>;
};

export default function ModalDeployProject({
  projectId,
  onClose,
  onDeployProject,
}: ModalDeployProjectProps) {
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const [selectedServer, setSelectedServer] = useState<ServerDto | null>();
  const serverUrl = selectedServer?.domain || selectedServer?.ip;
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getServersListApi().then((servers) => {
      setServerList(servers);
      if (servers.length > 0) {
        setSelectedServer(servers[0]);
      }
    });
  }, []);

  if (!selectedServer) {
    return null;
  }

  return (
    <Modal onClose={onClose}>
      <div className="flex flex-col">
        <div className="flex flex-row items-center gap-2 p-1">
          <Server className="w-4 h-4" />
          <div className="font-bold">Server</div>
          <SelectServerList
            serverList={serverList}
            onServerSelected={setSelectedServer}
          />
        </div>
        <div className="flex flex-row items-center gap-2 p-1">
          <Globe className="w-4 h-4" />
          <div className="font-bold">Domain</div>
          <div className="underline">
            <a href={"http://" + serverUrl} target="_blank">
              {serverUrl}
            </a>
          </div>
        </div>
        <div className="flex flex-row items-center gap-2 p-1 mt-2">
          <Checkbox />
          <Label className="font-bold">Setup TLS / HTTPS</Label>
        </div>
        <div className="w-full mt-2">
          <Button
            className="w-full"
            onClick={async () => {
              if (!selectedServer) return;
              setIsLoading(true);
              await onDeployProject({
                serverId: selectedServer.id,
                projectId: projectId,
                isTLSDomain: false,
                domain: "",
              });
              setIsLoading(false);
            }}
          >
            {isLoading ? <SpinnerIcon color="text-white" /> : "Deploy"}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
