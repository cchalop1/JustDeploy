import { createBrowserRouter } from "react-router-dom";
import Home from "./Home";
import ServerConfigForm from "./components/forms/ServerConfigForm";
import { DeployConfigForm } from "./components/forms/DeployConfigForm";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "server/create",
    element: <ServerConfigForm />,
  },
  {
    path: "deploy/create",
    element: <DeployConfigForm />,
  },
]);
