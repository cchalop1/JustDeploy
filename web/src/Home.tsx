import { useEffect, useState } from "react";
import DeployList from "./components/DeployList";
import ServerList from "./components/ServerList";
import { Button } from "./components/ui/button";
import { useNavigate } from "react-router-dom";
import { ServerDto, getServersListApi } from "./services/getServerListApi";

export default function Home() {
  const navigate = useNavigate();
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const [deployList, setDeployList] = useState([]);

  const canCreateDeploy = serverList.length > 0;

  function onClickNewServer() {
    navigate("server/create");
  }

  async function fetchServerList() {
    const serverList = await getServersListApi();
    // TODO: check error
    setServerList(serverList);
  }

  async function fetchDeployList() {}

  useEffect(() => {
    fetchServerList();
    fetchDeployList();
  }, []);

  return (
    <>
      <div className="h-52 mt-40">
        <div className="flex justify-between">
          <div className="text-2xl font-bold">Servers</div>
          <Button onClick={onClickNewServer}>New Server</Button>
        </div>
        <ServerList serverList={serverList} />
      </div>
      <div className="h-52">
        <div className="flex justify-between">
          <div className="text-2xl font-bold">Deploys</div>
          <Button disabled={!canCreateDeploy}>Create Deploy</Button>
        </div>
        <DeployList deployList={deployList} />
      </div>
    </>
  );
}
