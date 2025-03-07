import { Loader2, Settings } from "lucide-react";
import { useState } from "react";

type ProjectPageHeaderProps = {
  onClickDeploy: () => Promise<void>;
  onClickSettings: () => void;
  toDeploy: boolean;
};

export default function ProjectPageHeader({
  onClickDeploy,
  onClickSettings,
  toDeploy,
}: ProjectPageHeaderProps) {
  const [isDeploying, setIsDeploying] = useState(false);

  const handleDeploy = async () => {
    if (isDeploying) return;

    setIsDeploying(true);
    try {
      await onClickDeploy();
    } finally {
      setIsDeploying(false);
    }
  };

  return (
    <div className="flex justify-between pt-6 pb-6 pl-10 pr-10">
      <div className="font-bold text-lg border rounded-full shadow-lg bg-white pl-4 pr-4 flex items-center">
        🛵 JustDeploy
      </div>
      <div className="border rounded-full shadow-lg bg-white p-2 flex gap-3">
        <button
          className="rounded-full hover:bg-green-50 w-9 h-9 flex justify-center items-center"
          onClick={onClickSettings}
        >
          <Settings />
        </button>
        <button
          disabled={!toDeploy || isDeploying}
          className={`font-mono hover:opacity-80 font-bold bg-button text-white p-1 pl-4 pr-4 rounded-xl border border-green-200 ${
            !toDeploy || isDeploying
              ? "opacity-50 cursor-not-allowed"
              : "border-animate"
          } flex items-center justify-center gap-2`}
          onClick={handleDeploy}
        >
          {isDeploying ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>Deploying...</span>
            </>
          ) : (
            "Deploy"
          )}
        </button>
      </div>
    </div>
  );
}
