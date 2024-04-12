import { ServiceDto } from "@/services/getServicesApi";
import { CommandItem } from "../ui/command";

type DatabaseCardProps = {
  service: ServiceDto;
  onSelect: (serviceId: string) => void;
};

export default function NewServiceItem({
  service,
  onSelect,
}: DatabaseCardProps) {
  return (
    <CommandItem
      onSelect={() => onSelect(service.name)}
      className="flex gap-3"
      key={service.name}
      value={service.name}
    >
      <img src={service.icon} className="w-5" />
      <span className="h-4">{service.name}</span>
    </CommandItem>
  );
}
