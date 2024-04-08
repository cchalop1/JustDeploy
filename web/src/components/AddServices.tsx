import { ServiceDto, getServiceListApi } from "@/services/getServicesApi";
import { useEffect, useState } from "react";
import DatabaseCard from "./DatabaseCard";
import { Card } from "./ui/card";
import { Button } from "./ui/button";
import { Database } from "lucide-react";

type AddServiceProps = {
  dockerComposeIsFound: boolean;
};

export default function AddService({ dockerComposeIsFound }: AddServiceProps) {
  const [services, setServices] = useState<Array<ServiceDto>>([]);

  async function getServices() {
    const res = await getServiceListApi();
    setServices(res);
  }

  useEffect(() => {
    getServices();
  }, []);

  return (
    <div className="flex flex-col gap-2">
      {services.map((s) => (
        <DatabaseCard service={s} />
      ))}
      {dockerComposeIsFound && (
        <Card className="flex justify-between p-3">
          <div className="flex items-center gap-3">
            <Database />
            <div className="text-xl font-bold">
              From your Docker compose file
            </div>
          </div>
          <div>
            <Button>Create</Button>
          </div>
        </Card>
      )}
    </div>
  );
}