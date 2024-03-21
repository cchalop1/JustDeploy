import { callApi } from "./api";
import { Env } from "./postFormDetails";

export type DeployStatusType = "Runing" | "Installing";

export type DeployDto = {
  id: string;
  name: string;
  serverId: string;
  url: string;
  enableTls: boolean;
  email: string;
  pathToSource: string;
  envs: Env[];
  deployOnCommit: boolean;
  status: DeployStatusType;
};

export async function getDeployListApi(): Promise<Array<DeployDto>> {
  return await callApi<Array<DeployDto>>("/deploy", "GET");
}
