import { useParams } from "react-router-dom";
import useSubEvent from "./useSubEvent";

export default function useSubCreateServerEvent() {
  const { id } = useParams<{ id: string }>();
  const events = useSubEvent<string>(`server/${id}/installation`);
  return events;
}
