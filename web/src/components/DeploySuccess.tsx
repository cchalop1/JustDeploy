// import { Button } from "./ui/button";
// import { removeApplicationApi } from "@/services/removeApplication";
// import { DeployConfigDto } from "@/services/getDeployConfig";
// import { ButtonStateEnum } from "@/lib/utils";
// import { useState } from "react";
// import SpinnerIcon from "@/assets/SpinnerIcon";
// import FileIcon from "@/assets/fileIcon";
// import { Card } from "./ui/card";
// import { Badge } from "./ui/badge";
// import LinkIcon from "@/assets/linkIcon";
// import ModalApplicationLogs from "./modals/ModalLogs";
// import { reDeployAppApi } from "@/services/reDeployApp";
// import { stopApplicationApi } from "@/services/stopApplication";
// import { startApplicationApi } from "@/services/startApplication";
// import { Checkbox } from "./ui/checkbox";
// import { Label } from "./ui/label";
// import { editDeployementApi } from "@/services/editDeploymentApi";

// type DeploySuccessProps = {
//   deployConfig: DeployConfigDto;
//   fetchCurrentConfigData: () => void;
// };

// export default function DeploySuccess({
//   deployConfig,
//   fetchCurrentConfigData,
// }: DeploySuccessProps) {
//   const [connectButtonState, setConnectButtonState] = useState<ButtonStateEnum>(
//     ButtonStateEnum.INIT
//   );
//   const [redeployButtonState, setReDeployButtonState] =
//     useState<ButtonStateEnum>(ButtonStateEnum.INIT);
//   const [stopStartButtonState, setStopStartButtonState] =
//     useState<ButtonStateEnum>(ButtonStateEnum.INIT);

//   const [openLogs, setOpenLogs] = useState(false);

//   // TODO: fix this condition should be manage in the parent compoents
//   if (!deployConfig.appConfig) return null;
//   console.log(deployConfig);

//   async function removeApplication() {
//     if (!deployConfig.appConfig) return null;
//     setConnectButtonState(ButtonStateEnum.PENDING);
//     try {
//       await removeApplicationApi(deployConfig.appConfig.name);
//       setConnectButtonState(ButtonStateEnum.SUCESS);
//       fetchCurrentConfigData();
//     } catch (e) {
//       console.error(e);
//     }
//   }

//   async function reDeployApplication() {
//     if (!deployConfig.appConfig) return null;
//     setReDeployButtonState(ButtonStateEnum.PENDING);
//     try {
//       await reDeployAppApi(deployConfig.appConfig.name);
//       setReDeployButtonState(ButtonStateEnum.SUCESS);
//     } catch (e) {
//       console.error(e);
//     }
//   }

//   async function startStopApplication() {
//     if (!deployConfig.appConfig) return null;
//     setStopStartButtonState(ButtonStateEnum.PENDING);
//     try {
//       {
//         deployConfig.appStatus === "Runing"
//           ? await stopApplicationApi(deployConfig.appConfig.name)
//           : await startApplicationApi(deployConfig.appConfig.name);
//       }
//       fetchCurrentConfigData();
//       setStopStartButtonState(ButtonStateEnum.SUCESS);
//     } catch (e) {
//       console.error(e);
//     }
//   }

//   async function onCheckDeployOnCommit(checked: boolean) {
//     if (!deployConfig.appConfig) return null;
//     try {
//       await editDeployementApi(
//         {
//           deployOnCommit: checked,
//         },
//         deployConfig.appConfig.name
//       );
//       fetchCurrentConfigData();
//     } catch (e) {
//       console.error(e);
//     }
//   }

//   return (
//     <>
//       <ModalApplicationLogs
//         appName={deployConfig.appConfig.name}
//         open={openLogs}
//         onOpenChange={setOpenLogs}
//       />
//       <Card className="w-1/2 p-4">
//         <div className="flex justify-between">
//           <div className="font-bold">{deployConfig.appConfig.name}</div>
//           <div className="flex gap-2">
//             <Button variant="destructive" onClick={removeApplication}>
//               {connectButtonState === ButtonStateEnum.PENDING ? (
//                 <SpinnerIcon color="text-white" />
//               ) : (
//                 "Delete"
//               )}
//             </Button>
//             <Button onClick={startStopApplication}>
//               {stopStartButtonState === ButtonStateEnum.PENDING ? (
//                 <SpinnerIcon color="text-white" />
//               ) : deployConfig.appStatus === "Runing" ? (
//                 "Stop"
//               ) : (
//                 "Start"
//               )}
//             </Button>
//             <Button variant="secondary" onClick={() => reDeployApplication()}>
//               {redeployButtonState === ButtonStateEnum.PENDING ? (
//                 <SpinnerIcon color="text-black" />
//               ) : (
//                 "Redeploy"
//               )}
//             </Button>
//           </div>
//         </div>
//         <Badge
//           className={
//             deployConfig.appStatus === "Runing" ? "bg-green-600" : "bg-red-600"
//           }
//         >
//           {deployConfig.appStatus}
//         </Badge>
//         <div className="flex items-center gap-2 mt-4">
//           <LinkIcon />
//           <a href={deployConfig.url} target="_blank" className="underline">
//             {deployConfig.url}
//           </a>
//         </div>
//         <div className="flex items-center gap-2 mt-4">
//           <FileIcon />
//           <span className="text-sm text-gray-500 dark:text-gray-400">
//             {deployConfig.appConfig.pathToSource}
//           </span>
//         </div>
//         <div className="mt-4 flex items-center space-x-2">
//           <Checkbox
//             id="deploy-on-commit"
//             name="deploy-on-commit"
//             checked={deployConfig.appConfig.deployOnCommit}
//             onCheckedChange={(state) => {
//               if (typeof state === "boolean") {
//                 onCheckDeployOnCommit(state);
//               }
//             }}
//           />
//           <Label
//             htmlFor="deploy-on-commit"
//             className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
//           >
//             Deploy on every commit on <code>main</code> branch
//           </Label>
//         </div>
//         <Button className="mt-4" onClick={() => setOpenLogs(!openLogs)}>
//           Logs
//         </Button>
//       </Card>
//     </>
//   );
// }
