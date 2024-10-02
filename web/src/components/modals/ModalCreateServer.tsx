import ServerConfigForm from "../forms/ServerConfigForm";
import Modal from "./Modal";

type ModalCreateServerProps = {
  onClose: () => void;
};

export default function ModalCreateServerModalCreateServer({
  onClose,
}: ModalCreateServerProps) {
  return (
    <Modal onClose={onClose}>
      <ServerConfigForm onClose={onClose} />
    </Modal>
  );
}
