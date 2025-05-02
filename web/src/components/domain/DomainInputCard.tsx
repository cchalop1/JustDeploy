import { useState } from "react";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Button } from "../ui/button";
import { InfoIcon, XCircle } from "lucide-react";
import { Card } from "../ui/card";

interface DomainInputCardProps {
  domain: string;
  serverIp: string;
  onDomainChange: (domain: string) => void;
  className?: string;
  label?: string;
  placeholder?: string;
}

export default function DomainInputCard({
  domain,
  serverIp,
  onDomainChange,
  className = "",
  label = "Domain",
  placeholder = "Enter domain (e.g., example.com)",
}: DomainInputCardProps) {
  const [showInfoCard, setShowInfoCard] = useState<boolean>(false);

  function handleDomainChange(value: string) {
    onDomainChange(value);
  }

  return (
    <div className={`mb-4 relative ${className}`}>
      <Label className="mb-1">{label}</Label>
      <div className="flex items-center">
        <Input
          value={domain}
          onChange={(e) => handleDomainChange(e.target.value)}
          placeholder={placeholder}
          className="mb-2 font-mono"
        />
        <Button
          type="button"
          variant="ghost"
          className="ml-2 p-2 h-10 mb-2"
          onClick={() => setShowInfoCard(!showInfoCard)}
        >
          <InfoIcon className="h-5 w-5 text-blue-600" />
        </Button>
      </div>

      {showInfoCard && (
        <Card className="p-3 bg-blue-50 border-blue-200 absolute z-10 right-0 mt-1 w-full shadow-lg rounded-md">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-2">
              <InfoIcon className="mt-1 h-4 w-4 flex-shrink-0 text-blue-600" />
              <div className="text-sm text-blue-800">
                <p className="font-medium mb-1">
                  How to configure your domain name:
                </p>
                <p>
                  Create a DNS record of type A that points to the server's IP
                  address:
                </p>
                <code className="bg-blue-100 px-2 py-1 rounded font-mono text-blue-900 block my-2">
                  {domain || "example.com"} → {serverIp || "0.0.0.0"}
                </code>
                <p>
                  This configuration will direct your domain to this server.
                </p>
                <p className="mt-2">
                  Additionally, create a wildcard CNAME record that points to
                  your domain:
                </p>
                <code className="bg-blue-100 px-2 py-1 rounded font-mono text-blue-900 block my-2">
                  *.{domain + "." || "example.com."} →{" "}
                  {domain + "." || "example.com."}
                </code>
                <p>
                  This allows subdomains to work with your server automatically.
                </p>
              </div>
            </div>
            <Button
              variant="ghost"
              className="p-1 h-6 w-6"
              onClick={() => setShowInfoCard(false)}
            >
              <XCircle className="h-5 w-5 text-blue-600" />
            </Button>
          </div>
        </Card>
      )}

      <p className="text-sm text-gray-600 mt-1">
        Enter the domain you want to use with your deployment server.
      </p>
    </div>
  );
}
