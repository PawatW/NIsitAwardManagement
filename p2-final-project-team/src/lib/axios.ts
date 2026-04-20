import axios from "axios";

const baseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

const api = axios.create({
  baseURL: baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("nisit_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      if (typeof window !== "undefined") {
        localStorage.removeItem("nisit_token");
        window.location.href = "/"; 
      }
    }
    return Promise.reject(error);
  }
);

export default api;