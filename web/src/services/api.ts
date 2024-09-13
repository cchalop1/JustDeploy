export const baseUrl = import.meta.env.VITE_BACKEND_URL;

export async function callApi<T>(
  path: string,
  method: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body?: any
): Promise<T> {
  const apiUrl = baseUrl ? baseUrl : window.location.origin; // Fallback to the current URL
  const fullUrl = new URL("api" + path, apiUrl).toString();
  console.log(fullUrl);
  const res = await fetch(fullUrl, {
    method,
    headers: {
      Accept: "application/json",
      "Content-type": "application/json",
    },
    body: JSON.stringify(body),
  });
  console.log(res);
  if (res.status >= 400) {
    const body = (await res.json()) as ResponseApi;
    throw new Error(body.message);
  }
  return (await res.json()) as T;
}

export type ResponseApi = {
  message: string;
};
