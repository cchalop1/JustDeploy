import CommandModal, {
  CreateServiceFunc,
} from "@/components/databaseServices/CommandModal";
import DatabaseServiceCard from "@/components/databaseServices/DatabaseServiceCard";
import ServiceCard from "@/components/ServiceCard";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { createAppApi } from "@/services/createAppApi";
import { createServiceApi } from "@/services/createServiceApi";
import { deleteServiceByIdApi } from "@/services/deleteServiceApi";
import {
  AppDto,
  getProjectByIdApi,
  ProjectDto,
} from "@/services/getProjectById";
import {
  getPreConfiguredServiceListApi,
  ServiceDto,
} from "@/services/getServicesApi";
import { Service } from "@/services/getServicesByDeployId";
import { Cross, Folder, RemoveFormatting, X } from "lucide-react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function ProjectPage() {
  const { id } = useParams();
  const [openModal, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [project, setProject] = useState<ProjectDto | null>(null);
  const [preConfiguredServices, setPreConfiguredServices] = useState<
    Array<ServiceDto>
  >([]);
  const [serviceSelected, setServiceSelected] = useState<
    Service | AppDto | null
  >(null);

  async function getServices() {
    const res = await getPreConfiguredServiceListApi();
    setPreConfiguredServices(res);
  }

  async function getProjectById() {
    if (!id) return;
    const project = await getProjectByIdApi(id);
    setProject(project);
  }

  const create: CreateServiceFunc = async ({
    fromDockerCompose,
    path,
    serviceName,
  }) => {
    setLoading(true);
    setOpen(false);
    if (!project) return;
    if (path) {
      await createAppApi({
        path,
        projectId: project.id,
      });
    } else {
      await createServiceApi({
        serviceName,
        fromDockerCompose,
        projectId: project.id,
      });
    }
    await getProjectById();
    setLoading(false);
  };

  const deleteSelectedService = async () => {
    if (!serviceSelected) return;
    await deleteServiceByIdApi(serviceSelected.id);
    setServiceSelected(null);
    await getProjectById();
  };

  useEffect(() => {
    getProjectById();
    getServices();
  }, [id]);

  if (!project) return <div>Loading...</div>;

  return (
    <div className="p-6 bg-slate-100 h-screen">
      <div className="flex justify-between">
        <div className="font-bold text-3xl">JustDeploy 🛵</div>
        <div className="p-2 flex flex-row-reverse gap-3 items-center bg-white w-1/4 rounded-lg border shadow-lg">
          <Button onClick={() => setOpen(true)}>Deploy</Button>
          <Button variant="link" onClick={() => setOpen(true)}>
            Create +
          </Button>
        </div>
      </div>
      <div className="flex flex-col justify-center items-center h-3/6">
        <div className="grid grid-cols-2 gap-3 mt-3">
          {project.apps.map((app) => (
            <ServiceCard
              key={app.id}
              Name={app.name}
              logo={<Folder className="h4" />}
              status="running"
              onClick={() => setServiceSelected(app)}
            />
          ))}
        </div>
        <div className="grid grid-cols-2 gap-3 mt-3">
          {project.services.map((service) => (
            <ServiceCard
              key={service.id}
              Name={service.name}
              logo={service.imageUrl}
              status="running"
              onClick={() => setServiceSelected(service)}
            />
          ))}
          {loading && <Card className="p-3">Loading...</Card>}
        </div>
      </div>
      {/* Modal option */}
      <div
        hidden={!serviceSelected}
        className="absolute right-5 top-28 w-1/4 rounded-lg border shadow-lg h-full bg-white p-8"
      >
        <div className="flex justify-between">
          <div className="font-bold text-2xl">{serviceSelected?.name}</div>
          <X
            className="h-6 cursor-pointer"
            onClick={() => setServiceSelected(null)}
          />
        </div>

        <Button variant="destructive" onClick={() => deleteSelectedService()}>
          Delete
        </Button>
      </div>
      <CommandModal
        open={openModal}
        setOpen={setOpen}
        preConfiguredServices={preConfiguredServices}
        serviceFromDockerCompose={[]}
        currentPath={project.path}
        create={create}
      />
    </div>
  );
}
