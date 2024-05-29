import { DeployDto } from "@/services/getDeployListApi";
import EnvsManagements from "./forms/EnvsManagements";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import {
  EditDeployDto,
  editDeployementApi,
} from "@/services/editDeploymentApi";
import { Checkbox } from "./ui/checkbox";
import { Env } from "@/services/postFormDetails";

type DeploySettingsProps = {
  deploy: DeployDto;
  serverDomain: string;
  fetchDeployById: (id: string) => void;
};

export default function DeploySettings({
  deploy,
  serverDomain,
  fetchDeployById,
}: DeploySettingsProps) {
  const envs: Env[] =
    deploy.envs.length === 0 ? [{ name: "", value: "" }] : deploy.envs;
  const deploySetting: EditDeployDto = {
    deployOnCommit: deploy.deployOnCommit,
    envs: envs,
    id: deploy.id,
    subDomain: deploy.subDomain,
  };

  async function editDeploy(editDeployDto: EditDeployDto) {
    try {
      await editDeployementApi(editDeployDto);
      fetchDeployById(deploy.id);
    } catch (e) {
      console.error(e);
    }
  }

  return (
    <div>
      <div className="font-bold text-xl">Git Hooks</div>
      <div className="mt-4 flex items-center space-x-2">
        <Checkbox
          id="deploy-on-commit"
          name="deploy-on-commit"
          checked={deploy.deployOnCommit}
          onCheckedChange={(state) => {
            if (typeof state === "boolean") {
              editDeploy({ ...deploySetting, deployOnCommit: state });
            }
          }}
        />
        <Label
          htmlFor="deploy-on-commit"
          className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
        >
          Deploy on every commit on <code>main</code> branch
        </Label>
      </div>
      <div className="mt-4">
        <div className="font-bold text-xl">Domain</div>
        <div className="flex items-center gap-2">
          <div>{deploy.enableTls ? "https://" : "http://"}</div>
          <Input
            value={deploy.subDomain}
            onChange={(e) =>
              editDeploy({ ...deploySetting, subDomain: e.target.value })
            }
            type="text"
            id="subdomain"
            className="w-1/5"
          />
          <div>.{serverDomain}</div>
        </div>
      </div>
      <div className="font-bold text-xl mt-4">Env Variables</div>
      <EnvsManagements
        envs={deploySetting.envs}
        setEnvs={(newEnvs) => editDeploy({ ...deploySetting, envs: newEnvs })}
      />
    </div>
  );
}
