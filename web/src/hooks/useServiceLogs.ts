import { useState, useEffect } from "react";
import { getServiceBuildLogs } from "@/services/getServiceBuildLogs";
import { getServiceRunLogs } from "@/services/getServiceRunLogs";
import { Logs } from "@/services/getApplicationLogs";

export function useBuildLogs(serviceId: string) {
  const [logs, setLogs] = useState<Logs[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchLogs = async () => {
      setIsLoading(true);
      try {
        const data = await getServiceBuildLogs(serviceId);
        setLogs(data);
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Failed to fetch build logs")
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
  }, [serviceId]);

  return { logs, isLoading, error };
}

export function useRunLogs(serviceId: string) {
  const [logs, setLogs] = useState<Logs[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchLogs = async () => {
      setIsLoading(true);
      try {
        const data = await getServiceRunLogs(serviceId);
        setLogs(data);
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err : new Error("Failed to fetch run logs")
        );
      } finally {
        setIsLoading(false);
      }
    };

    fetchLogs();
  }, [serviceId]);

  return { logs, isLoading, error };
}
