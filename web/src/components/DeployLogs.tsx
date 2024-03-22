import SpinnerIcon from "@/assets/SpinnerIcon";
import { getApplicationLogs } from "@/services/getApplicationLogs";
import { useEffect, useState } from "react";

type DeployLogsProps = {
  id: string;
};

export default function DeployLogs({ id }: DeployLogsProps) {
  const [logs, setLogs] = useState<null | string[]>(null);

  useEffect(() => {
    getApplicationLogs(id).then(setLogs);
  }, []);

  return (
    <div className="flex flex-col ml-3 mr-3">
      {logs === null ? (
        <SpinnerIcon color="text-black" />
      ) : (
        logs.map((log, idx) => (
          <code key={idx} className=" text-xs w-full mb-2">
            {log}
          </code>
        ))
      )}
    </div>
  );
}
