import { Outlet } from "react-router-dom";
import Header from "@/components/Header";
import Footer from "@/components/Footer";

export default function Layout() {
  return (
    // <div className="flex justify-center">
    //   <div className="w-2/4 mt-10">

    <div className="">
      <div className="">
        {/* <Header /> */}
        <Outlet />
        {/* <Footer /> */}
      </div>
    </div>
  );
}
