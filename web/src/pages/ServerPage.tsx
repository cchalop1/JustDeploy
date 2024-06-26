import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import SpinnerIcon from "@/assets/SpinnerIcon";
import Status from "@/components/ServerStatus";
import { ServerDto } from "@/services/getServerListApi";
import { getServerByIdApi } from "@/services/getServerById";
import DeployList from "@/components/DeployList";
import LinkIcon from "@/assets/linkIcon";
import { DeployDto } from "@/services/getDeployListApi";
import ServerButtons from "@/components/ServerButtons";
import { getDeployListByServerIdApi } from "@/services/getDeployListByServerIdApi";

export default function ServerPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [server, setServer] = useState<ServerDto | null>(null);
  const [deployList, setDeployList] = useState<Array<DeployDto>>([]);

  async function fetchServerById(id: string) {
    const resServer = await getServerByIdApi(id);
    setServer(resServer);
  }

  async function fetchDeployListByServerId(id: string) {
    const resDeploys = await getDeployListByServerIdApi(id);
    setDeployList(resDeploys);
  }

  useEffect(() => {
    if (!id) {
      navigate("/");
    } else {
      fetchServerById(id);
      fetchDeployListByServerId(id);
    }
  }, [id, navigate]);

  if (server === null) {
    return <SpinnerIcon color="text-black" />;
  }

  return (
    <div className="mt-40">
      <div className="flex justify-between">
        <div className="text-xl font-bold mb-2">{server.name}</div>
        <ServerButtons server={server} fetchServerById={fetchServerById} />
      </div>
      <Status status={server.status} />
      <div className="flex items-center mt-2 gap-2">
        <LinkIcon />
        <a href={server.ip} target="_blank" className="underline">
          {server.ip}
        </a>
        <a href={server.domain} target="_blank" className="underline">
          {server.domain}
        </a>
      </div>
      <div className="mt-20 ">
        <div className="text-2xl font-bold mb-4">Deploys on this server</div>
        <DeployList deployList={deployList} />
      </div>
    </div>
  );
}
