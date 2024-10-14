import { useNotification } from "@/hooks/useNotifications";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Service } from "@/services/getServicesByDeployId";
import { use, useEffect, useState } from "react";
import { Input } from "./ui/input";
import { ProjectDto } from "@/services/getProjectById";
import { getServerByIdApi } from "@/services/getServerById";
import { ServerDto } from "@/services/getServerListApi";

type ServiceDeploySettingsProps = {
  project: ProjectDto;
  service: Service;
  onClose: () => void;
  getProjectById: () => Promise<void>;
};

export default function ServiceDeploySettings({
  project,
  service,
  onClose,
  getProjectById,
}: ServiceDeploySettingsProps) {
  const notif = useNotification();
  const [server, setServer] = useState<null | ServerDto>(null);
  // const server = use(getServerByIdApi(project.serverId));
  const [isServiceExposed, setIsServiceExposed] = useState<boolean>(
    service.isExposed
  );

  function onCheckedChange(value: boolean) {
    setIsServiceExposed(value);
  }

  useEffect(() => {
    getServerByIdApi(project.serverId).then(setServer);
  }, [project.serverId]);

  return (
    <>
      <div className="flex items-center space-x-2 w-[25vw]">
        <Switch
          checked={isServiceExposed}
          onCheckedChange={onCheckedChange}
          id="expose-service"
        />
        <Label htmlFor="expose-service">Expose this service</Label>
      </div>
      {isServiceExposed && (
        <div className="space-y-2">
          <Label htmlFor="subdomain">Subdomain</Label>
          <div className="flex items-center space-x-2">
            <Input id="subdomain" placeholder="Enter subdomain" />
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
