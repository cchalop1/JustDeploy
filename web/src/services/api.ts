import { getStoredToken, removeToken } from "./authStorage";
export const baseUrl = import.meta.env.VITE_BACKEND_URL;
export async function callApi<T>(path: string, method: string, body?: unknown): Promise<T> {
 const url = new URL(`/api${path}`, baseUrl || window.location.origin);
 const token = getStoredToken();
 const response = await fetch(url, {
  method,
  headers: { Accept: "application/json", "Content-Type": "application/json", ...(token ? {Authorization: `Bearer ${token}`} : {}) },
  body: body === undefined ? undefined : JSON.stringify(body),
 });
 const text = await response.text();
 let data: unknown;
 try { data = text ? JSON.parse(text) : undefined; } catch { data = text; }
 if (!response.ok) {
  if (response.status === 401) { removeToken(); window.dispatchEvent(new Event("auth-changed")); }
  let message = `Request failed (${response.status})`;
  if (typeof data === "string" && data.length < 500) message = data;
  else if (data && typeof data === "object") {
   const error = data as {message?: unknown; error?: unknown};
   if (typeof error.message === "string") message = error.message;
   else if (typeof error.error === "string") message = error.error;
  }
  throw new Error(message);
 }
 return data as T;
}
export type ResponseApi = { message: string };
