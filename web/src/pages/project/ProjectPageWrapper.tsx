import { useParams } from "react-router-dom";
import ProjectPage from "../ProjectPage";

export default function ProjectPageWrapper() {
  const { id } = useParams();

  if (!id) {
    return null;
  }

  return <ProjectPage id={id} />;
}
