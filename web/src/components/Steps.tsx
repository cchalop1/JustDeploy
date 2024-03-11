import DotsCircilsIcon from "@/assets/DotsCirlcleIcon";
import CheckIcon from "@/assets/checkIcon";
import { deployStatus } from "@/services/getDeployConfig";

type StepState = "pending" | "success" | "notstarted";

type StepProps = {
  state: StepState;
  title: string;
};

function Step({ state, title }: StepProps) {
  let icon = null;
  if (state === "notstarted" || state === "pending") {
    icon = <DotsCircilsIcon />;
  }
  if (state === "success") {
    icon = <CheckIcon size="10" />;
  }
  return (
    <div className="flex flex-col gap-1 justify-center items-center">
      {icon}
      <div>{title}</div>
    </div>
  );
}

type StepsProps = {
  status: deployStatus;
};

export default function Steps({ status }: StepsProps) {
  const toto: Array<StepState> = ["notstarted", "notstarted", "notstarted"];

  if (status === "appconfig") {
    toto[0] = "success";
  }
  if (status === "deployapp") {
    toto[0] = "success";
    toto[1] = "success";
    toto[2] = "success";
  }
  return (
    <div className="flex w-full justify-center mt-10">
      <div className="flex gap-2 w-1/3 justify-between">
        <Step state={toto[0]} title="Connect and Setup" />
        <Step state={toto[1]} title="Config App" />
        <Step state={toto[2]} title="Application Deploy" />
      </div>
    </div>
  );
}
