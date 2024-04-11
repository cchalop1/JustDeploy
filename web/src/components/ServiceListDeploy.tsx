import {
  getServicesByDeployIdApi,
  Service,
} from "@/services/getServicesByDeployId";
import { useEffect, useState } from "react";
import { Card } from "./ui/card";
import { Button } from "./ui/button";

type ServiceListDeployProps = {
  deployId: string;
};

export default function ServiceListDeploy({
  deployId,
}: ServiceListDeployProps) {
  const [services, setServices] = useState<Service[]>([]);

  async function fetchServiceList(deployId: string) {
    const res = await getServicesByDeployIdApi(deployId);
    setServices(res);
  }

  useEffect(() => {
    fetchServiceList(deployId);
  }, [deployId]);

  return (
    <div className="flex flex-col mb-10">
      {services.map((s) => (
        <Card className="flex justify-between p-3">
          <div className="flex flex-col gap-3">
            <div className="flex items-center gap-5">
              <div className="text-xl font-bold">{s.name}</div>
              <div>{s.status}</div>
            </div>
            <div>{s.imageName}</div>
          </div>
          <div>
            <Button variant="destructive">Delete</Button>
          </div>
        </Card>
      ))}
    </div>
  );
}
