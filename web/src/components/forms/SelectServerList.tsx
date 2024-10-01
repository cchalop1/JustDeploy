import { ServerDto } from "@/services/getServerListApi";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

type SelectServerListProps = {
  serverList: Array<ServerDto>;
  onServerSelected: (server: ServerDto) => void;
};

export default function SelectServerList({
  serverList,
  onServerSelected,
}: SelectServerListProps) {
  function onValueChange(value: string) {
    const server = serverList.find((server) => server.name === value);
    if (server) {
      onServerSelected(server);
    }
  }

  return (
    <Select onValueChange={onValueChange}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Server 1" />
      </SelectTrigger>
      <SelectContent>
        {serverList.map((server) => (
          <SelectItem key={server.id} value={server.name}>
            {server.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
