import { callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export type ServiceDto = {
  name: string;
  type: "llm" | "database";
  icon: string;
};

export async function getPreConfiguredServiceListApi(
  projectId?: string
): Promise<Array<ServiceDto>> {
  return await callApi<Array<ServiceDto>>("/service/" + projectId, "GET");
}

export async function getServicesApi() {
  return await callApi<Array<Service>>("/services", "GET");
}
