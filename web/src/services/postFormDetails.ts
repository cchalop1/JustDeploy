import { callApi } from "./api";

export type Env = {
  name: string;
  value: string;
};

export type CreateDeployDto = {
  serverId: string;
  enableTls: boolean;
  email: string | null;
  pathToSource: string;
  envs: Env[];
  deployOnCommit: boolean;
};

export type CreateDeployResponse = {
  id: string;
};

export async function createDeployApi(
  postCreateDeploymentDto: CreateDeployDto
): Promise<CreateDeployResponse> {
  return await callApi<CreateDeployResponse>(
    "/deploy",
    "POST",
    postCreateDeploymentDto
  );
}
