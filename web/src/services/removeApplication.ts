import { ResponseApi, callApi } from "./api";

export async function removeApplicationApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/remove/" + id, "DELETE");
}
