import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import EventList from "@/components/EventsList";
import { useEffect } from "react";
import useSubCreateDeployEvent from "@/hooks/useSubCreateDeployEvent";

export default function CreateDeployLoading() {
  const navigate = useNavigate();
  const event = useSubCreateDeployEvent();
  const isLoading = !event || event?.currentStep !== event?.eventsServer.length;

  useEffect(() => {
    if (
      event &&
      event.currentStep === event.eventsServer.length &&
      !event.eventsServer.find((e) => e.errorMessage)
    ) {
      navigate("/");
    }
  }, [event, navigate]);

  return (
    <div>
      <div className="mt-10 text-2xl font-bold ">
        Creating {event?.deployName} ...
      </div>
      {event && (
        <EventList
          eventList={event.eventsServer}
          currentStep={event.currentStep}
        />
      )}
      <div className="flex justify-center m-4">
        <Button disabled={isLoading} onClick={() => navigate("/")}>
          Click to go to your application
        </Button>
      </div>
    </div>
  );
}
