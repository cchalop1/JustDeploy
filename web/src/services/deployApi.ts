import { callApi, ResponseApi } from "./api";

export async function deployApi() {
  return await callApi<ResponseApi>("/deploy", "POST");
}
