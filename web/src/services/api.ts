export const baseUrl = import.meta.env.VITE_BACKEND_URL;

console.log(import.meta.env);
console.log(baseUrl);

import { getStoredApiKey, saveApiKey } from "./authStorage";

// Function to get the API key from URL or localStorage
export function getApiKey(): string | null {
  // First check if API key is in URL
  const urlParams = new URLSearchParams(window.location.search);
  const apiKeyParam = urlParams.get("api_key");

  if (apiKeyParam) {
    // Save API key to localStorage using authStorage function
    saveApiKey(apiKeyParam);
    return apiKeyParam;
  }

  // Otherwise, try to get from localStorage using authStorage function
  return getStoredApiKey();
}

export async function callApi<T>(
  path: string,
  method: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body?: any
): Promise<T> {
  const apiUrl = baseUrl ? baseUrl : window.location.origin; // Fallback to the current URL
  const fullUrl = new URL("api" + path, apiUrl).toString();

  // Get API key
  const apiKey = getApiKey();

  // Prepare headers
  const headers: HeadersInit = {
    Accept: "application/json",
    "Content-type": "application/json",
  };

  // Add API key header if available
  if (apiKey) {
    headers["X-API-Key"] = apiKey;
  }

  const res = await fetch(fullUrl, {
    method,
    headers,
    body: JSON.stringify(body),
  });

  if (res.status >= 400) {
    if (res.status === 401) {
      // Handle unauthorized access
      console.error("Unauthorized access. API key may be invalid.");
    }
    const body = (await res.json()) as ResponseApi;
    throw new Error(body.message);
  }
  return (await res.json()) as T;
}

export type ResponseApi = {
  message: string;
};
