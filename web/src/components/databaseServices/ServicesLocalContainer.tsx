import { useState, useCallback } from "react";
import { Service } from "@/services/getServicesByDeployId";
import { getServicesApi } from "@/services/getServicesApi";
import ServicesManagements from "./ServicesManagements";
export default function ServicesLocalContainer(){
 const [services,setServices]=useState<Service[]>([]);
 const fetchServices=useCallback(async()=>{setServices(await getServicesApi())},[]);
 return <ServicesManagements services={services} fetchServiceList={fetchServices}/>;
}
