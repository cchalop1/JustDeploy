export const baseUrl = import.meta.env.VITE_BACKEND_URL;

// Function to get the API key from URL or localStorage
export function getApiKey(): string | null {
  // First check if API key is in URL
  const urlParams = new URLSearchParams(window.location.search);
  const apiKeyParam = urlParams.get("api_key");

  if (apiKeyParam) {
    // Save API key to localStorage
    localStorage.setItem("api_key", apiKeyParam);
    return apiKeyParam;
  }

  // Otherwise, try to get from localStorage
  return localStorage.getItem("api_key");
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

// Specific API to retrieve a new API key if needed
export async function fetchApiKey(): Promise<string> {
  const response = await callApi<{ api_key: string }>("/apikey", "GET");

  if (response.api_key) {
    localStorage.setItem("api_key", response.api_key);
    return response.api_key;
  }

  throw new Error("Failed to retrieve API key");
}
