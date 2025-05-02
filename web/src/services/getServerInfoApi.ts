import { callApi } from "./api";
import { ServerDto } from "./getServerListApi";

export type VersionDto = {
  tagName: string;
  githubUrl: string;
};

export type InfoDto = {
  version: VersionDto;
  firstConnection: boolean;
  server: ServerDto;
};

export async function getServerInfoApi(): Promise<InfoDto> {
  return await callApi<InfoDto>("/info", "GET");
}
