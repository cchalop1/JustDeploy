import { callApi } from "./api";
import { ServerStatusType } from "./getServerListApi";
import { Env } from "./postFormDetails";

export type Service = {
  id: string;
  hostName: string;
  status: ServerStatusType;
  envs: Array<Env>;
  volumsNames: Array<string>;
  host: string;
  deployId: string;
  imageName: string;
  imageUrl: string;
  isDevContainer: boolean;
  currentPath: string;
  exposePort: string;
};

export async function getServicesByDeployIdApi(
  deployId: string
): Promise<Service[]> {
  return await callApi<Service[]>("/deploy/" + deployId + "/service", "GET");
}
