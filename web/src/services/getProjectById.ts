import { callApi } from "./api";
import { DeployDto } from "./getDeployListApi";
import { Service } from "./getServicesByDeployId";

export type ProjectDto = {
  id: string;
  name: string;
  path: string;
  services: Service[];
  deploys: DeployDto[];
};

export async function getProjectByIdApi(id: string): Promise<ProjectDto> {
  return await callApi<ProjectDto>("/project/" + id, "GET");
}
