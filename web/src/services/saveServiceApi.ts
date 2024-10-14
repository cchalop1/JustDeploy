import { ResponseApi, callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export async function saveServiceApi(service: Service, projectId: string) {
  return await callApi<ResponseApi>(
    `/project/${projectId}/service`,
    "PUT",
    service
  );
}
