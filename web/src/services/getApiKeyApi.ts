import { callApi } from "./api";

export type ApiKeyDto = {
  apiKey: string;
};

export async function getApiKeyApi(): Promise<ApiKeyDto> {
  return await callApi<ApiKeyDto>("/api-key", "GET");
}
