import { ResponseApi, callApi } from "./api";

export async function deleteServiceByIdApi(
  projectId: string,
  serviceId: string
): Promise<ResponseApi> {
  return await callApi<ResponseApi>(
    `/project/${projectId}/service/${serviceId}`,
    "DELETE"
  );
}
