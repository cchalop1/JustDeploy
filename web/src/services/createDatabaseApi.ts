import { callApi } from "./api";
import { ServiceDto } from "./getServicesApi";

export type CreateDatabaseRequest = {
  databaseName: string;
};

export async function createDatabaseApi(
  body: CreateDatabaseRequest
): Promise<ServiceDto> {
  return await callApi<ServiceDto>("/database/create", "POST", body);
}
