import api from "../../../lib/axios";

export const getAwardCategories = () => {
  return api.get("/api/v1/award-categories");
};

export const getAwardCategory = (id: number) => {
  return api.get(`/api/v1/award-categories/${id}`);
};

export const createAwardCategory = (data: any) => {
  return api.post("/api/v1/award-categories", data);
};

export const updateAwardCategory = (id: number, data: any) => {
  return api.put(`/api/v1/award-categories/${id}`, data);
};

export const deleteAwardCategory = (id: number) => {
  return api.delete(`/api/v1/award-categories/${id}`);
};

export const getHonorRoll = (params?: any) => {
  return api.get("/api/v1/honor-roll", { params });
};
