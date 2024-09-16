import { callApi } from "./api";

type VersionDto = {
  tagName: string;
  githubUrl: string;
};

export async function getVersion(): Promise<VersionDto> {
  return await callApi<VersionDto>("/version", "GET");
}
