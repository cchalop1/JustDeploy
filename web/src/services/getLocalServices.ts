import { callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export async function getLocalServices(): Promise<Array<Service>> {
  return await callApi<Array<Service>>("/local/services", "GET");
}
