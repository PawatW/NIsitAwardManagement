import api from "../../../lib/axios";

export const getFaculties = () => {
  return api.get("/faculties");
};

export const getFaculty = (id: number) => {
  return api.get(`/faculties/${id}`);
};

export const createFaculty = (data: any) => {
  return api.post("/faculties", data);
};

export const updateFaculty = (id: number, data: any) => {
  return api.put(`/faculties/${id}`, data);
};

export const deleteFaculty = (id: number) => {
  return api.delete(`/faculties/${id}`);
};
