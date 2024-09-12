import { useNavigate, useParams } from "react-router-dom";
import ProjectPage from "./ProjectPage";

export default function ProjectPageWrapper() {
  const { id } = useParams();
  const navigate = useNavigate();

  if (!id) {
    navigate("/");
    return null;
  }

  return <ProjectPage id={id} />;
}
