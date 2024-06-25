import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";

type EventServer = {
  eventType: string;
  title: string;
  errorMessage: string;
  time: Date;
};

type EventServerWrapper = {
  serverName: string;
  serverId: string;
  eventsServer: EventServer[];
  currentStep: number;
};

export default function useSubCreateServerEvent() {
  const { id } = useParams<{ id: string }>();
  const event = useSubEvent<EventServerWrapper>(`server/${id}/installation`);
  return event;
}
