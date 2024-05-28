import { use } from "react";
import DeployList from "./components/DeployList";
import ServerList from "./components/ServerList";
import { Button } from "./components/ui/button";
import { useNavigate } from "react-router-dom";
import { getServersListApi } from "./services/getServerListApi";
import { getDeployListApi } from "./services/getDeployListApi";

export default function Home() {
  const navigate = useNavigate();
  const serverList = use(getServersListApi());
  const deployList = use(getDeployListApi());

  const serverIsReady = serverList.length > 0;

  function onClickNewServer() {
    navigate("server/create");
  }

  function onClickNewDeploy() {
    navigate(`deploy/create`);
  }

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
    </>
  );
}
