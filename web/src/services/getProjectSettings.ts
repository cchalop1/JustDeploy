import { callApi } from "./api";

export type PathDto = {
  name: string;
  fullPath: string;
  folders: PathDto[];
};

export type ProjectSettingsDto = {
  currentPath: string;
  folders: PathDto[];
};

export async function getProjectSettingsByIdApi(
  id: string
): Promise<ProjectSettingsDto> {
  return await callApi<ProjectSettingsDto>(
    "/project/" + id + "/settings",
    "GET"
  );
}
