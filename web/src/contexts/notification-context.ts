import {createContext} from "react";
export type NotificationType = "success" | "error" | "info" | "warning";

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  content: string;
  onConfirm?: () => void;
}

type NotifFuncParams = {
  title: string;
  content: string;
  onConfirm?: () => void;
};

export type NotificationFunc = (params: NotifFuncParams) => void;

export interface NotificationContextProps {
  success: NotificationFunc;
  error: NotificationFunc;
  info: NotificationFunc;
  warning: NotificationFunc;
}

export const NotificationContext = createContext<
  NotificationContextProps | undefined
>(undefined);

