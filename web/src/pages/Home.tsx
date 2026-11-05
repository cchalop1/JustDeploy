import { useEffect, useState } from "react";
import AddService from "@/components/databaseServices/AddServices";
import ModalServiceSettings from "@/components/modals/ModalServiceSettings";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceCard from "@/components/ServiceCard";
import Version from "@/components/Version";
import { Service } from "@/services/getServicesByDeployId";
import { useNotification } from "@/hooks/useNotifications";
import { ServiceCardLoading } from "@/components/ServiceCardLoading";
import ModalGlobalSettings from "@/components/modals/ModalGlobalSettings";
import ModalFirstConnection from "@/components/modals/ModalFirstConnection";
import ModalLogin from "@/components/modals/ModalLogin";
import { getServicesApi } from "@/services/getServicesApi";
import { deployApi } from "@/services/deployApi";
import { useInfo } from "@/hooks/useInfo";
import { hasToken } from "@/services/authStorage";

export default function Home() {
  const notif = useNotification();
  const { serverInfo } = useInfo();

  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [isGlobalSettingsModalOpen, setIsGlobalSettingsModalOpen] = useState(false);
  const [serviceIsLoading, setServiceIsCreating] = useState(false);
  const [servicesLoaded, setServicesLoaded] = useState(false);

  const toDeploy = services.find((s) => s.status === "ready_to_deploy");

  const isFirstConnection = serverInfo?.firstConnection ?? false;
  const [isLoggedIn, setIsLoggedIn] = useState(hasToken);
  useEffect(() => {
    const update = () => setIsLoggedIn(hasToken());
    window.addEventListener("auth-changed", update);
    const timer = window.setInterval(update, 30000);
    return () => { window.removeEventListener("auth-changed", update); window.clearInterval(timer); };
  }, []);

  const [shouldShowFirstConnectionModal, setShouldShowFirstConnectionModal] = useState(false);
  const [shouldShowLoginModal, setShouldShowLoginModal] = useState(false);

  useEffect(() => {
    if (serverInfo === null) return;
    if (isFirstConnection && !isLoggedIn) {
      setShouldShowFirstConnectionModal(true);
    } else if (!isLoggedIn) {
      setShouldShowLoginModal(true);
    } else {
      setShouldShowFirstConnectionModal(false);
      setShouldShowLoginModal(false);
    }
  }, [serverInfo, isFirstConnection, isLoggedIn]);

  useEffect(() => {
    if (!shouldShowFirstConnectionModal && !shouldShowLoginModal && isLoggedIn) {
      fetchServices();
    }
  }, [shouldShowFirstConnectionModal, shouldShowLoginModal, isLoggedIn]);

  async function deploy() {
    notif.info({ title: "Deploying project", content: "This may take a few minutes..." });
    try {
      const res = await deployApi();
      if (res.message === "Deployed") {
        notif.success({ title: "Project deployed", content: "Your project has been deployed successfully" });
        fetchServices();
      }
    } catch (error: unknown) {
      notif.error({
        title: "Error deploying project",
        content: error instanceof Error ? error.message : "Unknown error occurred",
      });
    }
  }

  async function fetchServices() {
    try {
      const items = await getServicesApi();
      setServices(items);
      setServicesLoaded(true);
      setServiceSelected((selected) => selected ? items.find((item) => item.id === selected.id) || null : null);
    } catch (error) { console.error("Services unavailable", error); }
  }

  useEffect(() => {
    const onKeyDown = (ev: globalThis.KeyboardEvent) => {
      if (ev.key === "Escape") setServiceSelected(null);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);

  const githubRepoServices = services.filter((s) => s.type === "github_repo");
  const otherServices = services.filter((s) => s.type !== "github_repo");

  return (
    <div className="bg-grid-image min-h-screen">
      {serviceSelected && (
        <ModalServiceSettings
          key={serviceSelected.id}
          service={serviceSelected}
          onConfigureServer={() => setIsGlobalSettingsModalOpen(true)}
          onClose={() => setServiceSelected(null)}
          fetchServices={fetchServices}
        />
      )}
      {isGlobalSettingsModalOpen && (
        <ModalGlobalSettings onClose={() => setIsGlobalSettingsModalOpen(false)} />
      )}
      {shouldShowFirstConnectionModal && (
        <ModalFirstConnection
          onClose={() => {
            setShouldShowFirstConnectionModal(false);
            fetchServices();
          }}
        />
      )}
      {shouldShowLoginModal && (
        <ModalLogin
          onClose={() => {
            setShouldShowLoginModal(false);
            fetchServices();
          }}
        />
      )}
      <ProjectPageHeader
        onClickDeploy={deploy}
        onClickSettings={() => setIsGlobalSettingsModalOpen(!isGlobalSettingsModalOpen)}
        toDeploy={!!toDeploy}
      />
      <div className="flex flex-col justify-center items-center min-h-[60vh] px-4 pb-12">
        {isLoggedIn && !shouldShowFirstConnectionModal && !shouldShowLoginModal && (
          <div className="flex flex-col items-center gap-6 mt-3 w-full">
            {!servicesLoaded && <p role="status">Loading your services…</p>}
            {servicesLoaded && services.length === 0 && (
              <section className="text-center max-w-xl space-y-3">
                <p className="text-sm text-gray-500">Account ready · Next: your first service</p>
                <h1 className="text-3xl font-bold">Deploy your first service</h1>
                <p className="text-gray-600">Choose an application or a private database. Review its settings, then deploy from the same screen.</p>
              </section>
            )}
            <div className="flex flex-wrap justify-center gap-3">
              {githubRepoServices.map((service) => (
                <ServiceCard
                  key={service.id}
                  service={service}
                  onClick={() => { setServiceSelected(null); setServiceSelected(service); }}
                />
              ))}
            </div>
            <div className="flex flex-wrap justify-center gap-3">
              {otherServices.map((service) => (
                <ServiceCard
                  key={service.id}
                  service={service}
                  onClick={() => { setServiceSelected(null); setServiceSelected(service); }}
                />
              ))}
              {serviceIsLoading && <ServiceCardLoading />}
              {servicesLoaded && <AddService
                fetchServices={fetchServices}
                setLoading={setServiceIsCreating}
                firstService={services.length === 0}
                onCreated={setServiceSelected}
              />}
            </div>
          </div>
        )}
      </div>
      <div className="fixed bottom-6 right-4 pl-10 pr-10">
        <Version />
      </div>
    </div>
  );
}
