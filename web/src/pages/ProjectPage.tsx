import CommandModal from "@/components/databaseServices/CommandModal";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
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
      <div className="grid grid-cols-2 gap-3">
        {project.services.map((service) => (
          <Card key={service.id} className="p-3">
            <div className="text-lg font-bold">{service.name}</div>
            <div>{service.id}</div>
          </Card>
        ))}
        {loading && <Card className="p-3">Loading...</Card>}
      </div>
      {/* // TODO: move to a other components */}
      <CommandModal
        open={openModal}
        setOpen={setOpen}
        preConfiguredServices={preConfiguredServices}
        serviceFromDockerCompose={[]}
        currentPath={project.path}
        createService={async (serviceName, fromDockerCompose) => {
          setLoading(true);
          setOpen(false);
          await createServiceApi({
            serviceName,
            fromDockerCompose,
            projectId: project.id,
          });
          await getProjectById();
          setLoading(false);
        }}
      />
    </div>
  );
}
