import { ResponseApi, callApi } from "./api";

export async function startApplicationApi(
  appName: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/start/" + appName, "POST");
}
