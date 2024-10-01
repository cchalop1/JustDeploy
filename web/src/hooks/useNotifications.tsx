import { useContext } from "react";
import {
  NotificationContext,
  NotificationContextProps,
} from "@/contexts/Notifications";

export const useNotification = (): NotificationContextProps => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotification must be used within a NotificationProvider"
    );
  }
  return context;
};
