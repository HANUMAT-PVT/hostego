"use client";

import { useEffect } from "react";
import axios from "../utils/axiosClient";

export function usePolling(userId, onMessage) {
  useEffect(() => {
    if (!userId) return;
    let intervalId;

    const fetchData = async () => {
      try {
        const response = await axios.get(
          `/events?user=${userId}&roles=super_admin,order_manager,order_assign_manager`
        );

        if (response?.data && response.data.message) {
          const jsonResp = JSON.parse(response.data.message);
          console.log("✅ Polling message received:", jsonResp);
          onMessage(jsonResp);
        } else {
          console.log("🔄 No new message");
        }
      } catch (error) {
        console.log("❌ Polling error:", error);
      }
    };

    // Start polling every 2 minutes
    intervalId = setInterval(fetchData, 120 * 1000);

    // Run once immediately
    fetchData();

    return () => clearInterval(intervalId);
  }, [userId, onMessage]);
}
