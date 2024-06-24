import { baseUrl } from "./api";

export default function eventSubscription(path: string) {
  const source = new EventSource(`${baseUrl}/${path}`);
  return source;
}
