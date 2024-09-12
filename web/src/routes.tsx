import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";
import Home from "@/pages/Home";
import ServerConfigForm from "@/components/forms/ServerConfigForm";
import { CreateDeployForm } from "@/components/forms/CreateDeployForm";
import DeployPage from "@/pages/DeployPage";
import ServerPage from "@/pages/ServerPage";
import { Suspense } from "react";
import CreateServerLoading from "@/pages/CreateServerLoading";
import CreateDeployLoading from "@/pages/CreateDeployLoading";
import Layout from "@/Layout";
import ProjectPageWrapper from "./pages/project/ProjectPageWrapper";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Layout />}>
      <Route
        path="/"
        element={
          <Suspense>
            <Home />
          </Suspense>
        }
      />
      <Route path="server/create" element={<ServerConfigForm />} />
      <Route path="server/:id/installation" element={<CreateServerLoading />} />
      <Route path="server/:id" element={<ServerPage />} />

      <Route path="deploy/:id/installation" element={<CreateDeployLoading />} />
      <Route path="deploy/create" element={<CreateDeployForm />} />
      <Route path="deploy/:id" element={<DeployPage />} />
      <Route path="project/:id" element={<ProjectPageWrapper />} />
    </Route>
  )
);

export default router;
