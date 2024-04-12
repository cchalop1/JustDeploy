import { ResponseApi, callApi } from "./api";

export async function deleteServiceApi(
  deployId: string,
  serviceId: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>(
    `/deploy/${deployId}/service/${serviceId}`,
    "DELETE"
  );
}
