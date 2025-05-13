import { useEffect, useState } from "react";
import AddService from "@/components/databaseServices/AddServices";
import ModalWelcome from "@/components/modals/ModalWelcome";
import ModalServiceSettings from "@/components/modals/ModalServiceSettings";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceCard from "@/components/ServiceCard";
import Version from "@/components/Version";
import { Service } from "@/services/getServicesByDeployId";
import { useNotification } from "@/hooks/useNotifications";
import { useIsWelcome } from "@/hooks/useIsWelcome";
import { ServiceCardLoading } from "@/components/ServiceCardLoading";
import ModalGlobalSettings from "@/components/modals/ModalGlobalSettings";
import ModalFirstConnection from "@/components/modals/ModalFirstConnection";
import { getServicesApi } from "@/services/getServicesApi";
import { deployApi } from "@/services/deployApi";
import { useInfo } from "@/hooks/useInfo";
import { hasApiKey } from "@/services/authStorage";

export default function Home() {
  const notif = useNotification();
  const displayWelcomeModal = useIsWelcome();

  const { serverInfo } = useInfo();

  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [isGlobalSettingsModalOpen, setIsGlobalSettingsModalOpen] =
    useState(false);
  const [serviceIsLoading, setServiceIsCreating] = useState(false);

  const toDeploy = services.find(
    (service) => service.status === "ready_to_deploy"
  );

  // Check if we should display the first connection modal
  const [shouldShowFirstConnectionModal, setShouldShowFirstConnectionModal] =
    useState(serverInfo?.firstConnection || !hasApiKey());

  async function deploy() {
    notif.info({
      title: "Deploying project",
      content: "This may take a few minutes... Loading...",
    });
    try {
      const res = await deployApi();
      if (res.message === "Deployed") {
        notif.success({
          title: "Project deployed",
          content: "Your project has been deployed successfully",
        });
        fetchServices();
      }
    } catch (error: unknown) {
      notif.error({
        title: "Error deploying project",
        content:
          error instanceof Error ? error.message : "Unknown error occurred",
      });
    }
  }

  async function fetchServices() {
    // Only fetch services if authenticated
    const services = await getServicesApi();
    setServices(services);
  }

  useEffect(() => {
    // Only fetch services if we're not showing the first connection modal
    if (!shouldShowFirstConnectionModal) {
      fetchServices();
    }
  }, [shouldShowFirstConnectionModal]);

  useEffect(() => {
    window.addEventListener("keydown", (ev: globalThis.KeyboardEvent) => {
      if (ev.key === "Escape") {
        setServiceSelected(null);
      }
    });
    return () => {
      window.removeEventListener("keydown", () => {});
    };
  }, []);

  // Filter services by type
  const githubRepoServices = services.filter((s) => s.type === "github_repo");
  const otherServices = services.filter((s) => s.type !== "github_repo");

  // Authentication Modal component
  return (
    <div className="bg-grid-image h-screen">
      {serviceSelected && (
        <ModalServiceSettings
          service={serviceSelected}
          onClose={() => setServiceSelected(null)}
          fetchServices={fetchServices}
        />
      )}
      {isGlobalSettingsModalOpen && (
        <ModalGlobalSettings
          onClose={() => setIsGlobalSettingsModalOpen(false)}
        />
      )}
      {shouldShowFirstConnectionModal && (
        <ModalFirstConnection
          serverIp={serverInfo?.server?.ip || ""}
          onClose={() => setShouldShowFirstConnectionModal(false)}
        />
      )}
      <ProjectPageHeader
        onClickDeploy={deploy}
        onClickSettings={() =>
          setIsGlobalSettingsModalOpen(!isGlobalSettingsModalOpen)
        }
        toDeploy={!!toDeploy}
      />
      {displayWelcomeModal && <ModalWelcome />}
      <div className="flex flex-col justify-center items-center h-3/5">
        {!shouldShowFirstConnectionModal && (
          <div className="flex flex-col items-center gap-6 mt-3">
            {/* GitHub Repo Services */}
            <div className="flex gap-3">
              {githubRepoServices.map((service) => (
                <ServiceCard
                  key={service.id}
                  service={service}
                  onClick={() => {
                    setServiceSelected(null);
                    setServiceSelected(service);
                  }}
                />
              ))}
            </div>

            {/* Other Services */}
            <div className="flex gap-3">
              {otherServices.map((service) => (
                <ServiceCard
                  key={service.id}
                  service={service}
                  onClick={() => {
                    setServiceSelected(null);
                    setServiceSelected(service);
                  }}
                />
              ))}
              {serviceIsLoading && <ServiceCardLoading />}
              <AddService
                fetchServices={fetchServices}
                setLoading={setServiceIsCreating}
              />
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
