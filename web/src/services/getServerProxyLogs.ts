import { callApi } from "./api";
import { Logs } from "./getApplicationLogs";

export async function getServerProxyLogs(
  serverId: string
): Promise<Array<Logs>> {
  return await callApi<Array<Logs>>("/server/" + serverId + "/logs", "GET");
}
