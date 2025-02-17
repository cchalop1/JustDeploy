import { ResponseApi, callApi } from "./api";

type AddDomainToServerDto = {
  domain: string;
};

export async function addDomainToServerApi(body: AddDomainToServerDto) {
  return await callApi<ResponseApi>(`/server/domain`, "POST", body);
}
