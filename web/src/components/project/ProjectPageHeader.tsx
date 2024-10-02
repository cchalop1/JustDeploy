import { FileText, Settings } from "lucide-react";
import { Tooltip, TooltipProvider, TooltipTrigger } from "../ui/tooltip";
import { TooltipContent } from "@radix-ui/react-tooltip";

type ProjectPageHeaderProps = {
  onClickDeploy: () => void;
  onClickSettings: () => void;
};

export default function ProjectPageHeader({
  onClickDeploy,
  onClickSettings,
}: ProjectPageHeaderProps) {
  // TODO: add button to start in local and deploy to server
  return (
    <div className="flex justify-between pt-6 pb-6 pl-10 pr-10">
      <div className="font-semibold border rounded-full shadow-lg bg-white pl-4 pr-4 flex items-center">
        🛵 JustDeploy
      </div>
      <TooltipProvider>
        <div className="border rounded-full shadow-lg bg-white p-2 flex gap-3">
          <Tooltip>
            <TooltipTrigger>
              <button
                className="rounded-full bg-green-50 w-9 h-9 flex justify-center items-center"
                onClick={onClickSettings}
              >
                <Settings />
              </button>
            </TooltipTrigger>
            <TooltipContent>
              <div className="p-2">Settings</div>
            </TooltipContent>
          </Tooltip>
          <button className="rounded-full bg-green-50 w-9 h-9 flex justify-center items-center">
            <FileText />
          </button>
          <button
            className="font-mono bg-button text-white p-1 pl-4 pr-4 rounded-xl border border-green-200"
            onClick={onClickDeploy}
          >
            Deploy
          </button>
        </div>
      </TooltipProvider>
    </div>
  );
}
