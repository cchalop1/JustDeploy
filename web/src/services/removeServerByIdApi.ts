import { ResponseApi, callApi } from "./api";

export async function removeServerByIdApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/server/remove/" + id, "DELETE");
}
