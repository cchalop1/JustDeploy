import { callApi, ResponseApi } from "./api";

export type DeployProjectDto = {
  projectId: string;
  serverId: string;
  domain: string | null;
  isTLSDomain: boolean;
};

export async function deployProjectApi(deployProjectDto: DeployProjectDto) {
  return await callApi<ResponseApi>(
    "/deploy/project",
    "POST",
    deployProjectDto
  );
}
