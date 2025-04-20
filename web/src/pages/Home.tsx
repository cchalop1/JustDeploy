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
import ModalCreateServer from "@/components/modals/ModalCreateServer";
import { getServicesApi } from "@/services/getServicesApi";
import { deployApi } from "@/services/deployApi";

export default function Home() {
  const notif = useNotification();
  const displayWelcomeModal = useIsWelcome();

  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [showAuthModal, setShowAuthModal] = useState<boolean>(false);

  // Modals states
  const [isGlobalSettingsModalOpen, setIsGlobalSettingsModalOpen] =
    useState(false);
  const [serviceIsLoading, setServiceIsCreating] = useState(false);
  const [isCreateServiceModalOpen, setIsCreateServiceModalOpen] =
    useState(false);

  const toDeploy = services.find(
    (service) => service.status === "ready_to_deploy"
  );

  async function getProjectSettings() {
    // const projectSettingsResponse = await getProjectSettingsByIdApi(id);
    // setProjectSettings(projectSettingsResponse);
  }

  async function deploy() {
    if (!isAuthenticated) {
      setShowAuthModal(true);
      return;
    }

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
    // if (isAuthenticated) {
    const services = await getServicesApi();
    setServices(services);
    // }
  }

  async function getAuthKey() {
    const urlParams = new URLSearchParams(window.location.search);
    const apiKey = urlParams.get("api_key");

    if (apiKey) {
      // Store the API key in localStorage
      localStorage.setItem("api_key", apiKey);

      // Remove the api_key parameter from the URL without refreshing the page
      urlParams.delete("api_key");
      const newUrl =
        window.location.pathname +
        (urlParams.toString() ? `?${urlParams.toString()}` : "");
      window.history.replaceState({}, document.title, newUrl);
      setIsAuthenticated(true);
    } else {
      // Check if API key exists in localStorage
      const storedApiKey = localStorage.getItem("api_key");
      if (!storedApiKey) {
        setShowAuthModal(true);
        setIsAuthenticated(false);
      } else {
        setIsAuthenticated(true);
      }
    }
  }

  useEffect(() => {
    getAuthKey();
    getProjectSettings();
    // fetchServices will only run if authenticated
    fetchServices();
  }, []);

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

  // Authentication Modal component
  const AuthenticationModal = () => {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg p-8 max-w-lg w-full">
          <h2 className="text-2xl font-bold mb-4">Authentication Required</h2>
          <p className="mb-6">
            You are not authenticated. Please provide an API key to continue.
          </p>
          <div className="flex justify-end">
            <button
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              onClick={() => setShowAuthModal(false)}
            >
              Close
            </button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="bg-grid-image h-screen">
      {serviceSelected && (
        <ModalServiceSettings
          service={serviceSelected}
          onClose={() => setServiceSelected(null)}
          fetchServices={fetchServices}
        />
      )}
      {isCreateServiceModalOpen && (
        <ModalCreateServer onClose={() => setIsCreateServiceModalOpen(false)} />
      )}
      {isGlobalSettingsModalOpen && (
        <ModalGlobalSettings
          onClose={() => setIsGlobalSettingsModalOpen(false)}
          onClickNewServer={() => {
            setIsCreateServiceModalOpen(true);
          }}
        />
      )}
      {showAuthModal && <AuthenticationModal />}
      <ProjectPageHeader
        onClickDeploy={deploy}
        onClickSettings={() =>
          setIsGlobalSettingsModalOpen(!isGlobalSettingsModalOpen)
        }
        toDeploy={!!toDeploy}
      />
      {displayWelcomeModal && <ModalWelcome />}
      <div className="flex flex-col justify-center items-center h-3/5">
        <div className="flex flex-col gap-3 mt-3 ">
          {isAuthenticated ? (
            <>
              {services.map((service) => (
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
            </>
          ) : (
            <div className="text-center p-4 bg-gray-100 rounded-lg shadow">
              <p className="text-gray-700">
                Please authenticate to view and manage your services.
              </p>
            </div>
          )}
        </div>
      </div>
      <div className="fixed bottom-6 right-4 pl-10 pr-10">
        <Version />
      </div>
    </div>
  );
}
