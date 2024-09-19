import { ResponseApi, callApi } from "./api";

export type CreateServiceApi = {
  serviceName?: string;
  fromDockerCompose?: boolean;
  path?: string;
  projectId: string;
};

export async function createServiceApi(body: CreateServiceApi) {
  return await callApi<ResponseApi>(`/deploy/service`, "POST", body);
}
