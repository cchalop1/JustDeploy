import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";
import ServerConfigForm from "@/components/forms/ServerConfigForm";
import { CreateDeployForm } from "@/components/forms/CreateDeployForm";
import DeployPage from "@/pages/DeployPage";
import CreateServerLoading from "@/pages/CreateServerLoading";
import CreateDeployLoading from "@/pages/CreateDeployLoading";
import Layout from "@/Layout";
import { NotificationProvider } from "./contexts/Notifications";
import Background from "./components/ui/background";
import GithubRedirect from "./pages/GithubRedirect";
import Home from "./pages/Home";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Layout />}>
      <Route path="/github/auth/redirect" element={<GithubRedirect />} />
      {/* <Route path="server/create" element={<ServerConfigForm />} />
      <Route path="server/:id/installation" element={<CreateServerLoading />} />

      <Route path="deploy/:id/installation" element={<CreateDeployLoading />} />
      <Route path="deploy/create" element={<CreateDeployForm />} /> */}
      {/* <Route path="deploy/:id" element={<DeployPage />} /> */}
      <Route
        path="/"
        element={
          <NotificationProvider>
            <Background>
              <Home />
            </Background>
          </NotificationProvider>
        }
      />
    </Route>
  )
);

export default router;
