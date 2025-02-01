import { callApi } from "./api";

export type GithubRepo = {
  id: string;
  name: string;
  full_name: string;
  private: boolean;
};

export async function getGithubRepos(): Promise<Array<GithubRepo>> {
  return await callApi<Array<GithubRepo>>("/github/repos", "GET");
}
