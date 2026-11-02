import { useState } from "react";
import { Service } from "@/services/getServicesByDeployId";
import AddService from "./AddServices";
import ServiceListDeploy from "./ServiceListDeploy";
type Props={services:Service[];fetchServiceList:()=>Promise<void>};
export default function ServicesManagements({services,fetchServiceList}:Props){
 const [loading,setLoading]=useState(false);
 return <><AddService setLoading={setLoading} fetchServices={fetchServiceList}/><ServiceListDeploy services={services} loadingNewService={loading} fetchServiceList={fetchServiceList}/></>;
}
