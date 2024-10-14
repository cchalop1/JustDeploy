import { XCircle } from "lucide-react";
import { PropsWithChildren, useRef } from "react";

type ModalProps = {
  headerNode?: React.ReactNode;
  onClose: () => void;
  className?: string;
};

export default function Modal({
  headerNode,
  onClose,
  children,
  className,
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
        "absolute right-5 top-24 border border-spacing-3 rounded-lg bg-white shadow-lg z-20 p-3 " +
        className
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
