import { callApi } from "./api";
import { Service } from "./getServicesByDeployId";

export type CreateDatabaseRequest = {
  databaseName: string;
};

export async function createDatabaseApi(
  body: CreateDatabaseRequest
): Promise<Service> {
  return await callApi<Service>("/database/create", "POST", body);
}
