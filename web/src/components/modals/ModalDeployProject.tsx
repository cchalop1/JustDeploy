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
import { useNotification } from "@/hooks/useNotifications";
import { ProjectDto } from "@/services/getProjectById";

type ModalDeployProjectProps = {
  project: ProjectDto;
  onClose: () => void;
  onDeployProject: (deployProjectDto: DeployProjectDto) => Promise<void>;
};

export default function ModalDeployProject({
  project,
  onClose,
  onDeployProject,
}: ModalDeployProjectProps) {
  const notif = useNotification();
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const [selectedServer, setSelectedServer] = useState<ServerDto | null>();
  const serverUrl = selectedServer?.domain || selectedServer?.ip;
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    getServersListApi().then((servers) => {
      setServerList(servers);
      if (project.serverId) {
        setSelectedServer(servers.find((s) => s.id === project.serverId));
      } else {
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
            <a href={"https://" + serverUrl} target="_blank">
              {serverUrl}
            </a>
          </div>
        </div>
        {/* <div className="flex flex-row items-center gap-2 p-1 mt-2">
          <Checkbox />
          <Label className="font-bold">Setup TLS / HTTPS</Label>
        </div> */}
        <div className="w-full mt-2">
          <Button
            className="w-full"
            onClick={async () => {
              if (!selectedServer) return;
              setIsLoading(true);
              await onDeployProject({
                serverId: selectedServer.id,
                projectId: project.id,
                isTLSDomain: false,
                domain: "",
              });
              setIsLoading(false);
              onClose();
              notif.success({
                title: "Project deployed",
                content: "Your project has been deployed successfully",
              });
            }}
          >
            {isLoading ? <SpinnerIcon color="text-white" /> : "Deploy"}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
