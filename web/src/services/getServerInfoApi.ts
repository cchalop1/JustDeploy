import { callApi } from "./api";
import { ServerDto } from "./getServerListApi";

export async function getServerInfoApi(): Promise<ServerDto> {
  return await callApi<ServerDto>("/server/info", "GET");
}
