import { callApi } from "./api";

export async function getApplicationLogs(appName: string): Promise<string[]> {
  return await callApi<string[]>("/deploy/logs/" + appName, "GET");
}
