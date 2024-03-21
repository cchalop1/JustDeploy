import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getDeployByIdApi } from "./services/getDeployById";
import { DeployDto } from "./services/getDeployListApi";
import SpinnerIcon from "./assets/SpinnerIcon";
import Status from "./components/ServerStatus";
import LinkIcon from "./assets/linkIcon";
import FolderIcon from "./assets/FolderIcon";
import DeployButtons from "./components/DeployButtons";

export default function DeployPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [deploy, setDeploy] = useState<DeployDto | null>({
    name: "fezljkn",
    deployOnCommit: false,
    email: "",
    enableTls: false,
    id: "toto",
    pathToSource: "fkné jkf",
    serverId: "",
    envs: [],
    status: "Installing",
    url: "",
  });

  async function fetchDeployById(id: string) {
    const resDeploy = await getDeployByIdApi(id);
    console.log(resDeploy);
    setDeploy(resDeploy);
  }

  useEffect(() => {
    if (!id) {
      navigate("/");
    } else {
      fetchDeployById(id);
    }
  }, [id, navigate]);

  if (deploy === null) {
    return <SpinnerIcon color="text-black" />;
  }
  // console.log(deploy);
  return deploy.name;
}
/*
  return (
    <div className="mt-40">
      <div className="flex justify-between">
        <div className="text-xl font-bold mb-2">{deploy.name}</div>
        <DeployButtons deploy={deploy} />
      </div>
      <Status status={deploy.status} />
      <div className="flex items-center mt-2 gap-2">
        <LinkIcon />
        <a href={deploy.url} target="_blank" className="underline">
          {deploy.url}
        </a>
      </div>
      <div className="flex items-center mt-2 gap-2">
        <FolderIcon />
        {deploy.pathToSource}
      </div>
    </div>
  );
*/
