import { useEffect, useRef, useState } from "react";

import { useNotification } from "@/hooks/useNotifications";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Service } from "@/services/getServicesByDeployId";
import { Input } from "@/components/ui/input";
import { ServerDto } from "@/services/getServerListApi";
import { saveServiceApi } from "@/services/saveServiceApi";
import { getServerInfoApi } from "@/services/getServerInfoApi";

type ServiceDeploySettingsProps = {
  service: Service;
};

export default function ServiceDeploySettings({
  service,
}: ServiceDeploySettingsProps) {
  const notif = useNotification();
  const timeoutRef = useRef<number | null>(null);
  const [server, setServer] = useState<null | ServerDto>(null);

  const [isExposed, setIsServiceExposed] = useState<boolean>(
    service.exposeSettings.isExposed
  );
  const [subDomain, setSubdomain] = useState<string>(
    service.exposeSettings.subDomain
  );

  async function saveService(serviceUpdated: Service) {
    try {
      const res = await saveServiceApi(serviceUpdated, project.id);
      console.log(res);
    } catch (e) {
      notif.error({
        title: "Error",
        content: e.message,
      });
      return;
    }
    notif.success({
      title: "Settings saved",
      content: "Service settings have been saved !",
    });
  }

  function onCheckedChange(value: boolean) {
    setIsServiceExposed(value);

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = window.setTimeout(() => {
      saveService({
        ...service,
        exposeSettings: { ...service.exposeSettings, isExposed: value },
      });
    }, 1000);
  }

  function onSubdomainChange(value: string) {
    setSubdomain(value);

    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = window.setTimeout(() => {
      saveService({
        ...service,
        exposeSettings: { ...service.exposeSettings, subDomain: value },
      });
    }, 1000);
  }

  useEffect(() => {
    getServerInfoApi().then(setServer);
  }, []);

  return (
    <>
      <div className="flex items-center space-x-2 w-[25vw]">
        <Switch
          checked={isExposed}
          onCheckedChange={onCheckedChange}
          id="expose-service"
        />
        <Label htmlFor="expose-service">Expose this service</Label>
      </div>
      {isExposed && (
        <div className="space-y-2">
          <Label htmlFor="subdomain">Subdomain</Label>
          <div className="flex items-center space-x-2">
            <Input
              id="subdomain"
              placeholder="Enter subdomain"
              value={subDomain}
              onChange={(e) => onSubdomainChange(e.target.value)}
            />
            <span className="text-sm text-muted-foreground">
              {" "}
              .{server?.domain}
            </span>
          </div>
          <Label htmlFor="subdomain" className="text-gray-500">
            Leave the domain empty if you want only the base domain
          </Label>
        </div>
      )}
    </>
  );
}
