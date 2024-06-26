import { AlertCircle } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

type AlertDestructiveProps = {
  message: string;
};

export default function AlertDestructive({ message }: AlertDestructiveProps) {
  return (
    <Alert
      variant="destructive"
      className="fixed top-0 right-0 m-4 bg-white w-1/3 z-10"
    >
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  );
}
