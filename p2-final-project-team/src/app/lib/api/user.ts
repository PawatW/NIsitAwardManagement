import api from "../../../lib/axios";

export const getMe = () => {
  return api.get("/api/v1/user/me");
};
