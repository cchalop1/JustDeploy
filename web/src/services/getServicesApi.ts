import { callApi } from "./api";

export type ServiceDto = {
  name: string;
  icon: string;
};

export async function getPreConfiguredServiceListApi(
  projectId?: string
): Promise<Array<ServiceDto>> {
  return await callApi<Array<ServiceDto>>("/service/" + projectId, "GET");
}
