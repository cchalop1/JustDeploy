import { useState, useEffect, useRef } from "react";
import Modal from "./Modal";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { addDomainToServerApi } from "@/services/addDomainToServerApi";
import { saveTlsServerSettings } from "@/services/saveTlsServerSettings";
import { Github, Copy, Check } from "lucide-react";
import { useNotification } from "@/hooks/useNotifications";
import { useInfo } from "@/hooks/useInfo";
import DomainInputCard from "../domain/DomainInputCard";
import { copyToClipboard } from "@/lib/utils";

type ModalGlobalSettingsProps = {
  onClose: () => void;
};

export default function ModalGlobalSettings({
  onClose,
}: ModalGlobalSettingsProps) {
  const { serverInfo } = useInfo();
  const [serverIp, setServerIp] = useState<string>("");
  const [domain, setDomain] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [isGithubConnected, setIsGithubConnected] = useState<boolean>(false);
  const [useHttps, setUseHttps] = useState<boolean>(false);
  const [copied, setCopied] = useState<boolean>(false);
  const [isDomainLoading, setIsDomainLoading] = useState<boolean>(false);
  const notif = useNotification();
  const emailTimeoutRef = useRef<number | null>(null);

  useEffect(() => {
    async function fetchServerInfo() {
      if (serverInfo) {
        setServerIp(serverInfo.server.ip);
        setDomain(serverInfo.server.domain);
        setUseHttps(serverInfo.server.useHttps || false);
        setEmail(serverInfo.server.email || "");
      }
    }

    async function fetchGithubConnectionStatus() {
      const { isConnected } = await githubIsConnectedApi();
      setIsGithubConnected(isConnected);
    }

    fetchServerInfo();
    fetchGithubConnectionStatus();
  }, [serverInfo]);

  async function onDomainChange(value: string) {
    setDomain(value);
    setIsDomainLoading(true);

    try {
      // Add 300ms delay to make it feel more realistic
      await new Promise((resolve) => setTimeout(resolve, 300));
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
    } finally {
      setIsDomainLoading(false);
    }
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

  const copyServerIp = () => {
    copyToClipboard(serverIp, ({ title }) => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      notif.success({
        title: title,
        content: "Server IP has been copied to clipboard!",
      });
    });
  };

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
      headerNode={
        <div className="flex items-center">
          <h1 className="text-2xl font-bold">Global Settings</h1>
        </div>
      }
      className="w-[700px]"
    >
      <div className="flex flex-col h-[calc(100%-3rem)] px-6 py-4 w-full overflow-y-auto">
        <div className="mb-6">
          <div className="mb-6">
            <Label className="text-lg font-medium">Server Information</Label>
            <div className="mt-2 relative">
              <Input
                value={serverIp}
                readOnly
                disabled
                className="pr-12 font-mono"
                placeholder="Server IP Address"
              />
              <div className="absolute right-1 top-1">
                <Button
                  type="button"
                  variant="ghost"
                  className="h-8"
                  onClick={copyServerIp}
                  title="Copy Server IP"
                >
                  {copied ? (
                    <Check className="h-4 w-4 text-black" />
                  ) : (
                    <Copy className="h-4 w-4 text-black" />
                  )}
                </Button>
              </div>
            </div>
            <p className="text-sm text-gray-600 mt-1">
              Your server's IP address.
            </p>
          </div>

          <div className="relative">
            <DomainInputCard
              domain={domain}
              serverIp={serverIp}
              onDomainChange={onDomainChange}
              label="Domain Configuration"
              className="mb-6"
            />
            {isDomainLoading && (
              <div className="absolute top-0 right-0 mt-1 mr-12">
                <div className="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
              </div>
            )}
          </div>

          <div className="mb-8">
            <Label className="text-lg font-medium mb-2 block">
              HTTPS Configuration
            </Label>
            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <div className="flex items-center justify-between mb-2">
                <Label
                  htmlFor="https-switch"
                  className="cursor-pointer font-medium"
                >
                  Enable HTTPS
                </Label>
                <Switch
                  id="https-switch"
                  checked={useHttps}
                  onCheckedChange={onHttpsChange}
                />
              </div>

              {useHttps && (
                <div className="mt-3">
                  <Label htmlFor="email-input">
                    Email (for SSL certificate)
                  </Label>
                  <Input
                    id="email-input"
                    value={email}
                    onChange={(e) => onEmailChange(e.target.value)}
                    placeholder="Enter your email for Let's Encrypt"
                    className="mt-1 font-mono"
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    Your email is required to manage the HTTPS certificate and
                    receive important notifications from Let's Encrypt.
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
          </div>

          <div className="mb-2">
            <Label className="text-lg font-medium mb-3 block">
              External Connections
            </Label>
            <div className="flex flex-col space-y-4">
              <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    className={`${
                      isGithubConnected
                        ? "bg-green-100 border-green-400 text-green-700"
                        : "bg-red-100 border-red-400 text-red-700"
                    } font-medium w-52`}
                  >
                    <div className="flex items-center gap-2">
                      <Github className="h-5 w-5" />
                      <div>
                        {isGithubConnected
                          ? "GitHub Connected"
                          : "GitHub Not Connected"}
                      </div>
                    </div>
                  </Button>
                  {!isGithubConnected && (
                    <p className="text-sm text-gray-600 ml-2">
                      Connect GitHub to enable repository deployments.
                    </p>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
}
