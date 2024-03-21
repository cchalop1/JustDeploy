import { DeployDto } from "@/services/getDeployListApi";
import { Card, CardContent } from "./ui/card";
import Status from "./ServerStatus";

type DeployListProps = {
  deployList: Array<DeployDto>;
};

export default function DeployList({ deployList }: DeployListProps) {
  if (deployList.length === 0) {
    return (
      <div className="h-full flex justify-center pt-12">
        You have not created any deployments yet
      </div>
    );
  }

  return (
    <div className="flex gap-3 h-2/3 mt-2">
      {deployList.map((deploy) => (
        <Card
          className="hover:shadow-md cursor-pointer w-80 h-full pt-4 pl-2"
          key={deploy.id}
        >
          <CardContent className="flex flex-col gap-2">
            <div className="flex justify-between">
              <div className="font-bold">{deploy.name}</div>
              <Status status={deploy.status} />
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
