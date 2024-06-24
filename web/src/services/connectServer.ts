import { callApi } from "./api";
import { ServerDto } from "./getServerListApi";

export type ConnectServerDto = {
  ip: string;
  user: string;
  sshKey: string | null;
  password: string | null;
};

export async function connectServerApi(
  connectServerDto: ConnectServerDto
): Promise<ServerDto> {
  return await callApi<ServerDto>("/server", "POST", connectServerDto);
}
