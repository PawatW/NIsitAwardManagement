import api from "../../../lib/axios";

export const getRequestsForReview = (params?: any) => {
  return api.get("/api/v1/requests", { params });
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
