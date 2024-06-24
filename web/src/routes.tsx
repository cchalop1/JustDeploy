import { createBrowserRouter } from "react-router-dom";
import Home from "./pages/Home";
import ServerConfigForm from "./components/forms/ServerConfigForm";
import { CreateDeployForm } from "./components/forms/CreateDeployForm";
import DeployPage from "./pages/DeployPage";
import ServerPage from "./pages/ServerPage";
import { Suspense } from "react";
import CreateServerStatus from "./pages/CreateServerStatus";

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
    path: "server/:id/installation",
    element: <CreateServerStatus />,
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
