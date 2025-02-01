import { callApi } from "./api";

export async function saveGithubAccessToken(
  installationId: string
): Promise<void> {
  await callApi<void>(`/github/save-access-token/${installationId}`, "POST");
}
