import { useEffect, useState } from "react";
import { Dialog, DialogContent, DialogTitle } from "../ui/dialog";
import DisplayLogs from "../DisplayLogs";
import { Logs } from "@/services/getApplicationLogs";
type Props={open:boolean;onOpenChange:(open:boolean)=>void;fetchServerProxyLogs:()=>Promise<Logs[]>};
export default function ModalLogs({open,onOpenChange,fetchServerProxyLogs}:Props){
 const [logs,setLogs]=useState<Logs[]>([]);const [error,setError]=useState("");
 useEffect(()=>{if(!open)return;let active=true;fetchServerProxyLogs().then((data)=>{if(active)setLogs(data)}).catch((e:unknown)=>{if(active)setError(e instanceof Error ? e.message : "Logs unavailable")});return()=>{active=false}},[open,fetchServerProxyLogs]);
 return <Dialog open={open} onOpenChange={onOpenChange}><DialogContent><DialogTitle>Logs</DialogTitle>{error ? <p role="alert">{error}</p> : <div className="h-64 overflow-scroll"><DisplayLogs logs={logs}/></div>}</DialogContent></Dialog>;
}
