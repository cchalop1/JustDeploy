import { ResponseApi, callApi } from "./api";

export type Env = {
  name: string;
  secret: string;
};

export type CreateDeployDto = {
  serverId: string;
  name: string;
  enableTls: boolean;
  email: string | null;
  pathToSource: string;
  envs: Env[];
  deployOnCommit: boolean;
};

export async function createDeployApi(
  postCreateDeploymentDto: CreateDeployDto
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/deploy", "POST", postCreateDeploymentDto);
}
