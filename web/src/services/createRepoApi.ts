import { callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export type CreateRepoRequest = {
  repoUrl: string;
};

export async function createRepoApi(
  body: CreateRepoRequest
): Promise<Service> {
  return await callApi<Service>("/repo/create", "POST", body);
}
