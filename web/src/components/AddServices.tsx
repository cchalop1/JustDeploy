import { ServiceDto, getServiceListApi } from "@/services/getServicesApi";
import { useEffect, useState } from "react";
import DatabaseCard from "./DatabaseCard";
import { Button } from "./ui/button";
import { FileSlidersIcon } from "lucide-react";

import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "./ui/command";
import { createServiceApi } from "@/services/createServiceApi";

type AddServiceProps = {
  deployId: string;
  setLoading: (loading: boolean) => void;
  fetchServiceList: (deployId: string) => Promise<void>;
};

export default function AddService({
  deployId,
  fetchServiceList,
  setLoading,
}: AddServiceProps) {
  const [services, setServices] = useState<Array<ServiceDto>>([]);
  const [open, setOpen] = useState(false);

  async function getServices() {
    const res = await getServiceListApi();
    setServices(res);
  }

  useEffect(() => {
    getServices();
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  async function createService(serviceName: string) {
    try {
      setLoading(true);
      setOpen(false);
      await createServiceApi(serviceName, deployId);
      await fetchServiceList(deployId);
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  }

  return (
    <>
      <Button
        variant="outline"
        className="w-full mb-3"
        onClick={() => setOpen(true)}
      >
        Click here for create and connect new database or press{" "}
        <kbd className="ml-2 pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
          <span className="text-xs">⌘</span>K
        </kbd>
      </Button>
      <CommandDialog open={open} onOpenChange={setOpen}>
        <CommandInput placeholder="Type to search and lauch a service in the list..." />
        <CommandList onSelect={() => setOpen(false)}>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Local Config">
            <CommandItem className="flex gap-3" value="compose">
              <FileSlidersIcon className="w-5"></FileSlidersIcon>
              <span className="h-4">From your docker-compose file</span>
            </CommandItem>
          </CommandGroup>
          <CommandSeparator />
          <CommandGroup heading="Other sercices">
            {services.map((s) => (
              <DatabaseCard key={s.name} service={s} onSelect={createService} />
            ))}
            {/* TODO: command item to add new service link to github issue template */}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </>
  );
}
