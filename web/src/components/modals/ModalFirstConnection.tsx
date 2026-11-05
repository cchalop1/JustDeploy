import { useState } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { useNotification } from "@/hooks/useNotifications";
import { useInfo } from "@/hooks/useInfo";
import { saveInitialSetup } from "@/services/initialSetupApi";
import { saveToken } from "@/services/authStorage";

type ModalFirstConnectionProps = {
  onClose?: () => void;
};

export default function ModalFirstConnection({
  onClose = () => {},
}: ModalFirstConnectionProps) {
  const { fetchServerInfo } = useInfo();

  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [confirmPassword, setConfirmPassword] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const notif = useNotification();

  const handleSubmit = async () => {
    if (isSubmitting) return;
    if (!email.trim()) {
      notif.warning({ title: "Missing email", content: "Please enter your email address." });
      return;
    }

    if (password.length < 12 || new TextEncoder().encode(password).length > 72) {
      notif.warning({ title: "Missing password", content: "Use a password of at least 12 characters (maximum 72 bytes)." });
      return;
    }

    if (password !== confirmPassword) {
      notif.warning({ title: "Password mismatch", content: "Passwords do not match." });
      return;
    }

    setIsSubmitting(true);

    try {
      const { token } = await saveInitialSetup({ email, password, domain: "" });
      saveToken(token);
      await fetchServerInfo();
      notif.success({ title: "Setup complete", content: "Admin account created successfully." });
      onClose();
    } catch (error) {
      console.error("Error saving initial setup:", error);
      notif.error({ title: "Setup failed", content: error instanceof Error ? error.message : "Setup failed" });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div role="dialog" aria-modal="true" aria-labelledby="setup-title" className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <form onSubmit={(event) => { event.preventDefault(); void handleSubmit(); }} className="bg-white rounded-lg shadow-xl w-[700px] max-w-[90%] max-h-[90vh] overflow-hidden font-mono">
        <div className="p-6 border-b border-gray-200 flex items-center">
          <img src="/hand.png" className="w-16 mr-4" alt="Welcome" />
          <h1 id="setup-title" className="text-2xl font-bold">Welcome to your JustDeploy instance!</h1>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-180px)]">
          <div className="mb-6 space-y-4">
            <p className="text-gray-600 text-sm">Create your admin account. You can deploy a private database immediately and configure your domain later.</p>

            <div>
              <Label htmlFor="admin-email" className="text-sm font-medium">Email</Label>
              <Input
                id="admin-email" autoComplete="username" type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="admin@example.com"
                className="mt-1"
              />
            </div>

            <div>
              <Label htmlFor="admin-password" className="text-sm font-medium">Password</Label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                id="admin-password" autoComplete="new-password" placeholder="At least 12 characters"
                className="mt-1"
              />
            </div>

            <div>
              <Label htmlFor="admin-confirm-password" className="text-sm font-medium">Confirm password</Label>
              <Input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                id="admin-confirm-password" autoComplete="new-password" placeholder="Repeat your password"
                className="mt-1"
              />
            </div>

          </div>
        </div>

        <div className="p-6 border-t border-gray-200 flex justify-center">
          <Button
            type="submit"
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-8 py-2 rounded-lg"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Creating account..." : "Create account"}
          </Button>
        </div>
      </form>
    </div>
  );
}
