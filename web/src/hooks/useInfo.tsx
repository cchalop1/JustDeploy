import { useEffect } from "react";
import { useInfoStore } from "@/contexts/InfoStore";

/**
 * Hook to access server information
 * This hook will automatically fetch server info when mounted if it's not already loaded
 */
export function useInfo() {
  const { serverInfo, isLoading, error, fetchServerInfo, resetError } =
    useInfoStore();

  useEffect(() => {
    // Fetch server info if it's not already loaded
    if (!serverInfo && !isLoading && !error) {
      fetchServerInfo().catch(() => {
        // Error handling is already done in the store
      });
    }
  }, [serverInfo, isLoading, error, fetchServerInfo]);

  return {
    // Data
    serverInfo,
    isLoading,
    error,

    // Getters
    // Actions
    fetchServerInfo,
    resetError,
  };
}
