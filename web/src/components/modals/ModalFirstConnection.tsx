import { useState, useEffect } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { Copy, Check, RefreshCw } from "lucide-react";
import { useNotification } from "@/hooks/useNotifications";
import { generateRandomApiKey, copyToClipboard } from "@/lib/utils";
import { saveApiKey, getStoredApiKey } from "@/services/authStorage";
import DomainInputCard from "../domain/DomainInputCard";
import { useInfo } from "@/hooks/useInfo";
import { saveInitialSetup } from "@/services/initialSetupApi";

type ModalFirstConnectionProps = {
  serverIp: string;
  onClose?: () => void;
};

export default function ModalFirstConnection({
  serverIp,
  onClose = () => {},
}: ModalFirstConnectionProps) {
  const { serverInfo } = useInfo();

  const [domain, setDomain] = useState<string>("");
  const [copied, setCopied] = useState<boolean>(false);
  const [apiKey, setApiKey] = useState<string>("");
  const [isGeneratingKey, setIsGeneratingKey] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const notif = useNotification();

  // Generate a random API key on component mount or use existing one from storage
  useEffect(() => {
    const storedKey = getStoredApiKey();
    if (storedKey) {
      setApiKey(storedKey);
    } else {
      setApiKey(generateRandomApiKey());
    }
  }, []);

  useEffect(() => {
    setDomain(serverInfo?.server.domain || "");
  }, [serverInfo]);

  const copyApiKey = () => {
    copyToClipboard(apiKey, ({ title }) => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      notif.success({
        title: title,
        content: "API key has been copied to clipboard!",
      });
    });
  };

  const handleGenerateNewKey = () => {
    setIsGeneratingKey(true);
    setTimeout(() => {
      setApiKey(generateRandomApiKey());
      setIsGeneratingKey(false);
    }, 300);
  };

  const handleSubmit = async () => {
    if (!apiKey.trim()) {
      notif.warning({
        title: "Missing API Key",
        content: "Please enter or generate an API key.",
      });
      return;
    }

    if (domain.trim() === "") {
      notif.warning({
        title: "Missing Domain",
        content: "Please enter a domain before continuing.",
      });
      return;
    }

    setIsSubmitting(true);

    try {
      await saveInitialSetup({
        apiKey,
        domain,
      });

      saveApiKey(apiKey);

      notif.success({
        title: "Setup Complete",
        content: "Your API key and domain have been saved successfully.",
      });

      onClose();
    } catch (error) {
      console.error("Error saving initial setup:", error);
      notif.error({
        title: "Setup Failed",
        content: "Failed to save your setup. Please try again.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-[700px] max-w-[90%] max-h-[90vh] overflow-hidden font-mono">
        <div className="p-6 border-b border-gray-200 flex items-center">
          <img src="/hand.png" className="w-16 mr-4" alt="Welcome" />
          <h1 className="text-2xl font-bold">
            Welcome your JustDeploy instance !
          </h1>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-180px)]">
          <div className="mb-6">
            <div className="mb-6">
              <div className="flex justify-between items-center">
                <Label className="text-lg font-medium">Your API Key</Label>
              </div>
              <div className="mt-2 relative">
                <Input
                  type={"text"}
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="Enter your API key or generate one"
                  className="pr-24 font-mono bg-white"
                />
                <div className="absolute right-1 top-1 flex">
                  <Button
                    type="button"
                    variant="ghost"
                    className="h-8 mr-1"
                    onClick={handleGenerateNewKey}
                    disabled={isGeneratingKey}
                    title="Generate new API key"
                  >
                    <RefreshCw
                      color="black"
                      className={`h-4 w-4 ${
                        isGeneratingKey ? "animate-spin" : ""
                      } `}
                    />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    className="h-8"
                    onClick={copyApiKey}
                    title="Copy API key"
                  >
                    {copied ? (
                      <Check className="h-4 w-4 text-black" />
                    ) : (
                      <Copy className="h-4 w-4 text-black" />
                    )}
                  </Button>
                </div>
              </div>
              <div className="text-sm mt-2">
                <p className="text-gray-600">
                  This key is securely stored and will be used for
                  authentication when reconnecting to your JustDeploy instance.
                </p>
              </div>
            </div>

            <DomainInputCard
              domain={domain}
              serverIp={serverIp}
              onDomainChange={setDomain}
              label="Domain Configuration"
            />
          </div>
        </div>

        <div className="p-6 border-t border-gray-200 flex justify-center">
          <Button
            onClick={handleSubmit}
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-8 py-2 rounded-lg"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Saving..." : "Continue"}
          </Button>
        </div>
      </div>
    </div>
  );
}
