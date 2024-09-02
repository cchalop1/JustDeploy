import { Service } from "@/services/getServicesByDeployId";
import AddService from "./AddServices";
import ServiceListDeploy from "./ServiceListDeploy";
import { useState } from "react";
import { CreateServiceApi } from "@/services/createServiceApi";

type ServicesManagementsProps = {
  deployId?: string;
  services: Service[];
  createService: (serviceParams: CreateServiceApi) => Promise<void>;
  fetchServiceList: (deployId?: string) => Promise<void>;
};

export default function ServicesManagements({
  deployId,
  services,
  createService,
  fetchServiceList,
}: ServicesManagementsProps) {
  const [loadingNewService, setLoadingNewService] = useState(false);

  return (
    <>
      <AddService
        deployId={deployId}
        createService={createService}
        setLoading={setLoadingNewService}
        fetchServiceList={fetchServiceList}
      />
      <ServiceListDeploy
        deployId={deployId}
        services={services}
        loadingNewService={loadingNewService}
        fetchServiceList={fetchServiceList}
      />
    </>
  );
}
