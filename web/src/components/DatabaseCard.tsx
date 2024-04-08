import { ServiceDto } from "@/services/getServicesApi";
import { Button } from "./ui/button";
import { Card } from "./ui/card";

type DatabaseCardProps = {
  service: ServiceDto;
};

export default function DatabaseCard({ service }: DatabaseCardProps) {
  return (
    <Card className="flex justify-between p-3">
      <div className="flex items-center gap-3">
        <img className="w-10" src={service.image} />
        <div className="text-xl font-bold">{service.name}</div>
      </div>
      <div>
        <Button>Create</Button>
      </div>
    </Card>
  );
}
