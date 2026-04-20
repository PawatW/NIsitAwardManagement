import api from "../../../lib/axios";

export const getApplications = (params?: any) => {
  return api.get("/api/v1/requests", { params });
};

export const getApplication = (id: string) => {
  return api.get(`/api/v1/requests/${id}`);
};

export const getMyApplications = () => {
  return api.get("/api/v1/requests/my");
};

export const getApplicationsByStatus = (status: string, params?: any) => {
  return api.get(`/api/v1/requests/status/${status}`, { params });
};

export const approveApplication = (id: string, data?: any) => {
  return api.post(`/api/v1/requests/${id}/approve`, data);
};

export const rejectApplication = (id: string, data?: any) => {
  return api.post(`/api/v1/requests/${id}/reject`, data);
};
