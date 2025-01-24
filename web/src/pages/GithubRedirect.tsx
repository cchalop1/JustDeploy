import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";

export default function GithubRedirect() {
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const code = params.get("code");
    if (code) {
      //   navigate("/");
    }
  }, [location, navigate]);

  return <div></div>;
}
