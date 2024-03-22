import { ResponseApi, callApi } from "./api";
import { Env } from "./postFormDetails";

export type EditDeployDto = {
  deployOnCommit: boolean;
  id: string;
  envs: Array<Env>;
  subDomain: string;
};

export async function editDeployementApi(editDeployDto: EditDeployDto) {
  return await callApi<ResponseApi>("/deploy/edit", "PUT", editDeployDto);
}
