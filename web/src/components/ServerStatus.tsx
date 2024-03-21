import { ServerStatusType } from "@/services/getServerListApi";
import { Badge } from "./ui/badge";

type ServerStatusProps = {
  status: ServerStatusType;
};

export default function Status({ status }: ServerStatusProps) {
  let color = "";

  if (status === "Runing") {
    color = "bg-green-500";
  } else if (status === "Installing") {
    color = "bg-blue-500";
  } else {
    color = "bg-red-500";
  }

  return <Badge className={color}>{status}</Badge>;
}
