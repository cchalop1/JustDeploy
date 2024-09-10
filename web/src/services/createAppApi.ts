import { callApi, ResponseApi } from "./api";

export type CreateAppApi = {
  path: string;
  projectId: string;
};

export async function createAppApi(body: CreateAppApi) {
  return await callApi<ResponseApi>(`/app`, "POST", body);
}
