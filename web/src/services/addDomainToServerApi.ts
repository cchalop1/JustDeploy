import { ResponseApi, callApi } from "./api";

type AddDomainToServerDto = {
  domain: string;
};

export async function addDomainToServerApi(
  serverId: string,
  body: AddDomainToServerDto
) {
  return await callApi<ResponseApi>(`/server/${serverId}/domain`, "POST", body);
}
