import { callApi } from "./api";

export type ServerStatusType =
  | "Runing"
  | "NotConnected"
  | "Installing"
  | "Error";

export type ServerDto = {
  id: string;
  ip: string;
  name: string;
  port: string;
  domain: string;
  createdDate: string;
  status: ServerStatusType;
  useHttps?: boolean;
  email?: string;
};

export async function getServersListApi(): Promise<Array<ServerDto>> {
  return await callApi<Array<ServerDto>>("/server", "GET");
}
