import { ResponseApi, callApi } from "./api";

export async function stopApplicationApi(
  appName: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/stop/" + appName, "POST");
}
