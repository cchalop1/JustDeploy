import { ResponseApi, callApi } from "./api";

export async function addServiceDatabasesForDeployApi(
  deployId: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>(
    "/deploy/" + deployId + "/addService",
    "POST"
  );
}
