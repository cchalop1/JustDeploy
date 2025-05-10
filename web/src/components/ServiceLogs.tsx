import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import DisplayLogs from "./DisplayLogs";
import { useBuildLogs, useRunLogs } from "@/hooks/useServiceLogs";

type ServiceLogsProps = {
  serviceId: string;
};

export default function ServiceLogs({ serviceId }: ServiceLogsProps) {
  const {
    logs: buildLogs,
    isLoading: isBuildLoading,
    error: buildError,
  } = useBuildLogs(serviceId);
  const {
    logs: runLogs,
    isLoading: isRunLoading,
    error: runError,
  } = useRunLogs(serviceId);

  return (
    <div className="w-full max-w-full overflow-hidden">
      <Tabs defaultValue="build" className="w-full">
        <TabsList className="w-full">
          <TabsTrigger value="build" className="flex-1">
            Build Logs
          </TabsTrigger>
          <TabsTrigger value="run" className="flex-1">
            Run Logs
          </TabsTrigger>
        </TabsList>
        <TabsContent
          value="build"
          className="h-64 overflow-auto mt-4 w-full overflow-x-auto"
        >
          {isBuildLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
          ) : buildError ? (
            <div className="text-red-500 p-4">Error: {buildError.message}</div>
          ) : (
            <div className="w-full whitespace-pre overflow-x-auto">
              <DisplayLogs logs={buildLogs} />
            </div>
          )}
        </TabsContent>
        <TabsContent
          value="run"
          className="h-64 overflow-auto mt-4 w-full overflow-x-auto"
        >
          {isRunLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
          ) : runError ? (
            <div className="text-red-500 p-4">Error: {runError.message}</div>
          ) : (
            <div className="w-full whitespace-pre overflow-x-auto">
              <DisplayLogs logs={runLogs} />
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
