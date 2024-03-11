import { ResponseApi, callApi } from "./api";

export async function reDeployAppApi(appName: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/redeploy/" + appName, "POST");
}
