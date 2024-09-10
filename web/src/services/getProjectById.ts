import { callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export type AppDto = {
  id: string;
  path: string;
  name: string;
  isDockerFile: boolean;
  isDockerCompose: boolean;
};

export type ProjectDto = {
  id: string;
  name: string;
  path: string;
  services: Service[];
  apps: AppDto[];
};

export async function getProjectByIdApi(id: string): Promise<ProjectDto> {
  return await callApi<ProjectDto>("/project/" + id, "GET");
}
