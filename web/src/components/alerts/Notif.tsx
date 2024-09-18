import { Notification } from "@/contexts/Notifications";

type NotifProps = {
  notification: Notification;
  removeNotif: (notif: Notification) => void;
};

export default function Notif({ notification, removeNotif }: NotifProps) {
  const iconPath = "/icons/" + notification.type + ".jpg";

  return (
    <div className="bg-white border rounded-lg p-4 flex gap-4 justify-between">
      <img src={iconPath} className="h-8" />
      <div className="flex flex-col items-start w-3/4">
        <div className="font-bold">{notification.title}</div>
        <div>{notification.content}</div>
      </div>
      <button
        className="rounded-full bg-green-50 w-8 h-8 flex justify-center items-center"
        onClick={() => removeNotif(notification)}
      >
        <img src="/icons/X.svg" className="h-4" />
      </button>
    </div>
  );
}
