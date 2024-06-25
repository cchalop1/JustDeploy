import { useEffect, useState } from "react";
import eventSubscription from "@/services/eventSubsctiptions";

export default function useSubEvent<T>(eventPath: string) {
  const [events, setSubEvent] = useState<T | null>(null);

  useEffect(() => {
    const source = eventSubscription(eventPath);
    source.onmessage = (e: MessageEvent<T>) => {
      console.log(e);
      setSubEvent(JSON.parse<T>(e.data));
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
