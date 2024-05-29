export default function Header() {
  return (
    <div className="flex justify-between w-full">
      <div className="text-2xl font-bold cursor-pointer no-select">
        🛵 JustDeploy
      </div>
      <div className="flex gap-4 items-center">
        <a
          className="underline"
          href="https://github.com/cchalop1/JustDeploy"
          target="_blank"
        >
          Docs
        </a>
      </div>
    </div>
  );
}
