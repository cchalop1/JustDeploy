import { useState, FormEvent } from "react";
import { Dialog, DialogContent, DialogTitle, DialogFooter, DialogHeader } from "../ui/dialog";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { addDomainToServerApi } from "@/services/addDomainToServerApi";
type Props = {open:boolean; onOpenChange:(open:boolean)=>void;serverId:string;fetchServerById:(id:string)=>void};
export default function ModalAddDnsSettings({open,onOpenChange,serverId,fetchServerById}:Props){
 const [error,setError]=useState("");const [saving,setSaving]=useState(false);
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();const domain=String(new FormData(event.currentTarget).get("server-domain") || "").trim();setSaving(true);setError("");
  try{await addDomainToServerApi({domain});fetchServerById(serverId);onOpenChange(false)}catch(e){setError(e instanceof Error ? e.message : "Save failed")}finally{setSaving(false)}
 }
 return <Dialog open={open} onOpenChange={onOpenChange}><DialogContent><form onSubmit={submit}><DialogHeader><DialogTitle>Connect a domain</DialogTitle><Label htmlFor="server-domain">Domain name</Label><Input id="server-domain" name="server-domain" required disabled={saving}/>{error && <p role="alert">{error}</p>}</DialogHeader><DialogFooter><Button type="submit" disabled={saving}>Save domain</Button></DialogFooter></form></DialogContent></Dialog>;
}
