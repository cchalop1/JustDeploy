import { callApi } from "./api";

export type ServiceDto = {
  name: string;
  image: string;
};

export async function getServiceListApi(): Promise<Array<ServiceDto>> {
  return await callApi<Array<ServiceDto>>("/service", "GET");
}
