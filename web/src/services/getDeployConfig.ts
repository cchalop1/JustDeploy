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

type GetDeployConfigParams = {
  deployId?: string;
  path?: string;
};

export async function getDeployConfig(
  params: GetDeployConfigParams
): Promise<DeployConfigDto> {
  if (!params.deployId) {
    params.deployId = "";
  }

  return await callApi<DeployConfigDto>(
    `/deploy/config/${params.deployId}`,
    "POST",
    params
  );
}
