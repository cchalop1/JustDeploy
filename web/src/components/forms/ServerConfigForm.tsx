import { ConnectServerDto, connectServerApi } from "@/services/connectServer";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { useState } from "react";
import SpinnerIcon from "@/assets/SpinnerIcon";
import { ButtonStateEnum } from "@/lib/utils";
import { useNavigate } from "react-router-dom";

export default function ServerConfigForm() {
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [connectServerData, setConnectServerData] = useState<ConnectServerDto>({
    ip: "",
    password: null,
    sshKey: "",
    user: "root",
  });
  const navigate = useNavigate();

  async function readSSHKeyUpload(event: React.ChangeEvent<HTMLInputElement>) {
    const target = event.target;
    const files = target.files;

    if (!files || files.length === 0) {
      return;
    }

    const sshKeyFile = files[0];
    const fileContent = await sshKeyFile.text();
    setConnectServerData({ ...connectServerData, sshKey: fileContent });
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!event.target) return;

    setConnectButtonState(ButtonStateEnum.PENDING);

    try {
      const { id } = await connectServerApi(connectServerData);
      setConnectButtonState(ButtonStateEnum.SUCESS);
      navigate(`/server/${id}/installation`);
    } catch (e) {
      console.error(e);
    }
  };

  const changeAuthMethode = (value: string) => {
    if (value === "SSH key") {
      setConnectServerData({
        ...connectServerData,
        password: null,
        sshKey: "",
      });
    } else if (value === "password") {
      setConnectServerData({
        ...connectServerData,
        password: "",
        sshKey: null,
      });
    }
  };
  return (
    <div className="mt-16 flex justify-center">
      <Card className="w-[500px]">
        <CardHeader>
          <CardTitle>Connect Your Server</CardTitle>
          <CardDescription>
            Connect your server to deploy a application. The connect process can
            take a few minutes beacause we need to install some dependencies.
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent>
            <div className="grid w-full items-center gap-4">
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="server-user">Server user</Label>
                <Input
                  id="server-user"
                  name="server-user"
                  placeholder="root"
                  value={connectServerData.user}
                  disabled
                  onChange={(e) =>
                    setConnectServerData({
                      ...connectServerData,
                      user: e.target.value,
                    })
                  }
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="server-ip">Server Ip</Label>
                <Input
                  id="server-ip"
                  name="server-ip"
                  placeholder="Ip to your server"
                  value={connectServerData.ip}
                  onChange={(e) =>
                    setConnectServerData({
                      ...connectServerData,
                      ip: e.target.value,
                    })
                  }
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="auth-methode">Chose your auth methode</Label>
                <Select onValueChange={changeAuthMethode}>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="SSH key" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="SSH key">SSH key</SelectItem>
                    <SelectItem value="password">Password</SelectItem>
                    {/* <SelectItem value="arealy-access">Ssh connect</SelectItem> */}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex flex-col space-y-1.5">
                {connectServerData.sshKey !== null && (
                  <Input
                    type="file"
                    id="ssh-key"
                    name="ssh-key"
                    placeholder="Upload your ssh key"
                    onChange={readSSHKeyUpload}
                  ></Input>
                )}
                {connectServerData.password !== null && (
                  <Input
                    id="password"
                    name="password"
                    placeholder="Enter the password of your server"
                    type="password"
                    value={connectServerData.password}
                    onChange={(e) =>
                      setConnectServerData({
                        ...connectServerData,
                        password: e.target.value,
                      })
                    }
                  ></Input>
                )}
              </div>
            </div>
          </CardContent>
          <CardFooter className="flex justify-between">
            <Button type="submit" className="w-full">
              {connectButtonState === ButtonStateEnum.PENDING ? (
                <SpinnerIcon color="text-white" />
              ) : (
                "Connect server"
              )}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
