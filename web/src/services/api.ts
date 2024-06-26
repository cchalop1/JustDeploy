export const baseUrl = "http://localhost:8080/api";

export async function callApi<T>(
  path: string,
  method: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body?: any
): Promise<T> {
  const res = await fetch(baseUrl + path, {
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
