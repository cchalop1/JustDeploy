import { Settings } from "lucide-react";

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
          disabled={!toDeploy}
          className={`font-mono hover:opacity-80 font-bold bg-button text-white p-1 pl-4 pr-4 rounded-xl border border-green-200 ${
            !toDeploy ? "opacity-50 cursor-not-allowed" : "border-animate"
          }`}
          onClick={async () => {
            await onClickDeploy();
          }}
        >
          Deploy
        </button>
      </div>
    </div>
  );
}
