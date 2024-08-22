import { use } from "react";
import { getApplicationLogs } from "@/services/getApplicationLogs";
import DisplayLogs from "./DisplayLogs";

type DeployLogsProps = {
  id: string;
};

export default function DeployLogs({ id }: DeployLogsProps) {
  const logs = use(getApplicationLogs(id));

  return <DisplayLogs logs={logs} />;
}
