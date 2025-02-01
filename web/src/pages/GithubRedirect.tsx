import SpinnerIcon from "@/assets/SpinnerIcon";
import { finalizeGithubAppCreation } from "@/services/createGithubApp";
import { saveGithubAccessToken } from "@/services/saveGithubAccessToken";
import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";

export default function GithubRedirect() {
  const location = useLocation();
  const navigate = useNavigate();

  async function handleGithubAppCreation(code: string) {
    const res = await finalizeGithubAppCreation(code);
    console.log(res);

    const originalUrl = window.location.origin + "/github/auth/redirect";
    const stateValue = encodeURIComponent(originalUrl);
    const installUrl = `https://github.com/apps/${res.slug}/installations/new?state=${stateValue}`;

    console.log("Redirecting to:", installUrl);
    window.open(installUrl, "_self");
  }

  async function handleSaveAccessToken(installationId: string) {
    try {
      await saveGithubAccessToken(installationId);
      navigate("/");
    } catch (error) {
      console.error("Failed to save access token:", error);
    }
  }

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const code = params.get("code");
    const installationId = params.get("installation_id");

    if (code) {
      handleGithubAppCreation(code);
    }

    if (installationId) {
      console.log("Detected installation_id:", installationId);
      handleSaveAccessToken(installationId);
    }
  }, [location, navigate]);

  return (
    <div>
      <SpinnerIcon color="text-black"></SpinnerIcon>
    </div>
  );
}
