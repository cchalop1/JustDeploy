import { GithubIcon } from "lucide-react";
import { Button } from "./ui/button";
import {
  buildGithubAppManifest,
  createGithubAppsUrl,
} from "@/services/createGithubApp";

type BtnConnectGithubProps = {
  serverIp: string;
};

export default function BtnConnectGithub({ serverIp }: BtnConnectGithubProps) {
  function redirectToGithubAppRegistration() {
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

  return (
    <Button
      asChild
      className="w-full cursor-pointer font-bold mt-3"
      onClick={redirectToGithubAppRegistration}
    >
      <div>
        <GithubIcon className="mr-2 h-5 w-5" />
        Connect your GitHub
      </div>
    </Button>
  );
}
