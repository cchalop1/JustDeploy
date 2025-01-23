export function buildGithubAppManifest(serverIp: string) {
  const serverUrl = "http://" + serverIp + ":8080";
  return {
    name: "JustDeploy-1234",
    url: serverUrl,
    hook_attributes: {
      url: serverUrl + "/github/events",
    },
    redirect_url: serverUrl + "/redirect",
    callback_urls: [serverUrl + "/callback"],
    public: true,
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
