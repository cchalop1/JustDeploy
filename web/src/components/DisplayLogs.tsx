import { Logs } from "@/services/getApplicationLogs";
import dayjs from "dayjs";

type DisplayLogsProps = {
  logs: Array<Logs>;
};

export default function DisplayLogs({ logs }: DisplayLogsProps) {
  return (
    <div className="flex flex-col">
      {logs.map((log, idx) => (
        <code
          key={idx}
          className="flex gap-1 text-xs w-full mb-2 hover:bg-slate-50 rounded pt-1 pl-2 pr-2"
        >
          <div className="font-bold">
            {dayjs(log.date).format("DD/MM/YYYY - HH:mm:ss")}
          </div>
          {">"}
          <div>{log.message}</div>
        </code>
      ))}
    </div>
  );
}
