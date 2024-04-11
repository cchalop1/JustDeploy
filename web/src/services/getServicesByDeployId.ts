import { callApi } from "./api";
import { Env } from "./postFormDetails";

export type Service = {
  id: string;
  name: string;
  status: string;
  envs: Array<Env>;
  volumsNames: Array<string>;
  host: string;
  deployId: string;
  imageName: string;
};

export async function getServicesByDeployIdApi(
  deployId: string
): Promise<Service[]> {
  return await callApi<Service[]>("/deploy/" + deployId + "/service", "GET");
}
