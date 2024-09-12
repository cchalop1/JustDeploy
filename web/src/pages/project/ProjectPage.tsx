import SpinnerIcon from "@/assets/SpinnerIcon";
import AddService from "@/components/databaseServices/AddServices";
import CommandModal, {
  CreateServiceFunc,
} from "@/components/databaseServices/CommandModal";
import ProjectPageHeader from "@/components/project/ProjectPageHeader";
import ServiceSideBar from "@/components/project/ServiceSideBar";
import ServiceCard from "@/components/ServiceCard";
import { createAppApi } from "@/services/createAppApi";
import { createServiceApi } from "@/services/createServiceApi";
import { getProjectByIdApi, ProjectDto } from "@/services/getProjectById";
import {
  getPreConfiguredServiceListApi,
  ServiceDto,
} from "@/services/getServicesApi";
import { Service } from "@/services/getServicesByDeployId";
import { Folder } from "lucide-react";
import { useEffect, useState } from "react";

type ProjectPageProps = {
  id: string;
};

export default function ProjectPage({ id }: ProjectPageProps) {
  const [openModal, setOpen] = useState(false);
  const [project, setProject] = useState<ProjectDto | null>(null);
  const [serviceSelected, setServiceSelected] = useState<Service | null>(null);
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    ServiceDto[]
  >([]);
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
    setOpen(false);
    if (!project) return;
    await createServiceApi({
      serviceName,
      fromDockerCompose,
      projectId: project.id,
      localPath: path,
    });
    await getProjectById();
  };

  useEffect(() => {
    getProjectById();
    getPreConfiguredServiceListApi().then(setPreConfiguredServices);
  }, [id]);

  return (
    <div className="p-6 bg-slate-100 h-screen">
      <ProjectPageHeader setOpen={setOpen} />
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
