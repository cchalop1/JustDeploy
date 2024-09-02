import { getLocalServices } from "@/services/getLocalServices";
import ServicesManagements from "./ServicesManagements";
import { useState } from "react";
import { Service } from "@/services/getServicesByDeployId";
import {
  CreateServiceApi,
  createServiceApi,
} from "@/services/createServiceApi";

export default function ServicesLocalContainer() {
  const [localServices, setLocalServices] = useState<Array<Service>>([]);

  async function fetchServiceList() {
    const res = await getLocalServices();
    setLocalServices(res);
  }

  async function createService(createServiceParams: CreateServiceApi) {
    await createServiceApi(createServiceParams);
  }

  return (
    <div className="mt-10 mb-20">
      <div className="text-2xl font-bold mb-2">Local Databases</div>
      <ServicesManagements
        services={localServices}
        createService={createService}
        fetchServiceList={fetchServiceList}
      />
    </div>
  );
}
