import api from "../../../lib/axios";

export const getFacultiesByCampus = (campusId: string) => {
  return api.get(`/api/v1/organization/campuses/${campusId}/faculties`);
};

export const createFaculty = (data: any) => {
  return api.post("/api/v1/admin/faculties", data);
};
