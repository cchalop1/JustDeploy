import { useEffect, useState } from "react";
import DeployList from "./components/DeployList";
import ServerList from "./components/ServerList";
import { Button } from "./components/ui/button";
import { useNavigate } from "react-router-dom";
import { ServerDto, getServersListApi } from "./services/getServerListApi";
import { DeployDto, getDeployListApi } from "./services/getDeployListApi";

export default function Home() {
  const navigate = useNavigate();
  const [serverList, setServerList] = useState<Array<ServerDto>>([]);
  const [deployList, setDeployList] = useState<Array<DeployDto>>([]);

  const serverIsReady = serverList.length > 0;

  function onClickNewServer() {
    navigate("server/create");
  }

  function onClickNewDeploy() {
    navigate(`deploy/create`);
  }

  async function fetchServerList() {
    const serverList = await getServersListApi();
    // TODO: check error
    setServerList(serverList);
  }

  async function fetchDeployList() {
    const deployList = await getDeployListApi();
    // TODO: check error
    setDeployList(deployList);
  }

  useEffect(() => {
    fetchServerList();
    fetchDeployList();
    // setInterval(() => {
    //   fetchServerList();
    //   fetchDeployList();
    // }, 2000);
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
          <Button disabled={!serverIsReady} onClick={onClickNewDeploy}>
            Create Deploy
          </Button>
        </div>
        <DeployList deployList={deployList} />
      </div>
      <div className="h-52">
        <div className="flex justify-between">
          <div className="text-2xl font-bold">Databases</div>
        </div>
        {/* TODO: add database list */}
      </div>
    </>
  );
}
