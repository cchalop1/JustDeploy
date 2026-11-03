import { callApi, baseUrl } from "./api";
const STATE_KEY = "github-state";
export async function newGithubState(): Promise<string> {
 const {state} = await callApi<{state:string}>("/github/state", "POST");
 sessionStorage.setItem(STATE_KEY,state); return state;
}
export function checkGithubState(state: string | null): string {
 const expected = sessionStorage.getItem(STATE_KEY);
 if (!state || !expected || state !== expected) throw new Error("Invalid GitHub callback. Start the connection again.");
 return state;
}
export function buildGithubAppManifest() {
 const origin = new URL(baseUrl || window.location.origin).origin;
 if (!origin.startsWith("https://")) throw new Error("Connect to JustDeploy using its public HTTPS URL before connecting GitHub.");
 return {
  name: `JustDeploy-${crypto.randomUUID().slice(0,8)}`, url:origin,
  hook_attributes:{url:`${origin}/api/github/events`},
  redirect_url:`${origin}/github/auth/redirect`,callback_urls:[`${origin}/github/auth/redirect`],setup_url:`${origin}/github/auth/redirect`,public:false,
  default_permissions:{contents:"read",metadata:"read"},default_events:["push"],
 };
}
export async function finalizeGithubAppCreation(code: string, state: string): Promise<{slug:string}> {
 return callApi<{slug:string}>(`/github/connect/${encodeURIComponent(code)}`, "POST", {state});
}
export async function redirectToGithubAppRegistration() {
 const manifest = buildGithubAppManifest();
 const state = await newGithubState();
 const url = new URL("https://github.com/settings/apps/new"); url.searchParams.set("state",state);
 const form = document.createElement("form"); form.method="POST"; form.action=url.toString();
 const input = document.createElement("input"); input.type="hidden"; input.name="manifest";input.value=JSON.stringify(manifest);
 form.appendChild(input); document.body.appendChild(form); form.submit(); form.remove();
}
