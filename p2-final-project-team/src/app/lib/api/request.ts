import api from "../../../lib/axios";

export const getRequests = (params?: any) => {
  return api.get("/api/v1/requests", { params });
};

export const getMyRequests = () => {
  return api.get("/api/v1/requests/my");
};

export const getRequest = (id: string) => {
  return api.get(`/api/v1/requests/${id}`);
};

export const createRequest = (data: any) => {
  return api.post("/api/v1/requests", data);
};

export const approveRequest = (id: string, data?: any) => {
  return api.post(`/api/v1/requests/${id}/approve`, data);
};

export const rejectRequest = (id: string, data?: any) => {
  return api.post(`/api/v1/requests/${id}/reject`, data);
};

export const changeAwardCategory = (id: string, data: any) => {
  return api.post(`/api/v1/requests/${id}/change-category`, data);
};

export const uploadDocument = (id: string, formData: FormData) => {
  return api.post(`/api/v1/requests/${id}/documents`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
};

export const getRequestsByStatus = (status: string, params?: any) => {
  return api.get(`/api/v1/requests/status/${status}`, { params });
};
