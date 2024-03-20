import { ServerStatusType } from "@/services/getServerListApi";

type ServerStatusProps = {
  status: ServerStatusType;
};

export default function ServerStatus({ status }: ServerStatusProps) {
  let color = "";

  if (status === "Runing") {
    color = "bg-green-500";
  } else if (status === "Installing") {
    color = "bg-blue-500";
  } else {
    color = "bg-red-500";
  }

  return (
    <div className="flex items-center gap-2">
      <div className={`h-3 w-3 rounded-full ${color}`}></div>
      <div>{status}</div>
    </div>
  );
}
