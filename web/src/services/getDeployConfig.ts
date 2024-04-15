import { callApi } from "./api";
import { Env } from "./postFormDetails";

type SourceType = "Github" | "Local Folder";

export type DeployConfigDto = {
  sourceType: SourceType;
  pathToSource: string;
  dockerFileFound: boolean;
  composeFileFound: boolean;
  deployName: string;
  envFileFound: boolean;
  envs: Env[];
};

export async function getDeployConfig(
  deployId?: string,
): Promise<DeployConfigDto> {
  if (!deployId) {
    deployId = "";
  }
  return await callApi<DeployConfigDto>(`/deploy/config/${deployId}`, "GET");
}
