import { callApi } from "./api";

type LoginDto = {
  email: string;
  password: string;
};

type AuthResponseDto = {
  token: string;
};

export async function loginApi(data: LoginDto): Promise<AuthResponseDto> {
  return await callApi<AuthResponseDto>("/login", "POST", data);
}
