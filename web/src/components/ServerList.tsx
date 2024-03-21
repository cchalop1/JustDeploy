import { ServerDto } from "@/services/getServerListApi";
import { Card, CardContent } from "./ui/card";
import Status from "./ServerStatus";

type ServerListProps = {
  serverList: Array<ServerDto>;
};

export default function ServerList({ serverList }: ServerListProps) {
  if (serverList.length === 0) {
    return (
      <div className="h-full flex justify-center pt-12">
        You don't have any connected server yet
      </div>
    );
  }
  return (
    <div className="flex gap-3 h-2/3 mt-2">
      {serverList.map((server) => (
        <Card
          className="hover:shadow-md cursor-pointer w-80 h-full pt-4 pl-2"
          key={server.id}
        >
          <CardContent className="flex flex-col gap-2">
            <div className="flex justify-between ">
              <div className="font-bold">{server.name}</div>
              <Status status={server.status} />
            </div>
            <div className="underline">{server.domain}</div>
            {/* <div>ip: {server.ip}</div> */}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
