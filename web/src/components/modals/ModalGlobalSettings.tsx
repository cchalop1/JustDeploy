import ServerList from "../ServerList/ServerList";
import { Suspense } from "react";
import ServerListSkeleton from "../ServerList/ServerListSkeleton";
import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";
import Modal from "./Modal";

type ModalGlobalSettingsProps = {
  onClose: () => void;
  onClickNewServer: () => void;
};

export default function ModalGlobalSettings({
  onClose,
  onClickNewServer,
}: ModalGlobalSettingsProps) {
  return (
    <Modal
      onClose={onClose}
      headerNode={<h1 className="text-2xl font-bold">Global Settings</h1>}
    >
      <div
        className="flex 
       flex-col h-[calc(100%-3rem)] pl-4 pr-4 pt-2 pb-2"
      >
        <div className="mt-4 mb-20">
          <div className="flex justify-between">
            <div className="text-2xl font-bold">Servers</div>
            <Button
              onClick={() => {
                onClose();
                onClickNewServer();
              }}
            >
              New Server
            </Button>
          </div>
          <Suspense fallback={<ServerListSkeleton />}>
            <ServerList />
          </Suspense>
        </div>
      </div>
    </Modal>
  );
}
