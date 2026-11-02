import { useEffect, useState } from "react";
import { useNotification } from "@/hooks/useNotifications";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Service } from "@/services/getServicesByDeployId";
import { saveServiceApi } from "@/services/saveServiceApi";
import { getServicesApi } from "@/services/getServicesApi";
import { useInfo } from "@/hooks/useInfo";
import EnvsManagements from "./forms/EnvsManagements";

type Props = { service: Service; fetchServices: () => Promise<void> };
export default function ServiceDeploySettings({ service, fetchServices }: Props) {
 const notif = useNotification();
 const { serverInfo } = useInfo();
 const [draft, setDraft] = useState(service);
 const [services, setServices] = useState<Service[]>([]);
 const [saving, setSaving] = useState(false);
 const [dirty, setDirty] = useState(false);
 useEffect(() => {
  let active = true;
  getServicesApi().then((items) => { if (active) setServices(items); }).catch(() => { /* The current group remains selectable if loading fails. */ });
  return () => { active = false; };
 }, []);
 function change(update: Partial<Service>) { setDraft((previous) => ({...previous, ...update})); setDirty(true); }
 async function save() {
  setSaving(true);
  try {
   const saved = await saveServiceApi(draft);
   setDraft(saved); setDirty(false);
   await fetchServices();
   notif.success({title:"Saved", content:"Settings saved. Deploy the service to apply them."});
  } catch (error) { notif.error({title:"Save failed", content: error instanceof Error ? error.message : "Save failed"}); }
  finally { setSaving(false); }
 }
 const groups = new Set(services.map((item) => item.projectId || item.id));
 groups.add(draft.projectId || draft.id);
 return (
  <fieldset disabled={saving} className="space-y-5 py-4">
   <div className="flex items-center justify-between">
    <label htmlFor="service-exposed">Expose this service</label>
    <Switch id="service-exposed" checked={draft.exposeSettings.isExposed} onCheckedChange={(isExposed) => change({exposeSettings:{...draft.exposeSettings,isExposed}})} />
   </div>
   {draft.exposeSettings.isExposed && <div className="space-y-2">
    <label htmlFor="service-subdomain">Subdomain</label>
    <input id="service-subdomain" className="border rounded p-2 w-full" value={draft.exposeSettings.subDomain} onChange={(e) => change({exposeSettings:{...draft.exposeSettings,subDomain:e.target.value}})} />
    <p className="text-sm text-gray-500">Leave empty for {serverInfo?.server.domain || "the root domain"}.</p>
    <label htmlFor="service-port">Application port</label>
    <input id="service-port" className="border rounded p-2 w-full" type="number" min="1" max="65535" value={draft.exposeSettings.exposePort} onChange={(e) => change({exposeSettings:{...draft.exposeSettings,exposePort:e.target.value}})} />
   </div>}
   <div className="space-y-2">
    <label htmlFor="service-project">Connection group</label>
    <select id="service-project" className="border rounded p-2 w-full" value={draft.projectId || draft.id} onChange={(e) => change({projectId:e.target.value})}>
     {Array.from(groups).map((group) => <option key={group} value={group}>{services.find((item) => item.id === group)?.name || group}</option>)}
    </select>
    <p className="text-sm text-gray-500">Only services in the same group can connect through their private network. Add database credentials only to the apps that need them.</p>
   </div>
   <div><p className="mb-2 font-medium">Environment variables</p><EnvsManagements envs={draft.envs || []} setEnvs={(envs) => change({envs})} canEdit /></div>
   <Button onClick={save} disabled={saving || !dirty}>{saving ? "Saving…" : "Save settings"}</Button>
  </fieldset>
 );
}
