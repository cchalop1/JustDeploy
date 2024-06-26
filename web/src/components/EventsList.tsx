import { EventServer } from "@/hooks/useSubCreateServerEvent";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import SpinnerIcon from "@/assets/SpinnerIcon";

type EventListProps = { eventList: EventServer[]; currentStep: number };

export default function EventList({ eventList, currentStep }: EventListProps) {
  return (
    <Accordion type="single" collapsible className="w-full">
      {eventList.map((eventServer, index) => {
        const textColor =
          index < currentStep
            ? "text-black"
            : eventServer.errorMessage
            ? "text-red-500"
            : "text-gray-500";

        return (
          <AccordionItem key={index} value={`item-${index}`}>
            <AccordionTrigger>
              <div className="flex gap-5 items-center">
                <div>{index + 1}.</div>
                <div className={`${textColor}`}>{eventServer.title}</div>
                {currentStep === index && !eventServer.errorMessage ? (
                  <SpinnerIcon color="text-black" />
                ) : null}
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <code>{eventServer.errorMessage}</code>
            </AccordionContent>
          </AccordionItem>
        );
      })}
    </Accordion>
  );
}
