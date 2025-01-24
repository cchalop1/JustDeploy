import BtnConnectGithub from "@/components/BtnConnectGithub";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CheckCircle } from "lucide-react";
import { useEffect, useState } from "react";
import { getServerInfoApi } from "@/services/getServerInfoApi";

export default function Home() {
  const [serverIp, setServerIp] = useState<string>("");

  useEffect(() => {
    async function fetchServerInfo() {
      const serverInfo = await getServerInfoApi();
      setServerIp(serverInfo.ip);
    }
    fetchServerInfo();
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-center font-bold">
            JustDeploy 🛵
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            <p className="text-gray-600">
              Deploy your projects on your own servers with a single click.
            </p>
            <div className="flex items-center space-x-2 text-sm text-gray-500">
              <CheckCircle className={`h-4 w-4 text-green-500`} />
              <span>{"Ready to deploy"}</span>
            </div>
          </div>
          <BtnConnectGithub serverIp={serverIp} />
        </CardContent>
      </Card>
    </div>
  );
}
