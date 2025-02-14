import { GithubIcon } from "lucide-react";

export function CardIcon({ service }: { service: Service }) {
  if (service.type === "github_repo") {
    return <GithubIcon size={24} />;
  }
  return <img src="/icons/service.png" className="h-6" />;
}
