import eventSubscription from "@/services/eventSubsctiptions";
import { useEffect, useState } from "react";

export default function useSubEvent<T>(eventPath: string) {
  const [events, setSubEvent] = useState<Array<T>>([]);

  useEffect(() => {
    const source = eventSubscription(eventPath);
    source.onmessage = (e: MessageEvent<T>) => {
      setSubEvent((prev) => [...prev, e.data]);
    };

    source.onerror = (e) => {
      console.error("EventSource failed:", e);
      // TODO: find a better way to close the connection
      source.close();
    };

    return () => {
      source.close();
    };
  }, [eventPath]);

  return events;
}
