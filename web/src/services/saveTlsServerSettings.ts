import { ResponseApi, callApi } from "./api";

type ServerTlsSettings = {
  useHttps: boolean;
  email: string;
};

export async function saveTlsServerSettings(body: ServerTlsSettings) {
  return await callApi<ResponseApi>(`/server/tls-settings`, "PUT", body);
}
