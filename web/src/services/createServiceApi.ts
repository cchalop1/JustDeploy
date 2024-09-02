import { ResponseApi, callApi } from "./api";

type CreateServiceApi = {
  serviceName: string;
  fromDockerCompose: boolean;
  deployId: string | undefined;
};

export async function createServiceApi(body: CreateServiceApi) {
  return await callApi<ResponseApi>(`/deploy/service`, "POST", body);
}
