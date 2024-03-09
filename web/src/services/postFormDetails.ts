import { ResponseApi, callApi } from "./api";

export type Env = {
  name: string;
  secret: string;
};

export type PostCreateDeploymentDto = {
  name: string;
  enableTls: boolean;
  email: string | null;
  pathToSource: string;
  envs: Env[];
};

export async function postFormDetails(
  postCreateDeploymentDto: PostCreateDeploymentDto
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/deploy", "POST", postCreateDeploymentDto);
}
