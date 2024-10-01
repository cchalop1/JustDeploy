import DisplayNotifications from "@/components/alerts/DisplayNotifications";
import React, { createContext, useState, ReactNode } from "react";

const NOTIF_DISPLAY_TIME = 7000;

type NotificationType = "success" | "error" | "info" | "warning";

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

type NotificationFunc = (params: NotifFuncParams) => void;

export interface NotificationContextProps {
  success: NotificationFunc;
  error: NotificationFunc;
  info: NotificationFunc;
  warning: NotificationFunc;
  removeNotification: (id: string) => void;
}

export const NotificationContext = createContext<
  NotificationContextProps | undefined
>(undefined);

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const addNotification = (
    type: NotificationType,
    title: string,
    content: string,
    onConfirm?: () => void
  ) => {
    const id = Math.random().toString(36).substr(2, 9);
    setNotifications([
      ...notifications,
      { id, type, title, content, onConfirm },
    ]);
    setTimeout(() => {
      setNotifications((prevNotifications) => prevNotifications.slice(1));
    }, NOTIF_DISPLAY_TIME);
  };

  const removeNotification = (id: string) => {
    setNotifications(
      notifications.filter((notification) => notification.id !== id)
    );
  };

  const success: NotificationFunc = ({ content, title, onConfirm }) =>
    addNotification("success", title, content, onConfirm);

  const error: NotificationFunc = ({ content, title, onConfirm }) =>
    addNotification("error", title, content, onConfirm);

  const info: NotificationFunc = ({ content, title, onConfirm }) =>
    addNotification("info", title, content, onConfirm);

  const warning: NotificationFunc = ({ content, title, onConfirm }) =>
    addNotification("warning", title, content, onConfirm);

  return (
    <NotificationContext.Provider
      value={{ success, error, info, warning, removeNotification }}
    >
      <DisplayNotifications
        notifications={notifications}
        removeNotif={(notif) => removeNotification(notif.id)}
      />
      {children}
    </NotificationContext.Provider>
  );
};

// Hook pour utiliser le contexte
