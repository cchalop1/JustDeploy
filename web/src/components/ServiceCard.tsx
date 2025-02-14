import { Service } from "@/services/getServicesByDeployId";
import { Badge } from "./ui/badge";
import { GithubIcon } from "lucide-react";
import { CardIcon } from "./CardIcon";

type ServiceCardProps = {
  service: Service;
  onClick: () => void;
};

export default function ServiceCard({ service, onClick }: ServiceCardProps) {
  return (
    <div
      className={`relative w-80 h-36 bg-white border rounded shadow-lg hover:shadow-xl cursor-pointer`}
      onClick={onClick}
    >
      <div className="flex justify-between items-center m-4 ">
        <CardIcon service={service} />
        <div className="flex gap-3">
          <Badge>{service.status}</Badge>
        </div>
      </div>
      <div className="ml-4 mr-4 flex items-center">
        <div className="font-bold">{service.imageName}</div>
      </div>
      <div className="underline font-bold text-xl mt-2 ml-4 mr-4 mb-2">
        {service.hostName}
      </div>
    </div>
  );
}

// return (
//   <div
//     className={`relative w-80 h-36 bg-white border ${
//       isDevContainer ? "rounded-b-lg rounded-e-lg" : "rounded-lg pt-3"
//     }  shadow-lg hover:shadow-xl cursor-pointer`}
//     onClick={onClick}
//   >
//     {isDevContainer && (
//       <div className="absolute -top-5 -left-[0.80px] w-1/4 h-5 z-10 border-t border-l border-r bg-white rounded-t-lg"></div>
//     )}
//     <div className="flex justify-between items-center ml-4 mr-4 ">
//       <GithubIcon size={24} />
//       <div className="flex gap-3">
//         {!isDevContainer && <Badge>{service.status}</Badge>}
//       </div>
//     </div>
//     <div className="font-bold text-xl mt-2 ml-4 mr-4 mb-2">
//       {service.hostName}
//     </div>
//     {isDevContainer ? (
//       <div className="font-mono ml-4 mr-4 text-sm">{service.currentPath}</div>
//     ) : (
//       <div className="border-t h-2/6 flex items-center">
//         <div className="ml-4 mr-4 mt-2 flex items-center">
//           <img src={service.imageUrl} className="h-6" />
//           <div className="ml-2">{service.imageName}</div>
//         </div>
//       </div>
//     )}
//   </div>
// );
