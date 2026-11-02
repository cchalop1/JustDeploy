import { useEffect, useState } from "react";
import { getApplicationLogs, Logs } from "@/services/getApplicationLogs";
import DisplayLogs from "./DisplayLogs";
export default function DeployLogs({id}:{id:string}){
 const [logs,setLogs]=useState<Logs[]>([]);const [error,setError]=useState("");
 useEffect(()=>{let active=true;getApplicationLogs(id).then((data)=>{if(active)setLogs(data)}).catch((e:unknown)=>{if(active)setError(e instanceof Error ? e.message : "Logs unavailable")});return()=>{active=false}},[id]);
 return error ? <p role="alert">{error}</p> : <DisplayLogs logs={logs}/>;
}
