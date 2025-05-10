import { callApi } from "./api";
import { Logs } from "./getApplicationLogs";

export async function getServiceRunLogs(serviceId: string): Promise<Logs[]> {
  return await callApi<Logs[]>(`/service/${serviceId}/run-logs`, "GET");
}
