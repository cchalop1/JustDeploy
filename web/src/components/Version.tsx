import SpinnerIcon from "@/assets/SpinnerIcon";
import { getVersionApi, VersionDto } from "@/services/getVersion";
import { useEffect, useState } from "react";

export default function Version() {
  const [version, setVersion] = useState<VersionDto | null>(null);

  useEffect(() => {
    getVersionApi().then(setVersion);
  }, []);

  return (
    <div className="flex gap-2 items-center font-normal border rounded-full shadow-lg bg-white pt-2 pb-2 pl-4 pr-4">
      {version === null ? (
        <SpinnerIcon color="text-black" />
      ) : (
        <>
          <div className="w-2 h-2 bg-green-500 rounded-full"></div>
          <a target="_blank" href={version.githubUrl}>
            <div className="underline">{version.tagName}</div>
          </a>
        </>
      )}
    </div>
  );
}
