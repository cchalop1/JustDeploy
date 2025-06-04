import { useEffect, useState } from "react";
import {
  getServiceCommitInfo,
  type ServiceCommitInfo,
} from "@/services/getServiceCommitInfo";

type ServiceCommitInfoProps = {
  serviceId: string;
};

export default function ServiceCommitInfo({
  serviceId,
}: ServiceCommitInfoProps) {
  const [commitInfo, setCommitInfo] = useState<ServiceCommitInfo | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchCommitInfo() {
      try {
        setLoading(true);
        const info = await getServiceCommitInfo(serviceId);
        setCommitInfo(info);
      } catch (err) {
        console.error("Error fetching commit info:", err);
      } finally {
        setLoading(false);
      }
    }

    fetchCommitInfo();
  }, [serviceId]);

  if (loading) {
    return <div className="text-sm text-muted-foreground">Chargement...</div>;
  }

  if (!commitInfo) {
    return null; // Ne rien afficher pour les services non-GitHub
  }

  return (
    <div className="text-sm">
      <div className="text-muted-foreground mb-2">Source</div>
      <div className="flex items-center gap-2 mb-1">
        <span>main</span>
      </div>
      <div className="flex items-center gap-2">
        <span className="text-muted-foreground">⚬</span>
        <span className="text-muted-foreground">{commitInfo.hash}</span>
        <a
          href={commitInfo.githubUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="text-foreground hover:text-blue-600 hover:underline cursor-pointer"
        >
          {commitInfo.message}
        </a>
      </div>
    </div>
  );
}
