import { XCircle } from "lucide-react";
import { PropsWithChildren, useEffect, useRef } from "react";

type ModalProps = {
  headerNode?: React.ReactNode;
  onClose: () => void;
};

export default function Modal({
  headerNode,
  onClose,
  children,
}: PropsWithChildren<ModalProps>) {
  const modalRef = useRef<HTMLDivElement>(null);

  // useEffect(() => {
  //   function handleClickOutside(event: MouseEvent) {
  //     if (
  //       modalRef.current &&
  //       !modalRef.current.contains(event.target as Node)
  //     ) {
  //       onClose();
  //     }
  //   }
  //   document.addEventListener("mousedown", handleClickOutside);
  //   return () => {
  //     document.removeEventListener("mousedown", handleClickOutside);
  //   };
  // }, [onClose]);

  return (
    <div
      className={
        "absolute right-5 top-24 border border-spacing-3 rounded-lg bg-white shadow-lg z-20 p-3"
      }
      ref={modalRef}
    >
      {headerNode && (
        <div className="flex justify-between p-3 bg-white rounded-t-lg items-center">
          {headerNode}
          <XCircle
            className="w-6 h-6 cursor-pointer"
            onClick={() => onClose()}
          />
        </div>
      )}
      {children}
    </div>
  );
}
