import { callApi } from "./api";
import { ConnectServerDto } from "./connectServer";
import { PostCreateDeploymentDto } from "./postFormDetails";

export type DeployFormStatus = "serverconfig" | "appconfig" | "deployapp";

export type GetDeployConfigResponse = {
  pathToProject: string;
  dockerFileValid: boolean;
  serverConfig: ConnectServerDto | null;
  appConfig: PostCreateDeploymentDto | null;
  deployFormStatus: DeployFormStatus;
  url: string;
};

export async function getDeployConfig(): Promise<GetDeployConfigResponse> {
  return await callApi<GetDeployConfigResponse>("/deploy", "GET");
}
