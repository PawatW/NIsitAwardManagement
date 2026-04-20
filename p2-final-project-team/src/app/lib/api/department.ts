import api from "../../../lib/axios";

export const getDepartmentsByFaculty = (facultyId: string) => {
  return api.get(`/api/v1/organization/faculties/${facultyId}/departments`);
};

export const createDepartment = (data: any) => {
  return api.post("/api/v1/admin/departments", data);
};
