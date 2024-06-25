import useSubCreateServerEvent from "@/hooks/useSubCreateServerEvent";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";

export default function CreateServerStatus() {
  const navigate = useNavigate();
  const event = useSubCreateServerEvent();
  const isLoading = !event || event?.currentStep !== event?.eventsServer.length;

  return (
    <div>
      <div className="mt-10 text-2xl font-bold ">
        Creating {event?.serverName} ...
      </div>
      {event && (
        <Accordion type="single" collapsible className="w-full">
          {event.eventsServer.map((eventServer, index) => (
            <AccordionItem key={index} value={`item-${index}`}>
              <AccordionTrigger>
                <div className="flex gap-5 items-center">
                  <div
                    className={`${
                      eventServer.errorMessage ? "bg-red-400" : ""
                    }`}
                  >
                    {eventServer.title}
                  </div>
                  {event.currentStep === index + 1 ? (
                    <SpinnerIcon color="text-black" />
                  ) : null}
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <code>{eventServer.errorMessage}</code>
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      )}
      <div className="flex justify-center m-4">
        <Button
          disabled={isLoading}
          onClick={() => navigate("/server/" + event?.serverId)}
        >
          Go to server
        </Button>
      </div>
    </div>
  );
}
