import { callApi } from "./api";

export async function saveGithubAccessToken(
  installationId: string, state: string
): Promise<void> {
  await callApi<void>(`/github/save-access-token/${encodeURIComponent(installationId)}`, "POST", {state});
}
