import { useEffect, useState } from "react";
import Modal from "./Modal";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Switch } from "../ui/switch";
import { githubIsConnectedApi } from "@/services/githubIsConnected";
import { useNotification } from "@/hooks/useNotifications";
import { useInfo } from "@/hooks/useInfo";
import DomainInputCard from "../domain/DomainInputCard";
import { callApi, ResponseApi } from "@/services/api";
import { removeToken } from "@/services/authStorage";
export default function ModalGlobalSettings({onClose}:{onClose:()=>void}) {
 const {serverInfo,fetchServerInfo}=useInfo();
 const [domain,setDomain]=useState(serverInfo?.server.domain || "");
 const [email,setEmail]=useState(serverInfo?.server.email || "");
 const [useHttps,setUseHttps]=useState(serverInfo?.server.useHttps || false);
 const [connected,setConnected]=useState(false);const [saving,setSaving]=useState(false);
 const notif=useNotification();
 useEffect(()=>{let active=true;githubIsConnectedApi().then(({isConnected})=>{if(active)setConnected(isConnected)}).catch(()=>{});return()=>{active=false}},[]);
 async function save(){
  setSaving(true);
  try{await callApi<ResponseApi>("/server/settings","PUT",{domain,email,useHttps});await fetchServerInfo();notif.success({title:"Saved",content:"Server settings saved. Deploy exposed services to apply routing changes."})}
  catch(e){notif.error({title:"Save failed",content:e instanceof Error ? e.message : "Save failed"})}
  finally{setSaving(false)}
 }
 async function logout(){
  setSaving(true);
  try{await callApi<ResponseApi>("/logout","POST");removeToken();window.dispatchEvent(new Event("auth-changed"));onClose()}
  catch(e){notif.error({title:"Sign out failed",content:e instanceof Error ? e.message : "Sign out failed"})}
  finally{setSaving(false)}
 }
 return <Modal onClose={onClose} headerNode={<h1 className="font-semibold">Global settings</h1>} className="w-[520px]">
  <fieldset disabled={saving} className="p-5 space-y-5">
   <p>Server IP: <span className="font-mono">{serverInfo?.server.ip || "—"}</span></p>
   <DomainInputCard domain={domain} serverIp={serverInfo?.server.ip || ""} onDomainChange={setDomain}/>
   <div className="flex items-center justify-between"><Label htmlFor="https-switch">HTTPS for deployed services</Label><Switch id="https-switch" checked={useHttps} onCheckedChange={setUseHttps}/></div>
   <div><Label htmlFor="acme-email">Certificate email</Label><Input id="acme-email" type="email" value={email} onChange={(e)=>setEmail(e.target.value)}/></div>
   <p className="text-sm">GitHub {connected ? "connected" : "not connected"}</p>
   <div className="flex justify-between"><Button onClick={save}>{saving ? "Saving…" : "Save settings"}</Button><Button variant="outline" onClick={logout}>Sign out all sessions</Button></div>
  </fieldset>
 </Modal>;
}
