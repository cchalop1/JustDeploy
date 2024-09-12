import { Button } from "../ui/button";

export default function ProjectPageHeader() {
  return (
    <div className="flex justify-between">
      <div className="font-bold text-3xl">JustDeploy 🛵</div>
      <div className="p-2 flex flex-row-reverse gap-3 items-center bg-white w-1/4 rounded-lg border shadow-lg">
        <Button onClick={() => {}}>Deploy</Button>
        {/* <Button variant="link" onClick={() => setOpen(true)}>
          Create +
        </Button> */}
      </div>
    </div>
  );
}
