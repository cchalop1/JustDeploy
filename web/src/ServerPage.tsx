import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import SpinnerIcon from "./assets/SpinnerIcon";
import Status from "./components/ServerStatus";
import { ServerDto } from "./services/getServerListApi";
import { getServerByIdApi } from "./services/getServerById";

export default function ServerPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [server, setServer] = useState<ServerDto | null>(null);

  async function fetchServerById(id: string) {
    const resServer = await getServerByIdApi(id);
    setServer(resServer);
  }

  useEffect(() => {
    if (!id) {
      navigate("/");
    } else {
      fetchServerById(id);
    }
  }, [id, navigate]);

  if (server === null) {
    return <SpinnerIcon color="text-black" />;
  }

  return (
    <div className="mt-40">
      <div className="flex justify-between">
        <div className="text-xl font-bold mb-2">{server.name}</div>
      </div>
      <Status status={server.status} />
    </div>
  );
}
