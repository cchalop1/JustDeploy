import { Service } from "@/services/getServicesByDeployId";
import { Card } from "../ui/card";
import Status from "../ServerStatus";
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
    <Card
      key={service.name}
      className="flex justify-between p-3 mb-3 hover:shadow-md cursor-pointer"
    >
      <div className="flex flex-col gap-3">
        <div className="flex items-center gap-5">
          <img className="w-10" src={service.imageUrl}></img>
          <p className="font-bold">{service.name}</p>
          <Status status={service.status} />
        </div>
      </div>
      <div className="flex gap-3 items-center">
        Click to copy
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
