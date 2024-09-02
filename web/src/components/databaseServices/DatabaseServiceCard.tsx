import { Service } from "@/services/getServicesByDeployId";
import { Card } from "../ui/card";
import Status from "../ServerStatus";
import {
  TooltipProvider,
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "../ui/tooltip";
import { Copy } from "lucide-react";
import { Button } from "../ui/button";
import { useState } from "react";

type DatabaseServiceCardProps = {
  service: Service;
  setServiceToDelete: (s: Service) => void;
};

export default function DatabaseServiceCard({
  service,
  setServiceToDelete,
}: DatabaseServiceCardProps) {
  const [isCopied, setIsCopied] = useState(false);

  function copyEnv() {
    const env = service.envs.map((e) => `${e.name}=${e.value}`).join("\n");
    navigator.clipboard.writeText(env);
    setIsCopied(true);
  }

  return (
    <Card key={service.name} className="flex justify-between p-3 mb-3 ">
      <div className="flex flex-col gap-3">
        <div className="flex items-center gap-5">
          <img className="w-10" src={service.imageUrl}></img>
          <p className="font-bold">{service.name}</p>
          <Status status={service.status} />
        </div>
      </div>
      <div className="flex gap-3 items-center">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger
              className="hover:bg-gray-200 p-1 rounded cursor-pointer"
              onClick={copyEnv}
            >
              <Copy className="w-6 h-6" />
            </TooltipTrigger>
            <TooltipContent>
              <p>
                {isCopied
                  ? "Environment variable successfully copied!"
                  : "Click to copy environment variable"}
              </p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
        <Button
          variant="destructive"
          onClick={() => setServiceToDelete(service)}
        >
          Delete
        </Button>
      </div>
    </Card>
  );
}
