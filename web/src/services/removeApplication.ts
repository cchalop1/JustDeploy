import { ResponseApi, callApi } from "./api";

export async function removeApplicationApi(): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/remove", "DELETE");
}
