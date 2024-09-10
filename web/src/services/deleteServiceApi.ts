import { ResponseApi, callApi } from "./api";

export async function deleteServiceByIdApi(
  serviceId: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>(`/service/${serviceId}`, "DELETE");
}
