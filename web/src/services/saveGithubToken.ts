import { callApi } from "./api";

export async function saveGithubToken(
  githubAccessToken: string
): Promise<void> {
  return await callApi<void>("/github/token/", "POST", {
    token: githubAccessToken,
  });
}
