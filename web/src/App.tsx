import Footer from "@/components/Footer";
import Header from "@/components/Header";
import router from "@/routes";
import { RouterProvider } from "react-router-dom";

function App() {
  return <RouterProvider router={router} />;
}

export default App;
