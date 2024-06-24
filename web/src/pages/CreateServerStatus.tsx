import useSubCreateServerEvent from "@/hooks/useSubCreateServerEvent";

export default function CreateServerStatus() {
  const events = useSubCreateServerEvent();

  return (
    <>
      <div>Server Is Creating..</div>
      <div>{events.map((e) => e)}</div>
    </>
  );
}
