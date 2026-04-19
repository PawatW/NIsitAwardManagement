import api from "../../../lib/axios";

// ดึงรายการที่ต้อง review
export const getReviews = () => {
  return api.get("/api/v1/reviews");
};

// ดึง review ตาม ID
export const getReview = (id: number) => {
  return api.get(`/api/v1/reviews/${id}`);
};

// สร้าง review ใหม่
export const createReview = (data: any) => {
  return api.post("/api/v1/reviews", data);
};

// อัพเดท review
export const updateReview = (id: number, data: any) => {
  return api.put(`/api/v1/reviews/${id}`, data);
};

// อนุมัติใบสมัคร (ส่งต่อคณะกรรมการ)
export const approveApplication = (applicationId: number, reviewData: any) => {
  return api.post(`/api/v1/applications/${applicationId}/approve`, reviewData);
};

// ส่งกลับแก้ไข
export const rejectApplication = (applicationId: number, reviewData: any) => {
  return api.post(`/api/v1/applications/${applicationId}/reject`, reviewData);
};

// ดึง review history
export const getReviewHistory = (applicationId: number) => {
  return api.get(`/api/v1/applications/${applicationId}/reviews`);
};
