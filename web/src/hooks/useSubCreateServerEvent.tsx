import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";

export type EventServer = {
  eventType: string;
  title: string;
  errorMessage: string;
  time: Date;
};

export type EventServerWrapper = {
  serverName: string;
  serverId: string;
  eventServersList: EventServer[];
  currentStep: number;
};

export default function useSubCreateServerEvent() {
  const { id } = useParams<{ id: string }>();
  const event = useSubEvent<EventServerWrapper>(`server/${id}/loading`);
  return event;
}
