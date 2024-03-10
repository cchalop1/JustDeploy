import { useEffect, useState } from "react";
import { AppConfigForm } from "./components/forms/DeployForm";
import {
  GetDeployConfigResponse,
  getDeployConfig,
} from "./services/getDeployConfig";
import ServerConfigForm from "./components/forms/ServerConfigForm";
import Steps from "./components/Steps";
import SpinnerIcon from "./assets/SpinnerIcon";
import DeploySuccess from "./components/DeploySuccess";

function App() {
  const [deployConfig, setDeployConfig] =
    useState<null | GetDeployConfigResponse>(null);
  function fetchCurrentConfigData() {
    getDeployConfig().then(setDeployConfig);
  }

  useEffect(() => {
    fetchCurrentConfigData();
  }, []);

  if (!deployConfig) return <SpinnerIcon color="text-black" />;

  return (
    <div>
      <div className="text-center m-5 text-4xl font-bold">JustDeploy</div>
      <Steps status={deployConfig.deployFormStatus} />
      <div className="flex w-full justify-center mt-10">
        {deployConfig.deployFormStatus === "serverconfig" && (
          <ServerConfigForm fetchCurrentConfigData={fetchCurrentConfigData} />
        )}
        {deployConfig.deployFormStatus === "appconfig" && (
          <AppConfigForm fetchCurrentConfigData={fetchCurrentConfigData} />
        )}
        {deployConfig.deployFormStatus === "deployapp" && (
          <DeploySuccess
            fetchCurrentConfigData={fetchCurrentConfigData}
            deployConfig={deployConfig}
          />
        )}
      </div>
    </div>
  );
}

export default App;
