import { ResponseApi, callApi } from "./api";

export async function reDeployAppApi(id: string): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/deploy/redeploy/" + id, "POST");
}
