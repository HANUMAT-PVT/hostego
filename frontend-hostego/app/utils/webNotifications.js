
import axiosClient from "./axiosClient";

export const subscribeToNotifications = async (title, body) => {
    
  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey:
      "BGQRMk6dwGjrQHY47G4g1gphFGBdK11REbNsz8qUkMq9XJVkLO9VWs3a72ntetjKO5PRFEyRYrWggs8VJefqr7A",
  });

  // Send Subscription to Backend

  await axiosClient.post(`api/notifications`, {
    title,
    body,
    subscription: JSON.stringify(subscription),
  });
};
