import { useNavigate } from "react-router-dom";

export default function Header() {
  const navigate = useNavigate();
  return (
    <div className="flex justify-between w-full">
      <div
        className="text-2xl font-bold cursor-pointer no-select"
        onClick={() => navigate("/")}
      >
        🛵 JustDeploy
      </div>
      <div className="flex gap-4 items-center">
        <a
          className="underline"
          href="https://github.com/cchalop1/JustDeploy?tab=readme-ov-file#-justdeploy"
          target="_blank"
        >
          Docs
        </a>
      </div>
    </div>
  );
}
