import { ResponseApi, callApi } from "./api";

export async function createServiceApi(serviceName: string, deployId: string) {
  return await callApi<ResponseApi>(
    `/deploy/${deployId}/service/${serviceName}`,
    "POST"
  );
}
