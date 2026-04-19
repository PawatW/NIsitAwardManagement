import api from "../../../lib/axios";

export const getUsers = () => {
  return api.get("/api/v1/users");
};

export const getUser = (id: number) => {
  return api.get(`/api/v1/users/${id}`);
};

export const createUser = (data: any) => {
  return api.post("/api/v1/users", data);
};

export const updateUser = (id: number, data: any) => {
  return api.put(`/api/v1/users/${id}`, data);
};

export const deleteUser = (id: number) => {
  return api.delete(`/api/v1/users/${id}`);
};
