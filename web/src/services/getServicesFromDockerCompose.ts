import { callApi } from "./api";

export type ResponseServiceFromDockerComposeDto = Array<string> | null;

export async function getServicesFromDockerComposeApi(
  deployId: string,
): Promise<ResponseServiceFromDockerComposeDto> {
  return await callApi<ResponseServiceFromDockerComposeDto>(
    `/deploy/${deployId}/service-docker-compose`,
    "GET",
  );
}
