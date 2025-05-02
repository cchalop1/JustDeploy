import { ResponseApi, callApi } from "./api";

type InitialSetupDto = {
  apiKey: string;
  domain: string;
};

/**
 * Save the initial setup with API key and domain
 * @param setupData The setup data containing API key and domain
 * @returns Promise with the API response
 */
export async function saveInitialSetup(
  setupData: InitialSetupDto
): Promise<ResponseApi> {
  return await callApi<ResponseApi>("/setup", "POST", setupData);
}
