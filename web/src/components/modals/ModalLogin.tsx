import { useState } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { useNotification } from "@/hooks/useNotifications";
import { loginApi } from "@/services/loginApi";
import { saveToken } from "@/services/authStorage";

type ModalLoginProps = {
  onClose?: () => void;
};

export default function ModalLogin({ onClose = () => {} }: ModalLoginProps) {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const notif = useNotification();

  const handleSubmit = async () => {
    if (!email.trim() || !password.trim()) {
      notif.warning({ title: "Missing fields", content: "Please enter your email and password." });
      return;
    }

    setIsSubmitting(true);

    try {
      const { token } = await loginApi({ email, password });
      saveToken(token);
      notif.success({ title: "Logged in", content: "Welcome back!" });
      onClose();
    } catch {
      notif.error({ title: "Login failed", content: "Invalid email or password." });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") handleSubmit();
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-[400px] max-w-[90%] font-mono">
        <div className="p-6 border-b border-gray-200 flex items-center">
          <img src="/hand.png" className="w-12 mr-4" alt="Login" />
          <h1 className="text-xl font-bold">Sign in to JustDeploy</h1>
        </div>

        <div className="p-6 space-y-4">
          <div>
            <Label className="text-sm font-medium">Email</Label>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="admin@example.com"
              className="mt-1"
              autoFocus
            />
          </div>

          <div>
            <Label className="text-sm font-medium">Password</Label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Your password"
              className="mt-1"
            />
          </div>
        </div>

        <div className="p-6 border-t border-gray-200 flex justify-center">
          <Button
            onClick={handleSubmit}
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-8 py-2 rounded-lg w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Signing in..." : "Sign in"}
          </Button>
        </div>
      </div>
    </div>
  );
}
