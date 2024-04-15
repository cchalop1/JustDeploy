import { ResponseApi, callApi } from "./api";

type CreateServiceApi = {
  serviceName: string;
  fromDockerCompose: boolean;
};

export async function createServiceApi(
  deployId: string,
  body: CreateServiceApi,
) {
  return await callApi<ResponseApi>(
    `/deploy/${deployId}/service`,
    "POST",
    body,
  );
}
