import { callApi } from "./api";

export async function getApplicationLogs(appName: string): Promise<string[]> {
  return await callApi<string[]>("/logs/" + appName, "GET");
}
