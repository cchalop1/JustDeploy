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
import CreateServerLoading from "@/pages/CreateServerLoading";
import CreateDeployLoading from "@/pages/CreateDeployLoading";
import Layout from "@/Layout";
import ProjectPageWrapper from "./pages/project/ProjectPageWrapper";
import { NotificationProvider } from "./contexts/Notifications";
import Background from "./components/ui/background";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Layout />}>
      <Route path="/" element={<Home />} />
      <Route path="server/create" element={<ServerConfigForm />} />
      <Route path="server/:id/installation" element={<CreateServerLoading />} />
      <Route path="server/:id" element={<ServerPage />} />

      <Route path="deploy/:id/installation" element={<CreateDeployLoading />} />
      <Route path="deploy/create" element={<CreateDeployForm />} />
      <Route path="deploy/:id" element={<DeployPage />} />
      <Route
        path="project/:id"
        element={
          <NotificationProvider>
            <Background>
              <ProjectPageWrapper />
            </Background>
          </NotificationProvider>
        }
      />
    </Route>
  )
);

export default router;
