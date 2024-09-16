import { FileText, Settings } from "lucide-react";

export default function ProjectPageHeader() {
  return (
    <div className="flex justify-between pt-6 pb-6 pl-10 pr-10">
      <div className="font-semibold border rounded-full shadow-lg bg-white pl-4 pr-4 flex items-center">
        JustDeploy
      </div>
      <div className="border rounded-full shadow-lg bg-white p-2 flex gap-3">
        <button className="rounded-full bg-green-50 w-9 h-9 flex justify-center items-center">
          <Settings />
        </button>
        <button className="rounded-full bg-green-50 w-9 h-9 flex justify-center items-center">
          <FileText />
        </button>
        <button className="font-mono bg-button text-white p-1 pl-4 pr-4 rounded-xl border border-green-200">
          Deploy
        </button>
      </div>
    </div>
  );
}
