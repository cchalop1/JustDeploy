import { Toaster } from "@/components/ui/toaster";
import { useToast } from "@/hooks/use-toast";
import React, { createContext, useState, ReactNode } from "react";

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
}

export const NotificationContext = createContext<
  NotificationContextProps | undefined
>(undefined);

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const { toast } = useToast();

  const addNotification = (
    type: NotificationType,
    title: string,
    content: string
  ) => {
    let variant = "default";
    if (type === "error") {
      variant = "destructive";
    }
    if (type === "warning") {
      title = `⚠️ ${title}`;
    }
    if (type === "success") {
      title = `✅ ${title}`;
    }
    toast({
      title,
      description: content,
      variant: variant as any,
    });
  };

  const success: NotificationFunc = ({ content, title }) =>
    addNotification("success", title, content);

  const error: NotificationFunc = ({ content, title }) =>
    addNotification("error", title, content);

  const info: NotificationFunc = ({ content, title }) =>
    addNotification("info", title, content);

  const warning: NotificationFunc = ({ content, title }) =>
    addNotification("warning", title, content);

  return (
    <NotificationContext.Provider value={{ success, error, info, warning }}>
      {children}
      <Toaster />
    </NotificationContext.Provider>
  );
};

// Hook pour utiliser le contexte
