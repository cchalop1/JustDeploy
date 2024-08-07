import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import EventList from "@/components/EventsList";
import { useEffect } from "react";
import useSubCreateDeployEvent from "@/hooks/useSubCreateDeployEvent";

// TODO: rename to DeployLoadingPage
export default function CreateDeployLoading() {
  const navigate = useNavigate();
  const event = useSubCreateDeployEvent();
  const isLoading =
    !event || event?.currentStep !== event?.eventsDeployList.length;

  useEffect(() => {
    if (
      event &&
      event.currentStep === event.eventsDeployList.length &&
      !event.eventsDeployList.find((e) => e.errorMessage)
    ) {
      navigate("/deploy/" + event.deployId);
    }
  }, [event, navigate]);

  return (
    <div>
      <div className="mt-10 text-2xl font-bold ">
        Creating {event?.deployName} ...
      </div>
      {event && (
        <EventList
          eventList={event.eventsDeployList}
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
