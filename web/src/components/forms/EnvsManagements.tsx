import { Label } from "../ui/label";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Env } from "@/services/postFormDetails";

type EnvsManagementsProps = {
  envs: Env[];
  setEnvs: (envs: Env[]) => void;
};

export default function EnvsManagements({
  envs,
  setEnvs,
}: EnvsManagementsProps) {
  const addNewEnv = () => {
    setEnvs([...envs, { name: "", value: "" }]);
  };

  const removeEnv = (idx: number) => {
    setEnvs(envs.filter((_, index) => index !== idx));
  };

  return (
    <>
      <Label>Env Variables</Label>
      <div className="flex flex-col gap-2">
        {envs.map((env, idx) => {
          const lastEnv = envs.length - 1 === idx;
          return (
            <div className="flex gap-4" key={idx}>
              <Input
                id="envName"
                name="envName"
                type="envName"
                placeholder="Env Name"
                autoComplete="off"
                value={env.name}
                onChange={(e) => {
                  const updatedEnvs = [...envs];
                  updatedEnvs[idx] = {
                    ...updatedEnvs[idx],
                    name: e.target.value,
                  };
                  setEnvs(updatedEnvs);
                }}
              />
              <Input
                id="envSecret"
                name="envSecret"
                type="envSecret"
                placeholder="Env Secret"
                autoComplete="off"
                value={env.value}
                onChange={(e) => {
                  const updatedEnvs = [...envs];
                  updatedEnvs[idx] = {
                    ...updatedEnvs[idx],
                    value: e.target.value,
                  };
                  setEnvs(updatedEnvs);
                }}
              />
              <Button
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  if (lastEnv) {
                    addNewEnv();
                  } else {
                    removeEnv(idx);
                  }
                }}
              >
                {lastEnv ? "+" : "-"}
              </Button>
            </div>
          );
        })}
      </div>
    </>
  );
}
