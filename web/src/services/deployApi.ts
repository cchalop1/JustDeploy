import { callApi, ResponseApi } from "./api";

export async function deployApi() {
  return await callApi<ResponseApi>("/deploy", "POST");
}

export async function deployServiceApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>(`/deploy/service/${encodeURIComponent(id)}`, "POST");
}
