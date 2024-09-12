import { Check, X } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "../ui/alert";

type AlertModalProps = {
  title: string;
  message: string;
  type: "error" | "warning" | "info" | "success";
};

export default function AlertModal({ message, title, type }: AlertModalProps) {
  return (
    <Alert className="absolute top-4 left-1/2 transform -translate-x-1/2 w-full max-w-md ">
      {type === "error" && <X className="h-5 w-5" />}
      {type === "success" && <Check className="h-5 w-5" />}
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  );
}
