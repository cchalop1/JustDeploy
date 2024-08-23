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
  domain: string;
  createdDate: string;
  status: ServerStatusType;
};

export async function getServersListApi(): Promise<Array<ServerDto>> {
  return await callApi<Array<ServerDto>>("/server", "GET");
}
