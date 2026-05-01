import { useState, useEffect } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { useNotification } from "@/hooks/useNotifications";
import DomainInputCard from "../domain/DomainInputCard";
import { useInfo } from "@/hooks/useInfo";
import { saveInitialSetup } from "@/services/initialSetupApi";
import { saveToken } from "@/services/authStorage";

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
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [confirmPassword, setConfirmPassword] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const notif = useNotification();

  useEffect(() => {
    setDomain(serverInfo?.server.domain || "");
  }, [serverInfo]);

  const handleSubmit = async () => {
    if (!email.trim()) {
      notif.warning({ title: "Missing email", content: "Please enter your email address." });
      return;
    }

    if (!password.trim()) {
      notif.warning({ title: "Missing password", content: "Please enter a password." });
      return;
    }

    if (password !== confirmPassword) {
      notif.warning({ title: "Password mismatch", content: "Passwords do not match." });
      return;
    }

    if (domain.trim() === "") {
      notif.warning({ title: "Missing domain", content: "Please enter a domain before continuing." });
      return;
    }

    setIsSubmitting(true);

    try {
      const { token } = await saveInitialSetup({ email, password, domain });
      saveToken(token);
      notif.success({ title: "Setup complete", content: "Admin account created successfully." });
      onClose();
    } catch (error) {
      console.error("Error saving initial setup:", error);
      notif.error({ title: "Setup failed", content: "Failed to create admin account. Please try again." });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-[700px] max-w-[90%] max-h-[90vh] overflow-hidden font-mono">
        <div className="p-6 border-b border-gray-200 flex items-center">
          <img src="/hand.png" className="w-16 mr-4" alt="Welcome" />
          <h1 className="text-2xl font-bold">Welcome to your JustDeploy instance!</h1>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-180px)]">
          <div className="mb-6 space-y-4">
            <p className="text-gray-600 text-sm">Create your admin account to get started.</p>

            <div>
              <Label className="text-sm font-medium">Email</Label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@example.com"
                className="mt-1"
              />
            </div>

            <div>
              <Label className="text-sm font-medium">Password</Label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Choose a strong password"
                className="mt-1"
              />
            </div>

            <div>
              <Label className="text-sm font-medium">Confirm password</Label>
              <Input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Repeat your password"
                className="mt-1"
              />
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
            {isSubmitting ? "Creating account..." : "Create account"}
          </Button>
        </div>
      </div>
    </div>
  );
}
