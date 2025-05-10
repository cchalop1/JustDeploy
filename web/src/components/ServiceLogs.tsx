import { useEffect, useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import DisplayLogs from "./DisplayLogs";
import { getServiceBuildLogs } from "@/services/getServiceBuildLogs";
import { getServiceRunLogs } from "@/services/getServiceRunLogs";
import { Logs } from "@/services/getApplicationLogs";

type ServiceLogsProps = {
  serviceId: string;
};

export default function ServiceLogs({ serviceId }: ServiceLogsProps) {
  const [buildLogs, setBuildLogs] = useState<Logs[]>([]);
  const [runLogs, setRunLogs] = useState<Logs[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchLogs = async () => {
      setIsLoading(true);
      try {
        const [buildLogsData, runLogsData] = await Promise.all([
          getServiceBuildLogs(serviceId),
          getServiceRunLogs(serviceId),
        ]);
        setBuildLogs(buildLogsData);
        setRunLogs(runLogsData);
      } catch (error) {
        console.error("Error fetching logs:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
  }, [serviceId]);

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
          {isLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
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
          {isLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
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
