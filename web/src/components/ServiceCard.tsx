import React from "react";
import { Card } from "./ui/card";
import Status from "./ServerStatus";

type ServiceCardProps = {
  logo: React.ReactNode | string;
  Name: string;
  status: "running" | "loading" | "stopped";
  onClick: () => void;
};

export default function ServiceCard(props: ServiceCardProps) {
  return (
    <Card
      className="hover:shadow-md cursor-pointer pt-6 pb-6 pl-5 pr-5 flex gap-6 w-80 h-28 align-top"
      onClick={props.onClick}
    >
      {typeof props.logo === "string" ? (
        <img className="h-4" src={props.logo} />
      ) : (
        props.logo
      )}
      <div className="">
        <div className="font-bold">{props.Name}</div>
        <Status status={"Runing"} />
      </div>
    </Card>
  );
}
