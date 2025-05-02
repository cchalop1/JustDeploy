import { Service } from "@/services/getServicesByDeployId";
import { Card } from "../ui/card";
import { Button } from "../ui/button";
import { useState } from "react";
import { copyToClipboard } from "@/lib/utils";

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
    const env = service.envs
      .map((e: { name: string; value: string }) => `${e.name}=${e.value}`)
      .join("\n");
    copyToClipboard(env, () => {
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    });
  }

  return (
    <Card
      key={service.url}
      className="flex justify-between p-3 mb-3 hover:shadow-md cursor-pointer"
    >
      <div className="flex flex-col gap-3">
        <div className="flex items-center gap-5">
          <img className="w-10" src={service.imageUrl}></img>
          <p className="font-bold">{service.url}</p>
        </div>
      </div>
      <div className="flex gap-3 items-center">
        <Button variant="outline" onClick={copyEnv}>
          {isCopied ? "Copied!" : "Click to copy"}
        </Button>
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
