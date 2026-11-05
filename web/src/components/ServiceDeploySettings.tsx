import { useEffect, useState } from "react";
import { useNotification } from "@/hooks/useNotifications";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Service } from "@/services/getServicesByDeployId";
import { saveServiceApi } from "@/services/saveServiceApi";
import { getServicesApi } from "@/services/getServicesApi";
import { deployServiceApi } from "@/services/deployApi";
import { useInfo } from "@/hooks/useInfo";
import EnvsManagements from "./forms/EnvsManagements";

type Props = {
  service: Service;
  fetchServices: () => Promise<void>;
  onConfigureServer?: () => void;
};

export default function ServiceDeploySettings({ service, fetchServices, onConfigureServer }: Props) {
  const notif = useNotification();
  const { serverInfo } = useInfo();
  const [draft, setDraft] = useState(service);
  const [services, setServices] = useState<Service[]>([]);
  const [action, setAction] = useState<"saving" | "deploying" | null>(null);
  const [dirty, setDirty] = useState(false);
  const [error, setError] = useState("");
  const [deployed, setDeployed] = useState(false);
  const missingDomain = draft.exposeSettings.isExposed && !serverInfo?.server.domain;

  useEffect(() => {
    let active = true;
    getServicesApi().then((items) => { if (active) setServices(items); }).catch(() => { /* Keep the current group available. */ });
    return () => { active = false; };
  }, []);

  function change(update: Partial<Service>) {
    setDraft((previous) => ({ ...previous, ...update }));
    setDirty(true);
    setDeployed(false);
  }

  async function save(deploy = false) {
    if (action) return;
    setAction("saving");
    setError("");
    setDeployed(false);
    try {
      const saved = dirty ? await saveServiceApi(draft) : draft;
      setDraft(saved);
      setDirty(false);
      if (deploy) {
        setAction("deploying");
        await deployServiceApi(saved.id);
        setDeployed(true);
      }
      notif.success({ title: deploy ? "Deployment complete" : "Saved", content: deploy ? "Your service has started. You can review its logs." : "Settings saved." });
    } catch (cause) {
      const message = cause instanceof Error ? cause.message : "Operation failed";
      setError(message);
      notif.error({ title: "Operation failed", content: message });
    } finally {
      // A partially completed group deployment can update service revisions.
      try {
        const items = await getServicesApi();
        setServices(items);
        const current = items.find((item) => item.id === draft.id);
        // Preserve edits when saving itself failed.
        if (current) setDraft((previous) => ({ ...previous, revision: current.revision, status: current.status, url: current.url }));
        await fetchServices();
      } catch { /* The saved draft remains available for retry. */ }
      setAction(null);
    }
  }

  const groups = new Set(services.map((item) => item.projectId || item.id));
  groups.add(draft.projectId || draft.id);
  const group = draft.projectId || draft.id;
  const dependencies = services.filter((item) => item.id !== draft.id && (item.projectId || item.id) === group);

  return (
    <fieldset disabled={action !== null} className="space-y-5 py-4">
      {draft.type === "github_repo" && (
        <div className="rounded border p-3 text-sm space-y-1">
          <p>Branch: <strong>{draft.branch}</strong></p>
          <p>Build: {draft.dockerfile || "Dockerfile auto-detection, otherwise Nixpacks (must be installed on the VPS)"}</p>
          <p>Check the application port and environment variables below before deploying.</p>
        </div>
      )}
      <div className="flex items-center justify-between">
        <label htmlFor="service-exposed">Expose this service</label>
        <Switch id="service-exposed" checked={draft.exposeSettings.isExposed} onCheckedChange={(isExposed) => change({ exposeSettings: { ...draft.exposeSettings, isExposed } })} />
      </div>
      {draft.exposeSettings.isExposed && (
        <div className="space-y-2">
          {missingDomain && <div role="alert" className="rounded border border-amber-200 bg-amber-50 p-3 text-sm space-y-2">
            <p>Configure your server domain before exposing this service, or turn off public access to deploy privately.</p>
            {onConfigureServer && <button type="button" className="underline" onClick={onConfigureServer}>Configure server domain</button>}
          </div>}
          <label htmlFor="service-subdomain">Subdomain</label>
          <input id="service-subdomain" className="border rounded p-2 w-full" value={draft.exposeSettings.subDomain} onChange={(e) => change({ exposeSettings: { ...draft.exposeSettings, subDomain: e.target.value } })} />
          <p className="text-sm text-gray-500">Leave empty for {serverInfo?.server.domain || "the root domain"}. Point this hostname to your VPS before deploying.</p>
          <label htmlFor="service-port">Application port</label>
          <input id="service-port" className="border rounded p-2 w-full" type="number" min="1" max="65535" value={draft.exposeSettings.exposePort} onChange={(e) => change({ exposeSettings: { ...draft.exposeSettings, exposePort: e.target.value } })} />
        </div>
      )}
      <div><p className="mb-2 font-medium">Environment variables</p><EnvsManagements envs={draft.envs || []} setEnvs={(envs) => change({ envs })} canEdit /></div>
      <details>
        <summary className="cursor-pointer font-medium">Advanced settings</summary>
        <div className="space-y-2 mt-3">
          <label htmlFor="service-project">Connection group</label>
          <select id="service-project" className="border rounded p-2 w-full" value={group} onChange={(e) => change({ projectId: e.target.value })}>
            {Array.from(groups).map((id) => <option key={id} value={id}>{services.find((item) => item.id === id)?.name || id}</option>)}
          </select>
          <p className="text-sm text-gray-500">Services in the same group share a private network. Add database credentials only to applications that need them.</p>
        </div>
      </details>
      {dependencies.length > 0 && <p className="text-sm text-gray-600">This deployment includes pending services in this connection group: {dependencies.map((item) => item.name).join(", ")}. Running unchanged services are kept.</p>}
      {error && <p role="alert" className="text-sm text-red-700">{error} Your settings are preserved. Correct the issue and retry; deployment details are in the Logs tab.</p>}
      {action && <p role="status" className="text-sm">{action === "saving" ? "Saving settings…" : "Deploying your service and its dependencies. Downloading or building images may take a few minutes…"}</p>}
      {deployed && <div role="status" className="rounded border border-green-200 bg-green-50 p-3 text-sm space-y-2">
        <p>Deployment complete. The container passed its startup checks.</p>
        {draft.exposeSettings.isExposed && draft.url ? <a className="underline font-medium" href={draft.url} target="_blank" rel="noreferrer">Open my application</a> : <p>This service is private. Use its name <code>{draft.name}</code> and the credentials above from an application in the same connection group.</p>}
      </div>}
      <div className="flex flex-wrap gap-2">
        <Button onClick={() => save(true)} disabled={action !== null || missingDomain}>{action === "deploying" ? "Deploying…" : error ? "Retry deployment" : "Save and deploy"}</Button>
        <Button variant="outline" onClick={() => save()} disabled={action !== null || !dirty}>Save settings</Button>
      </div>
    </fieldset>
  );
}
