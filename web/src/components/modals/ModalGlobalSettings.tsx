import { useState, useEffect, useRef } from "react";
import Modal from "./Modal";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { getServerInfoApi } from "@/services/getServerInfoApi";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { addDomainToServerApi } from "@/services/addDomainToServerApi";
import { Github } from "lucide-react";
import { useNotification } from "@/hooks/useNotifications";

type ModalGlobalSettingsProps = {
  onClose: () => void;
  onClickNewServer: () => void;
};

export default function ModalGlobalSettings({
  onClose,
}: ModalGlobalSettingsProps) {
  const [serverIp, setServerIp] = useState<string>("");
  const [serverName, setServerName] = useState<string>("");
  const [domain, setDomain] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [isGithubConnected, setIsGithubConnected] = useState<boolean>(false);
  const notif = useNotification();
  const timeoutRef = useRef<number | null>(null);

  useEffect(() => {
    async function fetchServerInfo() {
      const serverInfo = await getServerInfoApi();
      setServerIp(serverInfo.ip);
      setServerName(serverInfo.name);
      setDomain(serverInfo.domain);
    }

    async function fetchGithubConnectionStatus() {
      const { isConnected } = await githubIsConnectedApi();
      setIsGithubConnected(isConnected);
    }

    fetchServerInfo();
    fetchGithubConnectionStatus();
  }, []);

  function onDomainChange(value: string) {
    setDomain(value);

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = window.setTimeout(async () => {
      try {
        await addDomainToServerApi({ domain: value });
        notif.success({
          title: "Domain saved",
          content: "Domain has been successfully saved!",
        });
      } catch (e) {
        notif.error({
          title: "Error",
          content: e.message,
        });
      }
    }, 1000);
  }

  return (
    <Modal
      onClose={onClose}
      headerNode={<h1 className="text-2xl font-bold">Global Settings</h1>}
      className="w-1/3"
    >
      <div className="flex flex-col h-[calc(100%-3rem)] pl-4 pr-4 pt-2 pb-2 w-full">
        <div className="mb-20">
          <div className="mb-4">
            <Label>IP Address</Label>
            <Input value={serverIp} readOnly disabled />
          </div>
          <div className="mb-4">
            <Label>Server Name</Label>
            <Input value={serverName} readOnly disabled />
          </div>
          <div className="mb-4">
            <Label>Domain</Label>
            <Input
              value={domain}
              onChange={(e) => onDomainChange(e.target.value)}
              placeholder="Enter domain"
            />
          </div>
          <div className="mb-4">
            <Label>Email</Label>
            <Input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter email"
            />
          </div>
          <div className="text-2xl font-bold mb-4">Connections</div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              className={`${
                isGithubConnected ? "bg-green-500" : "bg-red-500"
              } text-white font-bold`}
            >
              {isGithubConnected ? (
                <div className="flex items-center gap-2">
                  <Github />
                  <div>Connected</div>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <Github />
                  <div>Not Connected</div>
                </div>
              )}
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
