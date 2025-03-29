"use client";
import axios from "axios";

const axiosClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || "https://your-api.com/api",
  timeout: 10000, // 10 seconds timeout
  headers: {
    "Content-Type": "application/json",
  },
});

// Request Interceptor
axiosClient.interceptors.request.use(
  (config) => {
    // Retrieve token from localStorage or cookies (if applicable)
    const auth_response =
      typeof window !== "undefined"
        ? JSON.parse(localStorage.getItem("auth-response"))
        : null;
 
    if (auth_response) {
      config.headers.Authorization = `Bearer ${auth_response?.token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor
axiosClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      if (error.response.status === 401) {
        // Handle Unauthorized access (e.g., redirect to login)
        console.log("Unauthorized: Redirecting to login...");
      }
    }
    return Promise.reject(error);
  }
);

export default axiosClient;
