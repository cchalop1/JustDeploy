import { callApi } from "./api";
import { ConnectServerDto } from "./connectServer";
import { PostCreateDeploymentDto } from "./postFormDetails";

export type deployStatus = "serverconfig" | "appconfig" | "deployapp";

export type GetDeployConfigResponse = {
  pathToProject: string;
  dockerFileValid: boolean;
  serverConfig: ConnectServerDto | null;
  appConfig: PostCreateDeploymentDto | null;
  deployStatus: deployStatus;
  url: string;
};

export async function getDeployConfig(): Promise<GetDeployConfigResponse> {
  return await callApi<GetDeployConfigResponse>("/deploy", "GET");
}
