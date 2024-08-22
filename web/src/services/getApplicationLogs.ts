import { callApi } from "./api";

export type Logs = {
  date: string;
  message: string;
};

export async function getApplicationLogs(appName: string): Promise<Logs[]> {
  return await callApi<Logs[]>("/deploy/logs/" + appName, "GET");
}
