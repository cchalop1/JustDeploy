import {
  Service,
  getServicesByDeployIdApi,
} from "@/services/getServicesByDeployId";
import AddService from "./AddServices";
import ServiceListDeploy from "./ServiceListDeploy";
import { useState } from "react";

type ServicesManagementsProps = {
  deployId: string;
};

export default function ServicesManagements({
  deployId,
}: ServicesManagementsProps) {
  const [services, setServices] = useState<Service[]>([]);
  const [loadingNewService, setLoadingNewService] = useState(false);

  async function fetchServiceList() {
    const res = await getServicesByDeployIdApi(deployId);
    setServices(res);
  }

  return (
    <>
      <AddService
        deployId={deployId}
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
