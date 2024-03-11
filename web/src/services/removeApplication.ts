import { ResponseApi, callApi } from "./api";

export async function removeApplicationApi(
  appName: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/remove/" + appName, "DELETE");
}
