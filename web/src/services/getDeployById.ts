import { callApi } from "./api";
import { DeployDto } from "./getDeployListApi";

export async function getDeployByIdApi(id: string): Promise<DeployDto> {
  return await callApi<DeployDto>("/deploy/" + id, "GET");
}
