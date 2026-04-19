import api from "../../../lib/axios";

export const getAwards = () => {
  return api.get("/awards");
};

export const getAward = (id: number) => {
  return api.get(`/awards/${id}`);
};

export const createAward = (data: any) => {
  return api.post("/awards", data);
};

export const updateAward = (id: number, data: any) => {
  return api.put(`/awards/${id}`, data);
};

export const deleteAward = (id: number) => {
  return api.delete(`/awards/${id}`);
};
