import CommandModal, {
  CreateServiceFunc,
} from "@/components/databaseServices/CommandModal";
import DatabaseServiceCard from "@/components/databaseServices/DatabaseServiceCard";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { createAppApi } from "@/services/createAppApi";
import { createServiceApi } from "@/services/createServiceApi";
import { getProjectByIdApi, ProjectDto } from "@/services/getProjectById";
import {
  getPreConfiguredServiceListApi,
  ServiceDto,
} from "@/services/getServicesApi";
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

  useEffect(() => {
    getProjectById();
    getServices();
  }, [id]);

  if (!project) return <div>Loading...</div>;

  return (
    <div className="p-6 bg-slate-100 h-screen">
      <div className="flex gap-3 items-center">
        <Card className="p-3 w-1/5">
          <div className="text-xl font-bold">{project.name}</div>
          <div>{project.path}</div>
        </Card>
        <div>
          <Button onClick={() => setOpen(true)}>Create +</Button>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-3 mt-3">
        {project.apps.map((app) => (
          <Card key={app.id} className="p-3">
            <div className="text-xl font-bold">{app.name}</div>
            <div>{app.path}</div>
          </Card>
        ))}
      </div>
      <div className="grid grid-cols-2 gap-3 mt-3">
        {project.services.map((service) => (
          <DatabaseServiceCard
            key={service.id}
            service={service}
            setServiceToDelete={(s) => {}}
          />
        ))}
        {loading && <Card className="p-3">Loading...</Card>}
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
