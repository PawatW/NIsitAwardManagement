import api from "../../../lib/axios";

export const getAdmins = () => {
  return api.get("/admins");
};

export const getAdmin = (id: number) => {
  return api.get(`/admins/${id}`);
};

export const createAdmin = (data: any) => {
  return api.post("/admins", data);
};

export const updateAdmin = (id: number, data: any) => {
  return api.put(`/admins/${id}`, data);
};

export const deleteAdmin = (id: number) => {
  return api.delete(`/admins/${id}`);
};
