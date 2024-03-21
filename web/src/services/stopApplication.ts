import { ResponseApi, callApi } from "./api";

export async function stopApplicationByIdApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/deploy/stop/" + id, "POST");
}
