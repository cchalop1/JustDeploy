import { create } from "zustand";
import { getServerInfoApi, InfoDto } from "@/services/getServerInfoApi";

interface InfoState {
  // Data
  isLoading: boolean;
  error: string | null;
  serverInfo: InfoDto | null;

  // Actions
  fetchServerInfo: () => Promise<void>;
  resetError: () => void;
}

export const useInfoStore = create<InfoState>((set, get) => ({
  // Initial state
  isLoading: false,
  error: null,
  serverInfo: null,

  // Computed values / getters
  get version() {
    return get().serverInfo?.version || null;
  },
  get isFirstConnection() {
    return get().serverInfo?.firstConnection || false;
  },
  get server() {
    return get().serverInfo?.server || null;
  },

  // Actions
  fetchServerInfo: async () => {
    try {
      set({ isLoading: true, error: null });
      const data = await getServerInfoApi();
      console.log(data);
      set({ serverInfo: data, isLoading: false });
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch server info";
      set({ error: errorMessage, isLoading: false });
      throw error;
    }
  },

  resetError: () => set({ error: null }),
}));

// Hook for accessing server info
export function useServerInfo() {
  return useInfoStore();
}
