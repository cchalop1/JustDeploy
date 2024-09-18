import { Notification } from "@/contexts/Notifications";
import Notif from "./Notif";

type DisplayNotificationsProps = {
  notifications: Notification[];
  removeNotif: (notif: Notification) => void;
};

export default function DisplayNotifications({
  notifications,
  removeNotif,
}: DisplayNotificationsProps) {
  return (
    <div className="absolute top-3 left-1/2 transform -translate-x-1/2 flex flex-col gap-3 w-1/3 z-20">
      {notifications.map((notification) => (
        <Notif
          key={notification.id}
          notification={notification}
          removeNotif={removeNotif}
        />
      ))}
    </div>
  );
}
