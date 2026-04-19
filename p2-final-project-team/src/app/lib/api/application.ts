import api from "../../../lib/axios";

// ดึงรายการใบสมัครทั้งหมด
export const getApplications = (params?: any) => {
  return api.get("/api/v1/applications", { params });
};

// ดึงใบสมัครตาม ID
export const getApplication = (id: number) => {
  return api.get(`/api/v1/applications/${id}`);
};

// ดึงสถิติ Dashboard
export const getApplicationStats = () => {
  return api.get("/api/v1/applications/stats");
};

// ดึงใบสมัครที่รอพิจารณา
export const getPendingApplications = () => {
  return api.get("/api/v1/applications?status=pending");
};

// ดึงใบสมัครที่อนุมัติแล้ว
export const getApprovedApplications = () => {
  return api.get("/api/v1/applications?status=approved");
};

// ดึงใบสมัครที่ส่งกลับแก้ไข
export const getRejectedApplications = () => {
  return api.get("/api/v1/applications?status=rejected");
};

// อัพเดทสถานะใบสมัคร
export const updateApplicationStatus = (id: number, status: string, notes?: string) => {
  return api.patch(`/api/v1/applications/${id}/status`, { status, notes });
};

// ส่งต่อให้คณะกรรมการ
export const forwardToCommittee = (id: number) => {
  return api.post(`/api/v1/applications/${id}/forward`);
};

// ประกาศผลอย่างเป็นทางการ
export const publishResults = () => {
  return api.post("/api/v1/applications/publish");
};
