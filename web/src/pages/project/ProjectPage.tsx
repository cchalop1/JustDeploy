import SpinnerIcon from "@/assets/SpinnerIcon";
import AlertModal from "@/components/alerts/AlertModal";
import AddService from "@/components/databaseServices/AddServices";
import { CreateServiceFunc } from "@/components/databaseServices/CommandModal";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceSideBar from "@/components/project/ServiceSideBar";
import ServiceCard from "@/components/ServiceCard";
import { createServiceApi } from "@/services/createServiceApi";
import { getProjectByIdApi, ProjectDto } from "@/services/getProjectById";
import { Service } from "@/services/getServicesByDeployId";
import { Folder } from "lucide-react";
import { useEffect, useState } from "react";

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
    <div className="p-6 bg-slate-100 h-screen">
      {displayAlert && (
        <AlertModal
          type="success"
          title="Service Created !"
          message="Your service has been created successfully, you can find its variables in .env"
        />
      )}
      <ProjectPageHeader />
      <div className="flex flex-col justify-center items-center h-3/5">
        <div>
          {apps.map((app) => (
            <ServiceCard
              key={app.id}
              Name={app.name}
              logo={<Folder className="h4" />}
              status="running"
              onClick={() => setServiceSelected(app)}
            />
          ))}
        </div>
        {!project && <SpinnerIcon color="text-black" />}
        <div className="flex gap-3 mt-3 ">
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              Name={service.name}
              logo={service.imageUrl}
              status="running"
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
      {serviceSelected && (
        <ServiceSideBar
          serviceSelected={serviceSelected}
          setServiceSelected={setServiceSelected}
          getProjectById={getProjectById}
        />
      )}
    </div>
  );
}
