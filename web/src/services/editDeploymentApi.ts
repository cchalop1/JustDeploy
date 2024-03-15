import { ResponseApi, callApi } from "./api";

export type EditDeployementDto = {
  deployOnCommit: boolean;
};

export async function editDeployementApi(
  editDeployementDto: EditDeployementDto,
  appName: string
) {
  return await callApi<ResponseApi>(
    "/deploy/" + appName,
    "PUT",
    editDeployementDto
  );
}
