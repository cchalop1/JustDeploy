import { useEffect, useRef, useState } from "react";
import { useNotification } from "@/hooks/useNotifications";
import { Switch } from "@/components/ui/switch";
import { Service } from "@/services/getServicesByDeployId";
import { ServerDto } from "@/services/getServerListApi";
import { saveServiceApi } from "@/services/saveServiceApi";
import { getServerInfoApi } from "@/services/getServerInfoApi";
import EnvsManagements from "./forms/EnvsManagements";
import { Env } from "@/services/postFormDetails";

type ServiceDeploySettingsProps = {
  service: Service;
  fetchServices: () => Promise<void>;
};

function SectionTitle({ children }: { children: React.ReactNode }) {
  return (
    <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-3">
      {children}
    </p>
  );
}

export default function ServiceDeploySettings({ service, fetchServices }: ServiceDeploySettingsProps) {
  const notif = useNotification();
  const timeoutRef = useRef<number | null>(null);
  const [server, setServer] = useState<null | ServerDto>(null);
  const [isExposed, setIsExposed] = useState<boolean>(service.exposeSettings.isExposed);
  const [subDomain, setSubdomain] = useState<string>(service.exposeSettings.subDomain);
  const [envs, setLocalEnvs] = useState<Env[]>(service.envs || []);

  useEffect(() => {
    getServerInfoApi().then(setServer);
  }, []);

  function debounce(fn: () => void) {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = window.setTimeout(fn, 600);
  }

  async function saveService(updated: Service) {
    try {
      await saveServiceApi({ ...updated, status: "ready_to_deploy" });
      await fetchServices();
      notif.success({ title: "Saved", content: "Settings updated." });
    } catch (e) {
      notif.error({ title: "Error", content: (e as Error).message });
    }
  }

  function onExposedChange(value: boolean) {
    setIsExposed(value);
    debounce(() =>
      saveService({ ...service, exposeSettings: { ...service.exposeSettings, isExposed: value } })
    );
  }

  function onSubdomainChange(value: string) {
    setSubdomain(value);
    debounce(() =>
      saveService({ ...service, exposeSettings: { ...service.exposeSettings, subDomain: value } })
    );
  }

  function handleEnvsChange(updated: Env[]) {
    setLocalEnvs(updated);
    debounce(() => saveService({ ...service, envs: updated }));
  }

  return (
    <div className="divide-y divide-gray-100">

      {/* Expose */}
      <div className="py-4">
        <SectionTitle>Visibility</SectionTitle>
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-800">Expose this service</p>
            <p className="text-xs text-gray-400 mt-0.5">
              Make it accessible via a public URL.
            </p>
          </div>
          <Switch checked={isExposed} onCheckedChange={onExposedChange} />
        </div>

        {isExposed && (
          <div className="mt-4">
            <label className="text-xs text-gray-500 mb-1.5 block">Subdomain</label>
            <div className="flex items-center gap-0">
              <input
                className="flex-1 h-9 px-3 rounded-l-md border border-r-0 border-gray-200 bg-white text-sm font-mono focus:outline-none focus:ring-1 focus:ring-gray-300"
                placeholder="myapp"
                value={subDomain}
                onChange={(e) => onSubdomainChange(e.target.value)}
              />
              <span className="h-9 px-3 flex items-center rounded-r-md border border-gray-200 bg-gray-50 text-xs text-gray-400 font-mono whitespace-nowrap">
                .{server?.server.domain || "yourdomain.com"}
              </span>
            </div>
            <p className="text-xs text-gray-400 mt-1.5">
              Leave empty to use the root domain.
            </p>
          </div>
        )}
      </div>

      {/* Env vars */}
      <div className="py-4">
        <SectionTitle>Environment variables</SectionTitle>
        <EnvsManagements envs={envs} setEnvs={handleEnvsChange} canEdit />
      </div>

    </div>
  );
}
