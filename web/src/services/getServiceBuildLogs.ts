import { callApi } from "./api";
import { Logs } from "./getApplicationLogs";

export async function getServiceBuildLogs(serviceId: string): Promise<Logs[]> {
  return await callApi<Logs[]>(`/service/${serviceId}/build-logs`, "GET");
}
