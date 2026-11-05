import { callApi } from "./api";

type InitialSetupDto = {
  email: string;
  password: string;
  domain: string;
};

type AuthResponseDto = {
  token: string;
};

export async function saveInitialSetup(
  setupData: InitialSetupDto
): Promise<AuthResponseDto> {
  return await callApi<AuthResponseDto>("/setup", "POST", setupData);
}
