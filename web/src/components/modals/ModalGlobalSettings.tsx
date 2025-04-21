import { useState, useEffect, useRef } from "react";
import Modal from "./Modal";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { getServerInfoApi } from "@/services/getServerInfoApi";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { addDomainToServerApi } from "@/services/addDomainToServerApi";
import { saveTlsServerSettings } from "@/services/saveTlsServerSettings";
import { Github, InfoIcon, XCircle } from "lucide-react";
import { useNotification } from "@/hooks/useNotifications";
import { Card } from "../ui/card";

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
  const [useHttps, setUseHttps] = useState<boolean>(false);
  const [showInfoCard, setShowInfoCard] = useState<boolean>(false);
  const notif = useNotification();
  const timeoutRef = useRef<number | null>(null);
  const emailTimeoutRef = useRef<number | null>(null);

  useEffect(() => {
    async function fetchServerInfo() {
      const serverInfo = await getServerInfoApi();
      setServerIp(serverInfo.ip);
      setServerName(serverInfo.name);
      setDomain(serverInfo.domain);
      setUseHttps(serverInfo.useHttps || false);
      setEmail(serverInfo.email || "");
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
      } catch (e: unknown) {
        notif.error({
          title: "Error",
          content: e instanceof Error ? e.message : "An unknown error occurred",
        });
      }
    }, 1000);
  }

  function onHttpsChange(enabled: boolean) {
    setUseHttps(enabled);
    saveHttpsSettings(enabled, email);
  }

  function onEmailChange(value: string) {
    setEmail(value);

    if (emailTimeoutRef.current) {
      clearTimeout(emailTimeoutRef.current);
    }

    emailTimeoutRef.current = window.setTimeout(() => {
      saveHttpsSettings(useHttps, value);
    }, 1000);
  }

  async function saveHttpsSettings(httpsEnabled: boolean, emailValue: string) {
    try {
      await saveTlsServerSettings({
        useHttps: httpsEnabled,
        email: emailValue,
      });

      notif.success({
        title: "HTTPS settings saved",
        content: httpsEnabled
          ? "HTTPS is now enabled with the provided email"
          : "HTTPS settings updated",
      });
    } catch (e: unknown) {
      notif.error({
        title: "Error",
        content:
          e instanceof Error
            ? e.message
            : "An unknown error occurred while saving HTTPS settings",
      });
    }
  }

  return (
    <Modal
      onClose={onClose}
      headerNode={<h1 className="text-2xl font-bold">Global Settings</h1>}
      className="w-[600px]"
    >
      <div className="flex flex-col h-[calc(100%-3rem)] pl-4 pr-4 pt-2 pb-2 w-full overflow-y-auto">
        <div className="mb-6">
          <div className="mb-4">
            <Label>IP Address</Label>
            <Input value={serverIp} readOnly disabled />
          </div>
          <div className="mb-4">
            <Label>Server Name</Label>
            <Input value={serverName} readOnly disabled />
          </div>

          <div className="mb-4 relative">
            <Label className="mb-1">Domain</Label>
            <div className="flex items-center">
              <Input
                value={domain}
                onChange={(e) => onDomainChange(e.target.value)}
                placeholder="Enter domain (e.g., example.com)"
                className="mb-2"
              />
              <Button
                type="button"
                variant="ghost"
                className="ml-2 p-2 h-10 mb-2"
                onClick={() => setShowInfoCard(!showInfoCard)}
              >
                <InfoIcon className="h-5 w-5 text-blue-600" />
              </Button>
            </div>

            {showInfoCard && (
              <Card className="p-3 bg-blue-50 border-blue-200 absolute z-10 right-0 mt-1 w-full shadow-lg rounded-md">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-2">
                    <InfoIcon className="mt-1 h-4 w-4 flex-shrink-0 text-blue-600" />
                    <div className="text-sm text-blue-800">
                      <p className="font-medium mb-1">
                        How to configure your domain name:
                      </p>
                      <p>
                        Create a DNS record of type A that points to the
                        server's IP address:
                      </p>
                      <code className="bg-blue-100 px-2 py-1 rounded font-mono text-blue-900 block my-2">
                        {domain || "example.com"} → {serverIp || "0.0.0.0"}
                      </code>
                      <p>
                        This configuration will direct your domain to this
                        server.
                      </p>
                      <p className="mt-2">
                        Additionally, create a wildcard CNAME record that points
                        to your domain:
                      </p>
                      <code className="bg-blue-100 px-2 py-1 rounded font-mono text-blue-900 block my-2">
                        *.{domain + "." || "example.com."} →{" "}
                        {domain + "." || "example.com."}
                      </code>
                      <p>
                        This allows subdomains to work with your server
                        automatically.
                      </p>
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    className="p-1 h-6 w-6"
                    onClick={() => setShowInfoCard(false)}
                  >
                    <XCircle className="h-5 w-5 text-blue-600" />
                  </Button>
                </div>
              </Card>
            )}
          </div>

          <div className="mb-6">
            <div className="flex items-center justify-between mb-2">
              <Label htmlFor="https-switch" className="cursor-pointer">
                Enable HTTPS
              </Label>
              <Switch
                id="https-switch"
                checked={useHttps}
                onCheckedChange={onHttpsChange}
              />
            </div>

            {useHttps && (
              <div className="mt-3 ml-1">
                <Label htmlFor="email-input">Email (for SSL certificate)</Label>
                <Input
                  id="email-input"
                  value={email}
                  onChange={(e) => onEmailChange(e.target.value)}
                  placeholder="Enter your email for Let's Encrypt"
                  className="mt-1"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Your email is required to manage the HTTPS certificate and
                  receive important notifications.
                </p>
              </div>
            )}

            {!useHttps && (
              <p className="text-sm text-gray-500 mt-1">
                Enable HTTPS to secure your application and get a free SSL
                certificate via Let's Encrypt.
              </p>
            )}
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
