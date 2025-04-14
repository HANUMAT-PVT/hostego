"use client";

import { useEffect } from "react";
import axios from "../utils/axiosClient";

export function usePolling(userId, onMessage) {
  useEffect(() => {
    if (!userId) return;

    console.log("📡 Polling started for user:", userId);


    let intervalId;

    const fetchData = async () => {
      try {
        const response = await axios.get(`/events?user=${userId}`);

        if (response?.data && response.data.message) {
          console.log("✅ Polling message received:", response.data.message);
          onMessage(response.data.message);
        } else {
          console.log("🔄 No new message");
        }
      } catch (error) {
        console.error("❌ Polling error:", error.message);
      }
    };

    // Start polling every 5 seconds
    intervalId = setInterval(fetchData, 5000);

    // Run once immediately
    fetchData();

    return () => clearInterval(intervalId);
  }, [userId, onMessage]);
}
