import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";

type EventServer = {
  eventType: string;
  title: string;
  errorMessage: string;
  time: Date;
};

export type EventDeployWrapper = {
  deployName: string;
  deployId: string;
  eventsServer: EventServer[];
  currentStep: number;
};

export default function useSubCreateDeployEvent() {
  const { id } = useParams<{ id: string }>();
  const event = useSubEvent<EventDeployWrapper>(`deploy/${id}/loading`);
  return event;
}
