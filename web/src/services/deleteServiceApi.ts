import { ResponseApi, callApi } from "./api";

export async function deleteServiceApi(
  // TODO: delete deployId
  serviceId: string,
  deployId?: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>(
    `/deploy/${deployId}/service/${serviceId}`,
    "DELETE"
  );
}
