import SpinnerIcon from "@/assets/SpinnerIcon";
import { useInfo } from "@/hooks/useInfo";

export default function Version() {
  const { serverInfo } = useInfo();
  const version = serverInfo?.version;

  return (
    <div className="flex gap-2 items-center font-normal border rounded-full shadow-lg bg-white pt-2 pb-2 pl-4 pr-4">
      {version === null ? (
        <SpinnerIcon color="text-black" />
      ) : (
        <>
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <a target="_blank" href={version?.githubUrl}>
            <div className="underline">{version?.tagName}</div>
          </a>
        </>
      )}
    </div>
  );
}
