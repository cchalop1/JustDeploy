import { ResponseApi, callApi } from "./api";

export type ConnectServerDto = {
  domain: string;
  sshKey: string | null;
  password: string | null;
  user: string;
};

export async function connectServer(
  connectServerDto: ConnectServerDto
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/connect", "POST", connectServerDto);
}
