export const baseUrl = import.meta.env.VITE_BACKEND_URL;

import { getStoredToken } from "./authStorage";

export async function callApi<T>(
  path: string,
  method: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body?: any
): Promise<T> {
  const apiUrl = baseUrl ? baseUrl : window.location.origin;
  const fullUrl = new URL("api" + path, apiUrl).toString();

  const token = getStoredToken();

  const headers: HeadersInit = {
    Accept: "application/json",
    "Content-type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(fullUrl, {
    method,
    headers,
    body: JSON.stringify(body),
  });

  if (res.status >= 400) {
    if (res.status === 401) {
      console.error("Unauthorized. Token may be invalid or expired.");
    }
    const body = (await res.json()) as ResponseApi;
    throw new Error(body.message);
  }
  return (await res.json()) as T;
}

export type ResponseApi = {
  message: string;
};
