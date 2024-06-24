import eventSubscription from "@/services/eventSubsctiptions";
import { useEffect, useState } from "react";

export default function useSubEvent<T>(eventPath: string) {
  const [events, setSubEvent] = useState<Array<T>>([]);

  useEffect(() => {
    const source = eventSubscription(eventPath);
    source.onmessage = (e: MessageEvent<T>) => {
      setSubEvent((prev) => [...prev, e.data]);
    };

    return () => {
      source.close();
    };
  }, []);

  return events;
}
