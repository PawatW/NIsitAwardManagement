import api from "../../../lib/axios";

export const getRequests = () => {
  return api.get("/api/v1/requests");
};

export const getRequest = (id: number) => {
  return api.get(`/api/v1/requests/${id}`);
};

export const createRequest = (data: any) => {
  return api.post("/api/v1/requests", data);
};

export const updateRequest = (id: number, data: any) => {
  return api.put(`/api/v1/requests/${id}`, data);
};

export const updateRequestStatus = (id: number, status: string) => {
  return api.patch(`/api/v1/requests/${id}/status`, { status });
};

export const deleteRequest = (id: number) => {
  return api.delete(`/api/v1/requests/${id}`);
};
