import { callApi } from "./api";

export type ServiceDto = {
  name: string;
  icon: string;
};

export async function getPreConfiguredServiceListApi(): Promise<
  Array<ServiceDto>
> {
  return await callApi<Array<ServiceDto>>("/service", "GET");
}
