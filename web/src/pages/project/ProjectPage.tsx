import { Suspense, useEffect, useState } from "react";
import { useLocation } from "react-router-dom";

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

type ProjectPageProps = {
  id: string;
};

export default function ProjectPage({ id }: ProjectPageProps) {
  const notif = useNotification();

  const [project, setProject] = useState<ProjectDto | null>(null);
  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [projectSettings, setProjectSettings] =
    useState<ProjectSettingsDto | null>(null);

  const displayWelcomeModal = useIsWelcome();

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
    await createServiceApi({
      serviceName,
      fromDockerCompose,
      projectId: project.id,
      path,
    });
    await getProjectById();
    notif.success({
      title: "Service is started",
      content: `${serviceName} is started you can now connect to if the env is generate and store in .env file in your project folder`,
    });
  };

  async function getProjectSettings() {
    const res = await getProjectSettingsByIdApi(id);
    setProjectSettings(res);
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
          setServiceSelected={setServiceSelected}
          getProjectById={getProjectById}
        />
      )}
      <ProjectPageHeader />
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
          {projectSettings && (
            <AddService
              projectSettings={projectSettings}
              createService={async (serviceParams) => {
                create({
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
