import { createBrowserRouter } from "react-router-dom";
import Home from "./Home";
import ServerConfigForm from "./components/forms/ServerConfigForm";
import { CreateDeployForm } from "./components/forms/CreateDeployForm";
import DeployPage from "./DeployPage";
import ServerPage from "./ServerPage";
import { Suspense } from "react";

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <Suspense fallback={"loading..."}>
        <Home />
      </Suspense>
    ),
  },
  {
    path: "server/create",
    element: <ServerConfigForm />,
  },
  {
    path: "deploy/create",
    element: <CreateDeployForm />,
  },
  {
    path: "deploy/:id",
    element: <DeployPage />,
  },
  {
    path: "server/:id",
    element: <ServerPage />,
  },
]);
