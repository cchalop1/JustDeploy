import { Service } from "@/services/getServicesByDeployId";
import { Badge } from "./ui/badge";
import { CardIcon } from "./CardIcon";

type ServiceCardProps = {
  service: Service;
  onClick: () => void;
};

export default function ServiceCard({ service, onClick }: ServiceCardProps) {
  return (
    <div
      className={`relative w-80 h-36 bg-white border rounded shadow-lg hover:shadow-xl cursor-pointer p-4 flex flex-col`}
      onClick={onClick}
    >
      <div className="flex justify-between items-center">
        {service.imageUrl ? (
          <img
            src={service.imageUrl}
            alt={service.name}
            className="h-8 w-8 object-contain"
          />
        ) : (
          <CardIcon service={service} />
        )}
        <div className="flex gap-3">
          <Badge
            variant={
              service.status === "Stopped"
                ? "destructive"
                : service.status === "Running"
                ? "outline"
                : "secondary"
            }
          >
            {service.status}
          </Badge>
        </div>
      </div>

      <div className="mt-2 flex items-center">
        <div className="font-bold">{service.imageName}</div>
      </div>

      <div className="mt-1 flex items-center">
        <Badge variant="outline" className="text-xs">
          {service.type === "database"
            ? "Database"
            : service.type === "github_repo"
            ? "GitHub Repo"
            : service.type === "llm"
            ? "LLM"
            : service.type}
        </Badge>
      </div>

      {service.url && (
        <a
          href={service.url}
          target="_blank"
          rel="noopener noreferrer"
          className="mt-auto"
        >
          <div className="underline text-sm text-blue-600 hover:text-blue-800">
            {service.url}
          </div>
        </a>
      )}
    </div>
  );
}
