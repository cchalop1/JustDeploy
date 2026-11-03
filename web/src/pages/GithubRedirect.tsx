import { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { checkGithubState, finalizeGithubAppCreation, newGithubState } from "@/services/createGithubApp";
import { saveGithubAccessToken } from "@/services/saveGithubAccessToken";
import SpinnerIcon from "@/assets/SpinnerIcon";
export default function GithubRedirect() {
 const location = useLocation(); const navigate = useNavigate();
 const started = useRef(false); const [error,setError] = useState("");
 useEffect(() => {
  if (started.current) return;
  started.current = true;
  async function complete() {
   const params = new URLSearchParams(location.search);
   const state = checkGithubState(params.get("state"));
   const code = params.get("code"); const installation = params.get("installation_id");
   // Remove one-time codes from history before making network requests.
   window.history.replaceState(null,"",location.pathname);
   if (code) {
    const {slug} = await finalizeGithubAppCreation(code,state);
    const installState = await newGithubState();
    const url = new URL(`https://github.com/apps/${encodeURIComponent(slug)}/installations/new`);
    url.searchParams.set("state",installState); window.location.assign(url.toString());
   } else if (installation) {
    await saveGithubAccessToken(installation,state); sessionStorage.removeItem("github-state"); navigate("/",{replace:true});
   } else { throw new Error("GitHub did not return a code or installation."); }
  }
  complete().catch((e:unknown) => setError(e instanceof Error ? e.message : "GitHub connection failed"));
 },[location.search,location.pathname,navigate]);
 return error ? <div role="alert" className="p-8">{error} <a href="/">Return home</a></div> : <SpinnerIcon color="text-black" />;
}
