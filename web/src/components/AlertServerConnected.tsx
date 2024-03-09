import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import CheckIcon from "@/assets/checkIcon";

export default function AlertServerConnected() {
  return (
    <Alert className="m-10 w-[300px]">
      <CheckIcon></CheckIcon>
      <AlertTitle>Server connected</AlertTitle>
      <AlertDescription>
        <div>
          Your server is connected and setup correcly to deploy your app
        </div>
        <ul className="list-disc">
          <li>server connect</li>
          <li>docker install</li>
          <li>certificates generated</li>
          <li>docker tls port open</li>
          <li>docker client is connected</li>
        </ul>
      </AlertDescription>
    </Alert>
  );
}
