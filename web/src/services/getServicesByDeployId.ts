import { callApi } from "./api";
import { Env } from "./postFormDetails";

export type ServiceExposeSettings = {
  isExposed: boolean;
  subDomain: string;
  exposePort: string;
  tls: boolean;
};

export type Service = {
  id: string;
  type: string;
  projectId: string;
  revision: number;
  branch: string;
  url: string;
  name: string;
  envs: Array<Env>;
  volumsNames: Array<string>;
  status: string;
  host: string;
  imageName: string;
  imageUrl: string;
  currentPath: string;
  exposePort: string;
  exposeSettings: ServiceExposeSettings;
};

export async function getServicesByDeployIdApi(
  deployId: string
): Promise<Service[]> {
  return await callApi<Service[]>("/deploy/" + deployId + "/service", "GET");
}
