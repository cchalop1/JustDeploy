import { callApi } from "./api";
import { DeployDto } from "./getDeployListApi";

export async function getDeployListByServerIdApi(
  serverId: string
): Promise<Array<DeployDto>> {
  return await callApi<Array<DeployDto>>(
    "/server/" + serverId + "/deploy",
    "GET"
  );
}
