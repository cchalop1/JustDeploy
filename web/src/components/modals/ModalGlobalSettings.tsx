import { useState, useEffect, useRef } from "react";
import Modal from "./Modal";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
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

function SectionTitle({ children }: { children: React.ReactNode }) {
  return (
    <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-3">
      {children}
    </p>
  );
}

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
    if (serverInfo) {
      setServerIp(serverInfo.server.ip);
      setDomain(serverInfo.server.domain);
      setUseHttps(serverInfo.server.useHttps || false);
      setEmail(serverInfo.server.email || "");
    }

    githubIsConnectedApi().then(({ isConnected }) =>
      setIsGithubConnected(isConnected)
    );
  }, [serverInfo]);

  async function onDomainChange(value: string) {
    setDomain(value);
    setIsDomainLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 300));
      await addDomainToServerApi({ domain: value });
      notif.success({ title: "Domain saved", content: "Domain updated successfully." });
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
    if (emailTimeoutRef.current) clearTimeout(emailTimeoutRef.current);
    emailTimeoutRef.current = window.setTimeout(() => {
      saveHttpsSettings(useHttps, value);
    }, 1000);
  }

  const copyServerIp = () => {
    copyToClipboard(serverIp, ({ title }) => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      notif.success({ title, content: "Server IP copied to clipboard." });
    });
  };

  async function saveHttpsSettings(httpsEnabled: boolean, emailValue: string) {
    try {
      await saveTlsServerSettings({ useHttps: httpsEnabled, email: emailValue });
      notif.success({
        title: "HTTPS settings saved",
        content: httpsEnabled ? "HTTPS enabled." : "HTTPS settings updated.",
      });
    } catch (e: unknown) {
      notif.error({
        title: "Error",
        content: e instanceof Error ? e.message : "An unknown error occurred",
      });
    }
  }

  return (
    <Modal
      onClose={onClose}
      headerNode={
        <h1 className="text-base font-semibold text-gray-900">Global Settings</h1>
      }
      className="w-[520px]"
    >
      <div className="divide-y divide-gray-100">

        {/* Server IP */}
        <div className="px-5 py-4">
          <SectionTitle>Server</SectionTitle>
          <div className="flex items-center gap-2 bg-gray-50 border border-gray-200 rounded-md px-3 py-2">
            <span className="flex-1 font-mono text-sm text-gray-700 select-all">
              {serverIp || "—"}
            </span>
            <button
              onClick={copyServerIp}
              className="text-gray-400 hover:text-gray-700 transition-colors"
              title="Copy IP"
            >
              {copied ? (
                <Check className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
            </button>
          </div>
          <p className="text-xs text-gray-400 mt-1.5">Your server's public IP address.</p>
        </div>

        {/* Domain */}
        <div className="px-5 py-4 relative">
          <SectionTitle>Domain</SectionTitle>
          <DomainInputCard
            domain={domain}
            serverIp={serverIp}
            onDomainChange={onDomainChange}
            label=""
          />
          {isDomainLoading && (
            <div className="absolute top-4 right-5">
              <div className="w-3.5 h-3.5 border-2 border-blue-500 border-t-transparent rounded-full animate-spin" />
            </div>
          )}
        </div>

        {/* HTTPS */}
        <div className="px-5 py-4">
          <SectionTitle>HTTPS</SectionTitle>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-800">Enable HTTPS</p>
              <p className="text-xs text-gray-400 mt-0.5">
                Free SSL certificate via Let's Encrypt.
              </p>
            </div>
            <Switch
              id="https-switch"
              checked={useHttps}
              onCheckedChange={onHttpsChange}
            />
          </div>

          {useHttps && (
            <div className="mt-4 space-y-1.5">
              <Label htmlFor="email-input" className="text-xs text-gray-500">
                Email for SSL certificate
              </Label>
              <Input
                id="email-input"
                value={email}
                onChange={(e) => onEmailChange(e.target.value)}
                placeholder="you@example.com"
                className="text-sm"
              />
              <p className="text-xs text-gray-400">
                Used by Let's Encrypt to manage your certificate.
              </p>
            </div>
          )}
        </div>

        {/* GitHub */}
        <div className="px-5 py-4">
          <SectionTitle>Integrations</SectionTitle>
          <div className="flex items-center gap-3">
            <div
              className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium border ${
                isGithubConnected
                  ? "bg-green-50 border-green-200 text-green-700"
                  : "bg-red-50 border-red-200 text-red-600"
              }`}
            >
              <Github className="h-3.5 w-3.5" />
              {isGithubConnected ? "GitHub connected" : "GitHub not connected"}
            </div>
            {!isGithubConnected && (
              <p className="text-xs text-gray-400">
                Connect GitHub to enable repository deployments.
              </p>
            )}
          </div>
        </div>

      </div>
    </Modal>
  );
}
