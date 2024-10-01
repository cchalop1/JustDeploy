import { useEffect, useState } from "react";
import SpinnerIcon from "@/assets/SpinnerIcon";
import AddService from "@/components/databaseServices/AddServices";
import { CreateServiceFunc } from "@/components/databaseServices/CommandModal";
import ModalWelcome from "@/components/modals/ModalWelcome";
import ModalServiceSettings from "@/components/modals/ModalServiceSettings";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceCard from "@/components/ServiceCard";
import Version from "@/components/Version";
import { createServiceApi } from "@/services/createServiceApi";
import { getProjectByIdApi, ProjectDto } from "@/services/getProjectById";
import { Service } from "@/services/getServicesByDeployId";
import { useNotification } from "@/hooks/useNotifications";
import {
  getProjectSettingsByIdApi,
  ProjectSettingsDto,
} from "@/services/getProjectSettings";
import { useIsWelcome } from "@/hooks/useIsWelcome";
import {
  DeployProjectDto,
  deployProjectApi,
} from "@/services/deployProjectApi";
import { ServiceCardLoading } from "@/components/ServiceCardLoading";
import ModalGlobalSettings from "@/components/modals/ModalGlobalSettings";
import ModalDeployProject from "@/components/modals/ModalDeployProject";

type ProjectPageProps = {
  id: string;
};

export default function ProjectPage({ id }: ProjectPageProps) {
  const notif = useNotification();

  const [project, setProject] = useState<ProjectDto | null>(null);
  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [projectSettings, setProjectSettings] =
    useState<ProjectSettingsDto | null>(null);

  // Modals states
  const displayWelcomeModal = useIsWelcome();
  const [isGlobalSettingsModalOpen, setIsGlobalSettingsModalOpen] =
    useState(false);
  const [serviceIsLoading, setServiceIsCreating] = useState(false);
  const [displayDeployModal, setDisplayDeployModal] = useState(false);

  const apps = project?.services.filter((s) => s.isDevContainer) || [];
  const services = project?.services.filter((s) => !s.isDevContainer) || [];

  async function getProjectById() {
    const project = await getProjectByIdApi(id);
    setProject(project);
  }

  const create: CreateServiceFunc = async ({
    fromDockerCompose,
    path,
    serviceName,
  }) => {
    if (!project) return;
    setServiceIsCreating(true);
    await createServiceApi({
      serviceName,
      fromDockerCompose,
      projectId: project.id,
      path,
    });
    await getProjectById();
    await getProjectSettings();
    notif.success({
      title: "Service is started",
      content: `${serviceName} is started you can now connect to if the env is generate and store in .env file in your project folder`,
    });
    setServiceIsCreating(false);
  };

  async function getProjectSettings() {
    const projectSettingsResponse = await getProjectSettingsByIdApi(id);

    setProjectSettings(projectSettingsResponse);
  }

  async function deployProject(deployProjectDto: DeployProjectDto) {
    const response = await deployProjectApi(deployProjectDto);
  }

  useEffect(() => {
    getProjectById();
    getProjectSettings();
  }, [id]);

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
          getProjectById={getProjectById}
        />
      )}
      {isGlobalSettingsModalOpen && (
        <ModalGlobalSettings
          onClose={() => setIsGlobalSettingsModalOpen(false)}
        />
      )}
      {displayDeployModal && (
        <ModalDeployProject
          projectId={id}
          onClose={() => setDisplayDeployModal(false)}
          onDeployProject={deployProject}
        />
      )}
      <ProjectPageHeader
        onClickDeploy={() => setDisplayDeployModal(true)}
        onClickSettings={() =>
          setIsGlobalSettingsModalOpen(!isGlobalSettingsModalOpen)
        }
      />
      {displayWelcomeModal && <ModalWelcome />}
      <div className="flex flex-col justify-center items-center h-3/5">
        <div className="flex gap-3">
          {apps.map((app) => (
            <ServiceCard
              key={app.id}
              service={app}
              onClick={() => setServiceSelected(app)}
            />
          ))}
        </div>
        {!project && <SpinnerIcon color="text-black" />}
        <div className="flex gap-3 mt-3 ">
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
              onClick={() => setServiceSelected(service)}
            />
          ))}
          {serviceIsLoading && <ServiceCardLoading />}
          {projectSettings && (
            <AddService
              projectId={project?.id}
              projectSettings={projectSettings}
              createService={async (serviceParams) => {
                await create({
                  serviceName: serviceParams.serviceName,
                  fromDockerCompose: serviceParams.fromDockerCompose,
                  path: serviceParams.path,
                });
              }}
              fetchServiceList={getProjectById}
              setLoading={() => {}}
            />
          )}
        </div>
      </div>
      <div className="fixed bottom-6 right-4 pl-10 pr-10">
        <Version />
      </div>
    </div>
  );
}
