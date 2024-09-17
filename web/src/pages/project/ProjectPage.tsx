import SpinnerIcon from "@/assets/SpinnerIcon";
import AlertModal from "@/components/alerts/AlertModal";
import AddService from "@/components/databaseServices/AddServices";
import { CreateServiceFunc } from "@/components/databaseServices/CommandModal";
import ModalInfo from "@/components/modals/ModalInfo";
import ModalServiceSettings from "@/components/modals/ModalServiceSettings";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceCard from "@/components/ServiceCard";
import Version from "@/components/Version";
import { createServiceApi } from "@/services/createServiceApi";
import { getProjectByIdApi, ProjectDto } from "@/services/getProjectById";
import { Service } from "@/services/getServicesByDeployId";
import { Suspense, useEffect, useState } from "react";
import { motion } from "framer-motion";

type ProjectPageProps = {
  id: string;
};

export default function ProjectPage({ id }: ProjectPageProps) {
  const [project, setProject] = useState<ProjectDto | null>(null);
  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);

  const apps = project?.services.filter((s) => s.isDevContainer) || [];
  const services = project?.services.filter((s) => !s.isDevContainer) || [];

  // TODO: move to context
  const [displayAlert, setDisplayAlert] = useState(false);

  async function getProjectById() {
    const project = await getProjectByIdApi(id);
    setProject(project);
  }

  console.log("serviceSelected", serviceSelected);

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
      localPath: path,
    });
    await getProjectById();
    setDisplayAlert(true);
    setTimeout(() => {
      setDisplayAlert(false);
    }, 3000);
  };

  useEffect(() => {
    getProjectById();
  }, [id]);

  return (
    <div className="bg-grid-image h-screen">
      {displayAlert && (
        <AlertModal
          type="success"
          title="Service Created !"
          message="Your service has been created successfully, you can find its variables in .env"
        />
      )}

      {serviceSelected && (
        <ModalServiceSettings
          service={serviceSelected}
          setServiceSelected={setServiceSelected}
          getProjectById={getProjectById}
        />
      )}
      <ProjectPageHeader />
      <ModalInfo />
      <div className="flex flex-col justify-center items-center h-3/5">
        <div>
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
          <AddService
            createService={async (serviceParams) => {
              create({
                serviceName: serviceParams.serviceName,
                fromDockerCompose: serviceParams.fromDockerCompose,
              });
            }}
            fetchServiceList={getProjectById}
            setLoading={() => {}}
          />
        </div>
      </div>
      <div className="fixed bottom-6 right-4 pl-10 pr-10">
        <Suspense fallback={<SpinnerIcon color="text-black" />}>
          <Version />
        </Suspense>
      </div>
    </div>
  );
}