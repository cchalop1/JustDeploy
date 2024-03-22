import { ResponseApi, callApi } from "./api";

export type EditDeployDto = {
  deployOnCommit: boolean;
  id: string;
};

export async function editDeployementApi(editDeployDto: EditDeployDto) {
  return await callApi<ResponseApi>("/deploy/edit", "PUT", editDeployDto);
}
