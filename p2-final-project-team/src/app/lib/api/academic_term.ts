import api from "../../../lib/axios";

export const getAcademicTerms = () => {
  return api.get("/academic-terms");
};

export const getAcademicTerm = (id: number) => {
  return api.get(`/academic-terms/${id}`);
};

export const createAcademicTerm = (data: any) => {
  return api.post("/academic-terms", data);
};

export const updateAcademicTerm = (id: number, data: any) => {
  return api.put(`/academic-terms/${id}`, data);
};

export const deleteAcademicTerm = (id: number) => {
  return api.delete(`/academic-terms/${id}`);
};
