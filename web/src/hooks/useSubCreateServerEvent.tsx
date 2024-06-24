import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";

type EventCreateServer = {
  eventType: string;
  serverId: string;
  message: string;
  errorMessage?: string; // note the ? symbol for optional property
  step: number;
  totalSteps: number;
  time: Date;
};

export default function useSubCreateServerEvent() {
  const { id } = useParams<{ id: string }>();
  const events = useSubEvent<EventCreateServer>(`server/${id}/installation`);
  return events;
}
