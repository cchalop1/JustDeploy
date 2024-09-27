import { callApi, ResponseApi } from "./api";

type DeployProjectDto = {
  projectId: string;
  serverId: string;
};

export async function deployProjectApi(deployProjectDto: DeployProjectDto) {
  return await callApi<ResponseApi>(
    "/deploy/project",
    "POST",
    deployProjectDto
  );
}
