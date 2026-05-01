import { X } from "lucide-react";
import Typewriter from "typewriter-effect";

type ModalWelcomeProps = {
  onClose: () => void;
};

export default function ModalWelcome({ onClose }: ModalWelcomeProps) {
  return (
    <div className="fixed w-72 bg-white border rounded-xl shadow-lg left-6 bottom-6 font-mono p-5 z-30">
      <button
        onClick={onClose}
        className="absolute top-3 right-3 text-gray-400 hover:text-gray-700 transition-colors"
      >
        <X className="w-4 h-4" />
      </button>
      <img src="/hand.png" className="w-12 mb-3" />
      <div className="font-bold text-sm mb-2">Welcome!</div>
      <div className="text-xs text-gray-600 leading-relaxed">
        <Typewriter
          onInit={(typewriter) => {
            typewriter
              .typeString("You're connected to your JustDeploy instance.")
              .pauseFor(500)
              .typeString("<br/>Docker is up and running.")
              .pauseFor(500)
              .typeString("<br/>You can now configure your domain, connect your GitHub account, and start deploying.")
              .start();
          }}
          options={{ delay: 25 }}
        />
      </div>
    </div>
  );
}
