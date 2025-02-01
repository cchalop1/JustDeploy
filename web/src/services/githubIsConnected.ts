import { callApi } from "./api";

type GithubIsConnectedResponse = {
  isConnected: boolean;
};

export async function githubIsConnectedApi(): Promise<GithubIsConnectedResponse> {
  return await callApi<GithubIsConnectedResponse>(
    "/github/is-connected",
    "GET"
  );
}
