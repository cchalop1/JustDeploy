import { ResponseApi, callApi } from "./api";

export type ConnectServerDto = {
  ip: string;
  user: string;
  sshKey: string | null;
  password: string | null;
};

export async function connectServerApi(
  connectServerDto: ConnectServerDto
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/server", "POST", connectServerDto);
}
