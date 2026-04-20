import api from "../../../lib/axios";

export const getAcademicTerms = () => {
  return api.get("/api/v1/academic-terms");
};

export const getLatestAcademicTerm = () => {
  return api.get("/api/v1/academic-terms/latest");
};

export const getAcademicTerm = (id: number) => {
  return api.get(`/api/v1/academic-terms/${id}`);
};

export const createAcademicTerm = (data: any) => {
  return api.post("/api/v1/academic-terms", data);
};

export const updateAcademicTermStatus = (id: number, isOpen: boolean) => {
  return api.patch(`/api/v1/academic-terms/${id}`, { is_open: isOpen });
};
