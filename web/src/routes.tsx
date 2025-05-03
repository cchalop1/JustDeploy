import {
  Route,
  createBrowserRouter,
  createRoutesFromElements,
} from "react-router-dom";
import Layout from "@/Layout";
import { NotificationProvider } from "./contexts/Notifications";
import Background from "./components/ui/background";
import GithubRedirect from "./pages/GithubRedirect";
import Home from "./pages/Home";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route element={<Layout />}>
      <Route path="/github/auth/redirect" element={<GithubRedirect />} />
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
