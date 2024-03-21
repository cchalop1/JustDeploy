import { callApi } from "./api";
import { ConnectServerDto } from "./connectServer";
import { CreateDeployDto } from "./postFormDetails";

export type deployStatus = "serverconfig" | "appconfig" | "deployapp";
export type appStatus = "Stoped" | "Runing";

export type GetDeployConfigResponse = {
  dockerFileValid: boolean;
  serverConfig: ConnectServerDto | null;
  appConfig: CreateDeployDto | null;
  deployStatus: deployStatus;
  appStatus: appStatus;
  url: string;
};

export async function getDeployConfig(): Promise<GetDeployConfigResponse> {
  return await callApi<GetDeployConfigResponse>("/deploy", "GET");
}
