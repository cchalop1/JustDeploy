import { callApi } from "./api";
import { ServerDto } from "./getServerListApi";

export async function getServerByIdApi(id: string): Promise<ServerDto> {
  return await callApi<ServerDto>("/server/" + id, "GET");
}
