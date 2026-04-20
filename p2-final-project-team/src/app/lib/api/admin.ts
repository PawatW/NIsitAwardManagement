import api from "../../../lib/axios";

export const getUsers = () => {
  return api.get("/api/v1/admin/users");
};

export const createUser = (data: any) => {
  return api.post("/api/v1/admin/users", data);
};

export const updateUserStatus = (id: string, isActive: boolean) => {
  return api.patch(`/api/v1/admin/users/${id}/status`, { is_active: isActive });
};

export const createCampus = (data: any) => {
  return api.post("/api/v1/admin/campuses", data);
};

export const createFaculty = (data: any) => {
  return api.post("/api/v1/admin/faculties", data);
};

export const createDepartment = (data: any) => {
  return api.post("/api/v1/admin/departments", data);
};
