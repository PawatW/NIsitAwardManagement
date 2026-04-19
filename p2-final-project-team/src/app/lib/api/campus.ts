import api from "../../../lib/axios";

export const getCampuses = () => {
  return api.get("/campuses");
};

export const getCampus = (id: number) => {
  return api.get(`/campuses/${id}`);
};

export const createCampus = (data: any) => {
  return api.post("/campuses", data);
};

export const updateCampus = (id: number, data: any) => {
  return api.put(`/campuses/${id}`, data);
};

export const deleteCampus = (id: number) => {
  return api.delete(`/campuses/${id}`);
};
