import { ResponseApi, callApi } from "./api";

export function buildGithubAppManifest(serverIp: string) {
  const serverUrl = "http://" + serverIp + ":8080";
  const randomString = Math.random().toString(36).substring(2, 6);

  return {
    name: `JustDeploy-${randomString}`,
    url: serverUrl,
    hook_attributes: {
      url: serverUrl + "/github/events",
    },
    redirect_url: serverUrl + "/github/auth/redirect",
    callback_urls: [serverUrl + "/github/auth/redirect"],
    setup_url: serverUrl + "/github/auth/redirect",
    public: false,
    default_permissions: {
      contents: "read",
      metadata: "read",
    },
    default_events: ["push"],
  };
}

export function createGithubAppsUrl(): string {
  const baseUrl = "https://github.com/settings/apps/new";

  const url = new URL(baseUrl);
  url.searchParams.append("state", "1234");

  return url.toString();
}

export async function finalizeGithubAppCreation(code: string): Promise<any> {
  const path = `/github/connect/${code}`;
  return await callApi<ResponseApi>(path, "POST");
}

export function redirectToGithubAppRegistration(serverIp: string) {
  const form = document.createElement("form");
  form.method = "POST";
  form.action = createGithubAppsUrl();

  const input = document.createElement("input");
  input.type = "hidden";
  input.name = "manifest";
  const manifest = buildGithubAppManifest(serverIp);
  input.value = JSON.stringify(manifest);

  form.appendChild(input);
  document.body.appendChild(form);
  form.submit();
}
