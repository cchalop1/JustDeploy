import CommandModal from "@/components/databaseServices/CommandModal";
import ServiceCard from "@/components/ServiceCard";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { AppDto } from "@/services/getProjectById";
import { ServiceDto } from "@/services/getServicesApi";
import { Service } from "@/services/getServicesByDeployId";
import { Folder, X } from "lucide-react";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { useQuery } from "@tanstack/react-query";
import { graphql } from "@/graphql";
import { execute } from "@/graphql/execute";

const GetProjectQuery = graphql(`
  query getProjectQuery($id: ID!) {
    getProject(id: $id) {
      id
      name
      path
      apps {
        id
        name
        path
      }
      services {
        id
        name
        status
        imageUrl
      }
    }
  }
`);

type ProjectPageProps = {
  id: string;
};

export default function ProjectPage({ id }: ProjectPageProps) {
  const [openModal, setOpen] = useState(false);
  const { data, isLoading, error } = useQuery({
    queryKey: ["id", id],
    queryFn: () => execute(GetProjectQuery, { id }),
  });
  const project = data?.getProject;
  console.log(project);
  console.log(isLoading);
  console.log(error);

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
              onClick={() => {}}
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
              onClick={() => {}}
            />
          ))}
        </div>
      </div>
      {/* Modal option */}
      {/* <div
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
      </div> */}
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
