import { DeployDto } from "@/services/getDeployListApi";
import { Card, CardContent } from "./ui/card";
import Status from "./ServerStatus";
import FolderIcon from "@/assets/FolderIcon";
import { useNavigate } from "react-router-dom";

type DeployListProps = {
  deployList: Array<DeployDto>;
};

function cleanPathToSource(pathToSource: string): string {
  const arr = pathToSource.split("/");
  return arr.slice(arr.length - 3, arr.length - 1).join("/");
}

export default function DeployList({ deployList }: DeployListProps) {
  const navigate = useNavigate();

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
          onClick={() => navigate("deploy/" + deploy.id)}
        >
          <CardContent className="flex flex-col gap-2">
            <div className="flex justify-between">
              <div className="font-bold">{deploy.name}</div>
              <Status status={deploy.status} />
            </div>
            <div className="flex items-center gap-2">
              <FolderIcon /> {cleanPathToSource(deploy.pathToSource)}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
