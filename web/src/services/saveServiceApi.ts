import { ResponseApi, callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export async function saveServiceApi(service: Service) {
  return await callApi<ResponseApi>(`/service`, "PUT", service);
}
