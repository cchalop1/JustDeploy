import useSubCreateServerEvent from "@/hooks/useSubCreateServerEvent";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import EventList from "@/components/EventsList";

export default function CreateServerLoading() {
  const navigate = useNavigate();
  const event = useSubCreateServerEvent();
  const isLoading = !event || event?.currentStep !== event?.eventsServer.length;

  return (
    <div>
      <div className="mt-10 text-2xl font-bold ">
        Creating {event?.serverName} ...
      </div>
      {event && (
        <EventList
          eventList={event.eventsServer}
          currentStep={event.currentStep}
        />
      )}
      <div className="flex justify-center m-4">
        <Button disabled={isLoading} onClick={() => navigate("/")}>
          Deploy application
        </Button>
      </div>
    </div>
  );
}
