import api from "../../../lib/axios";

export const getCampuses = () => {
  return api.get("/api/v1/organization/campuses");
};

export const getFacultiesByCampus = (campusId: string) => {
  return api.get(`/api/v1/organization/campuses/${campusId}/faculties`);
};

export const getDepartmentsByFaculty = (facultyId: string) => {
  return api.get(`/api/v1/organization/faculties/${facultyId}/departments`);
};
