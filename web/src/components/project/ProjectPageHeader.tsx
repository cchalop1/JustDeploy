import { Button } from "../ui/button";

type ProjectPageHeaderProps = {
  setOpen: (open: boolean) => void;
};

export default function ProjectPageHeader({ setOpen }: ProjectPageHeaderProps) {
  return (
    <div className="flex justify-between">
      <div className="font-bold text-3xl">JustDeploy 🛵</div>
      <div className="p-2 flex flex-row-reverse gap-3 items-center bg-white w-1/4 rounded-lg border shadow-lg">
        <Button onClick={() => setOpen(true)}>Deploy</Button>
        {/* <Button variant="link" onClick={() => setOpen(true)}>
          Create +
        </Button> */}
      </div>
    </div>
  );
}
