import { Service } from "@/services/getServicesByDeployId";
import { EllipsisVertical } from "lucide-react";

type ServiceCardProps = {
  service: Service;
  onClick: () => void;
};

export default function ServiceCard({ service, onClick }: ServiceCardProps) {
  const isDevContainer = service.isDevContainer;
  const iconPath = service.isDevContainer
    ? "/icons/folder.png"
    : "/icons/service.png";

  return (
    <div
      className={`relative w-80 h-36 bg-white border ${
        isDevContainer ? "rounded-b-lg rounded-e-lg" : "rounded-lg pt-3"
      }  shadow-lg hover:shadow-xl cursor-pointer`}
      onClick={onClick}
    >
      {isDevContainer && (
        <div className="absolute -top-5 -left-[0.80px] w-1/4 h-5 z-10 border-t border-l border-r bg-white rounded-t-lg"></div>
      )}
      <div className="flex justify-between ml-4 mr-4 items-center">
        <img src={iconPath} className="w-8" />
        <div className="flex gap-3">
          <div>tag</div>
          {!isDevContainer && <EllipsisVertical />}
        </div>
      </div>
      <div className="font-bold text-xl mt-2 ml-4 mr-4 mb-2">
        {service.name}
      </div>
      {isDevContainer ? (
        <div className="font-mono ml-4 mr-4 text-sm">{service.currentPath}</div>
      ) : (
        <div className="border-t h-2/6 flex items-center">
          <div className="ml-4 mr-4 mt-2 flex items-center">
            <img src={service.imageUrl} className="h-6" />
            <div className="ml-2">{service.imageName}</div>
          </div>
        </div>
      )}
    </div>
  );
}
