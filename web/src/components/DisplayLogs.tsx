import { Logs } from "@/services/getApplicationLogs";
import dayjs from "dayjs";

type DisplayLogsProps = {
  logs: Array<Logs>;
};

export default function DisplayLogs({ logs }: DisplayLogsProps) {
  if (logs.length === 0) {
    return (
      <p className="text-xs text-gray-500 font-mono py-4 px-3">No logs yet.</p>
    );
  }

  return (
    <div className="flex flex-col">
      {logs.map((log, idx) => (
        <div
          key={idx}
          className="flex gap-3 px-3 py-0.5 hover:bg-gray-100 font-mono text-xs leading-5"
        >
          <span className="flex-shrink-0 text-gray-400 select-none">
            {dayjs(log.date).format("HH:mm:ss")}
          </span>
          <span className="text-gray-700 break-all">{log.message}</span>
        </div>
      ))}
    </div>
  );
}
