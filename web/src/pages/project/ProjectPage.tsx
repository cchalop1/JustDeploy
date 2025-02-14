import { useEffect, useState } from "react";
import AddService from "@/components/databaseServices/AddServices";
import { CreateServiceFunc } from "@/components/databaseServices/CommandModal";
import ModalWelcome from "@/components/modals/ModalWelcome";
import ModalServiceSettings from "@/components/modals/ModalServiceSettings";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceCard from "@/components/ServiceCard";
import Version from "@/components/Version";
import { Service } from "@/services/getServicesByDeployId";
import { useNotification } from "@/hooks/useNotifications";
import { useIsWelcome } from "@/hooks/useIsWelcome";
import { DeployProjectDto } from "@/services/deployProjectApi";
import { ServiceCardLoading } from "@/components/ServiceCardLoading";
import ModalGlobalSettings from "@/components/modals/ModalGlobalSettings";
import ModalCreateServer from "@/components/modals/ModalCreateServer";
import { getServicesApi } from "@/services/getServicesApi";

export default function ProjectPage() {
  const notif = useNotification();
  const displayWelcomeModal = useIsWelcome();

  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);

  const [services, setServices] = useState<Service[]>([]);

  // Modals states
  const [isGlobalSettingsModalOpen, setIsGlobalSettingsModalOpen] =
    useState(false);
  const [serviceIsLoading, setServiceIsCreating] = useState(false);
  const [displayDeployModal, setDisplayDeployModal] = useState(false);
  const [isCreateServiceModalOpen, setIsCreateServiceModalOpen] =
    useState(false);

  const toDeploy = services.find(
    (service) => service.status === "ready_to_deploy"
  );

  async function getProjectById() {
    // const project = await getProjectByIdApi(id);
    // setProject(project);
  }

  const create: CreateServiceFunc = async ({
    fromDockerCompose,
    path,
    serviceName,
  }) => {
    // if (!project) return;
    // setServiceIsCreating(true);
    // await createServiceApi({
    //   serviceName,
    //   fromDockerCompose,
    //   projectId: project.id,
    //   path,
    // });
    // await getProjectById();
    // await getProjectSettings();
    // notif.success({
    //   title: "Service is started",
    //   content: `${serviceName} is started you can now connect to if the env is generate and store in .env file in your project folder`,
    // });
    // setServiceIsCreating(false);
  };

  async function getProjectSettings() {
    // const projectSettingsResponse = await getProjectSettingsByIdApi(id);
    // setProjectSettings(projectSettingsResponse);
  }

  async function deployProject(deployProjectDto: DeployProjectDto) {
    // const response = await deployProjectApi(deployProjectDto);
  }

  async function fetchServices() {
    const services = await getServicesApi();
    setServices(services);
  }

  useEffect(() => {
    getProjectById();
    getProjectSettings();
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
      {/* {displayDeployModal && (
        <ModalDeployProject
          project={project}
          onClose={() => setDisplayDeployModal(false)}
          onDeployProject={deployProject}
        />
      )} */}
      <ProjectPageHeader
        onClickDeploy={() => setDisplayDeployModal(true)}
        onClickSettings={() =>
          setIsGlobalSettingsModalOpen(!isGlobalSettingsModalOpen)
        }
        toDeploy={!!toDeploy}
      />
      {displayWelcomeModal && <ModalWelcome />}
      <div className="flex flex-col justify-center items-center h-3/5">
        <div className="flex flex-col gap-3 mt-3 ">
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
              onClick={() => setServiceSelected(service)}
            />
          ))}
          {serviceIsLoading && <ServiceCardLoading />}
          <AddService
            fetchServices={fetchServices}
            setLoading={setServiceIsCreating}
          />
        </div>
      </div>
      <div className="fixed bottom-6 right-4 pl-10 pr-10">
        <Version />
      </div>
    </div>
  );
}
