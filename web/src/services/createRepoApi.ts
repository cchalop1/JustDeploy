import { callApi } from "./api";
import { ServiceDto } from "./getServicesApi";

export type CreateRepoRequest = {
  repoUrl: string;
};

export async function createRepoApi(
  body: CreateRepoRequest
): Promise<ServiceDto> {
  return await callApi<ServiceDto>("/repo/create", "POST", body);
}
