import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";
import { EventServer } from "./useSubCreateServerEvent";

export type EventDeployWrapper = {
  deployName: string;
  deployId: string;
  eventsDeployList: EventServer[];
  currentStep: number;
};

export default function useSubCreateDeployEvent() {
  const { id } = useParams<{ id: string }>();
  const event = useSubEvent<EventDeployWrapper>(`deploy/${id}/loading`);
  return event;
}
