import { ConnectServerDto, connectServer } from "@/services/connectServer";
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
import ModalDnsSettings from "../modals/ModalDnsSettings";

type ServerConfigFormProps = {
  fetchCurrentConfigData: () => void;
};

export default function ServerConfigForm({
  fetchCurrentConfigData,
}: ServerConfigFormProps) {
  const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
    ButtonStateEnum.INIT
  );
  const [connectServerData, setConnectServerData] = useState<ConnectServerDto>({
    domain: "",
    password: null,
    sshKey: "",
    user: "root",
  });
  const [modalIsOpen, setModalIsOpen] = useState<boolean>(false);

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
    if (!modalIsOpen) {
      setModalIsOpen(true);
      return;
    }
    setModalIsOpen(false);

    if (!event.target) return;

    setConnectButtonState(ButtonStateEnum.PENDING);

    try {
      const res = await connectServer(connectServerData);
      console.log(res);
      setConnectButtonState(ButtonStateEnum.SUCESS);
      fetchCurrentConfigData();
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
    <div className="flex">
      <Card className="w-[500px] m-10">
        <CardHeader>
          <CardTitle>Connect To Your Server</CardTitle>
          <CardDescription>
            Connect your server with your domain name before deploy your
            application.
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
                <Label htmlFor="server-domain">Server Domain</Label>
                <Input
                  id="server-domain"
                  name="server-domain"
                  placeholder="Dns domain server ex: mydomain.com"
                  value={connectServerData.domain}
                  onChange={(e) =>
                    setConnectServerData({
                      ...connectServerData,
                      domain: e.target.value,
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
          <ModalDnsSettings
            onClick={handleSubmit}
            open={modalIsOpen}
            onOpenChange={setModalIsOpen}
            domain={connectServerData.domain}
          />
          <CardFooter className="flex justify-between">
            <Button type="submit" className="w-full">
              {connectButtonState === ButtonStateEnum.PENDING ? (
                <SpinnerIcon color="text-white" />
              ) : (
                "Connect and setup server"
              )}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}

<div className="pl-10 pr-10">
  <div className="mb-2">
    Before click to connect you server make sure you have corrcly setup your Dns
  </div>
</div>;
