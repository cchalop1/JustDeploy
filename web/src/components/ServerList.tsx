import { ServerDto } from "@/services/getServerListApi";
import { Card, CardContent } from "./ui/card";
import Status from "./ServerStatus";
import { useNavigate } from "react-router-dom";

type ServerListProps = {
  serverList: Array<ServerDto>;
};

export default function ServerList({ serverList }: ServerListProps) {
  const navigate = useNavigate();

  if (serverList.length === 0) {
    return (
      <div className="h-full flex justify-center pt-12">
        You don't have any connected server yet
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-2 gap-3 h-2/3 mt-2">
      {serverList.map((server) => (
        <Card
          className="hover:shadow-md cursor-pointer h-full pt-4 pl-2"
          key={server.id}
          onClick={() => navigate("server/" + server.id)}
        >
          <CardContent className="flex flex-col gap-2">
            <div className="flex justify-between">
              <div className="font-bold">{server.name}</div>
              <Status status={server.status} />
            </div>
            <div className="underline">{server.domain}</div>
            <div>{server.ip}</div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
