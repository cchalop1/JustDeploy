import { callApi } from "./api";

export type VersionDto = {
  tagName: string;
  githubUrl: string;
};

export async function getVersionApi(): Promise<VersionDto> {
  return await callApi<VersionDto>("/version", "GET");
}
