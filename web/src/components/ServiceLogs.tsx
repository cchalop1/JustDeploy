import { useState } from "react";
import DisplayLogs from "./DisplayLogs";
import { useBuildLogs, useRunLogs } from "@/hooks/useServiceLogs";

type ServiceLogsProps = {
  serviceId: string;
};

type LogTab = "build" | "run";

export default function ServiceLogs({ serviceId }: ServiceLogsProps) {
  const [tab, setTab] = useState<LogTab>("build");

  const { logs: buildLogs, isLoading: isBuildLoading, error: buildError } = useBuildLogs(serviceId);
  const { logs: runLogs, isLoading: isRunLoading, error: runError } = useRunLogs(serviceId);

  const isLoading = tab === "build" ? isBuildLoading : isRunLoading;
  const error = tab === "build" ? buildError : runError;
  const logs = tab === "build" ? buildLogs : runLogs;

  return (
    <div className="flex flex-col gap-3">
      {/* Tab switcher */}
      <div className="flex gap-1 bg-gray-100 rounded-md p-0.5 w-fit">
        {(["build", "run"] as LogTab[]).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`px-3 py-1 text-xs font-medium rounded transition-colors ${
              tab === t
                ? "bg-white text-gray-900 shadow-sm"
                : "text-gray-500 hover:text-gray-700"
            }`}
          >
            {t === "build" ? "Build" : "Runtime"}
          </button>
        ))}
      </div>

      {/* Log content */}
      <div className="h-72 overflow-y-auto overflow-x-hidden rounded-md border border-gray-100 bg-gray-50">
        {isLoading ? (
          <div className="flex items-center justify-center h-full">
            <div className="w-5 h-5 border-2 border-gray-300 border-t-gray-600 rounded-full animate-spin" />
          </div>
        ) : error ? (
          <p className="font-mono text-xs text-red-500 px-3 py-2">
            Error: {error.message}
          </p>
        ) : (
          <DisplayLogs logs={logs} />
        )}
      </div>
    </div>
  );
}
