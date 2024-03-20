import Footer from "./components/Footer";
import Header from "./components/Header";
import { RouterProvider } from "react-router-dom";
import { router } from "./routes";

function App() {
  return (
    <div className="flex justify-center">
      <div className="w-2/3 mt-10">
        <Header />
        <RouterProvider router={router} />
        <Footer />
      </div>
    </div>
  );
}

export default App;
