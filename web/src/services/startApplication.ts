import { ResponseApi, callApi } from "./api";

export async function startApplicationApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/deploy/start/" + id, "POST");
}
