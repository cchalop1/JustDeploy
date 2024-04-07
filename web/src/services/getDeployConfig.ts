import { callApi } from "./api";
import { Env } from "./postFormDetails";

type SourceType = "Github" | "Local Folder";

export type DeployConfigDto = {
  sourceType: SourceType;
  pathToSource: string;
  dockerFileFound: boolean;
  composeFileFound: boolean;
  deployName: string;
  envs: Env[];
};

export async function getDeployConfig(): Promise<DeployConfigDto> {
  return await callApi<DeployConfigDto>("/deploy/config", "GET");
}
