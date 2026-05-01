import { motion } from "framer-motion";
import { Suspense, useState } from "react";
import { Trash2 } from "lucide-react";

import Modal from "@/components/modals/Modal";
import { Service } from "@/services/getServicesByDeployId";
import ServiceDeploySettings from "../ServiceDeploySettings";
import { CardIcon } from "../CardIcon";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import { useNotification } from "@/hooks/useNotifications";
import ServiceLogs from "../ServiceLogs";
import ServiceCommitInfo from "../ServiceCommitInfo";

type ModalServiceSettingsProps = {
  service: Service;
  onClose: () => void;
  fetchServices: () => Promise<void>;
};

function statusColor(status: string) {
  if (status === "running") return "bg-green-100 text-green-700";
  if (status === "error" || status === "failed") return "bg-red-100 text-red-600";
  if (status === "ready_to_deploy") return "bg-blue-50 text-blue-600";
  return "bg-gray-100 text-gray-500";
}

export default function ModalServiceSettings({ service, onClose, fetchServices }: ModalServiceSettingsProps) {
  const notif = useNotification();
  const [isDeleting, setIsDeleting] = useState(false);
  const [tab, setTab] = useState<"settings" | "logs">("settings");

  async function deleteServiceById() {
    setIsDeleting(true);
    try {
      await deleteServiceByIdApi(service.id);
      onClose();
      notif.success({ title: "Service deleted", content: `${service.name} has been deleted.` });
      await fetchServices();
    } catch {
      notif.error({ title: "Error", content: "Could not delete service." });
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, x: 40 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      style={{ position: "absolute", right: 5, top: 28, width: 520, zIndex: 20 }}
    >
      <Modal
        onClose={onClose}
        headerNode={
          <div className="flex items-center gap-3">
            <CardIcon service={service} />
            <div>
              <p className="text-sm font-semibold text-gray-900">{service.name}</p>
              <span className={`inline-block mt-0.5 text-xs font-medium px-2 py-0.5 rounded-full ${statusColor(service.status)}`}>
                {service.status}
              </span>
            </div>
          </div>
        }
        className="w-[520px] max-w-full"
      >
        {/* Tabs */}
        <div className="flex border-b border-gray-100 mt-1 px-5">
          {(["settings", "logs"] as const).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`pb-2 mr-5 text-sm font-medium border-b-2 transition-colors ${
                tab === t
                  ? "border-gray-900 text-gray-900"
                  : "border-transparent text-gray-400 hover:text-gray-600"
              }`}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>

        <div className="px-5 overflow-y-auto max-h-[60vh]">
          {tab === "settings" && (
            <>
              <ServiceCommitInfo serviceId={service.id} />
              <ServiceDeploySettings service={service} fetchServices={fetchServices} />

              {/* Danger zone */}
              <div className="pt-4 pb-2 border-t border-gray-100 mt-2">
                <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-3">
                  Danger zone
                </p>
                <button
                  onClick={deleteServiceById}
                  disabled={isDeleting}
                  className="flex items-center gap-2 text-sm text-red-500 hover:text-red-700 border border-red-200 hover:border-red-400 hover:bg-red-50 px-3 py-1.5 rounded-md transition-colors disabled:opacity-50"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  {isDeleting ? "Deleting…" : "Delete service"}
                </button>
              </div>
            </>
          )}

          {tab === "logs" && (
            <div className="py-4">
              <Suspense fallback={<p className="text-xs text-gray-400">Loading…</p>}>
                <ServiceLogs serviceId={service.id} />
              </Suspense>
            </div>
          )}
        </div>
      </Modal>
    </motion.div>
  );
}
