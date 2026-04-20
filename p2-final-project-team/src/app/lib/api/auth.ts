import api from "../../../lib/axios";

export const getGoogleAuthUrl = () => {
  return api.get("/api/v1/auth/google");
};

export const register = (data: any) => {
  return api.post("/api/v1/user/register", data);
};

export const getMe = () => {
  return api.get("/api/v1/user/me");
};
