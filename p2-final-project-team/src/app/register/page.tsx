"use client";

import { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import api from "../../lib/axios"; 
import "./register.css";

interface OrganizationItem {
  ID: string;   
  Name: string; 
}

interface RegisterFormData {
  first_name: string;
  last_name: string;
  nisit_id: string;
  phone: string;
  campus_id: string;
  faculty_id: string;
  department_id: string;
}

function RegisterFormContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  
  // ดึง Token จาก URL (?register_token=...)
  const registerToken = searchParams.get("register_token");

  const [isLoading, setIsLoading] = useState(false);
  
  // State เก็บข้อมูลตัวเลือกจาก API
  const [campuses, setCampuses] = useState<OrganizationItem[]>([]);
  const [faculties, setFaculties] = useState<OrganizationItem[]>([]);
  const [departments, setDepartments] = useState<OrganizationItem[]>([]);

  const [formData, setFormData] = useState<RegisterFormData>({
    first_name: "",
    last_name: "",
    nisit_id: "",
    phone: "",
    campus_id: "",
    faculty_id: "",
    department_id: "",
  });

  useEffect(() => {
    if (!registerToken) {
      alert("ไม่พบ Register Token ข้อมูลอาจไม่สมบูรณ์ กรุณาเข้าสู่ระบบผ่าน SSO ใหม่อีกครั้ง");
      router.push("/");
      return;
    }

    const fetchCampuses = async () => {
      try {
        const res = await api.get("/api/v1/organization/campuses");
        setCampuses(res.data.data || res.data || []);
      } catch (error) {
        console.error("Failed to fetch campuses", error);
      }
    };
    fetchCampuses();
  }, [registerToken, router]);

  useEffect(() => {
    const fetchFaculties = async () => {
      if (!formData.campus_id) {
        setFaculties([]);
        return;
      }
      try {
        const res = await api.get(`/api/v1/organization/campuses/${formData.campus_id}/faculties`);
        setFaculties(res.data.data || res.data || []);
      } catch (error) {
        console.error("Failed to fetch faculties", error);
      }
    };
    fetchFaculties();
  }, [formData.campus_id]);

  useEffect(() => {
    const fetchDepartments = async () => {
      if (!formData.faculty_id) {
        setDepartments([]);
        return;
      }
      try {
        const res = await api.get(`/api/v1/organization/faculties/${formData.faculty_id}/departments`);
        setDepartments(res.data.data || res.data || []);
      } catch (error) {
        console.error("Failed to fetch departments", error);
      }
    };
    fetchDepartments();
  }, [formData.faculty_id]);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => {
      const newData = { ...prev, [name]: value };
      
      // ล้างค่าเมื่อเปลี่ยนระดับที่สูงกว่า
      if (name === "campus_id") {
        newData.faculty_id = "";
        newData.department_id = "";
      }
      if (name === "faculty_id") {
        newData.department_id = "";
      }
      return newData;
    });
  };

  // จัดการตอนกดปุ่มลงทะเบียน
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!registerToken) {
      alert("ไม่สามารถลงทะเบียนได้เนื่องจากไม่มี Token (กรุณาล็อกอินใหม่อีกครั้ง)");
      return;
    }

    setIsLoading(true);

    try {
   
      const payload = {
        first_name: formData.first_name,
        last_name: formData.last_name,
        nisit_id: formData.nisit_id,
        phone_number: formData.phone, 
        campus_id: formData.campus_id,
        faculty_id: formData.faculty_id,
        department_id: formData.department_id,
        register_token: registerToken
      };

      const response = await api.post("/api/v1/user/register", payload, {
        headers: {
          Authorization: `Bearer ${registerToken}`
        }
      });

      if (response.status === 200 || response.status === 201) {
        alert("ลงทะเบียนสำเร็จ! กรุณาเข้าสู่ระบบใหม่");
        router.push("/"); 
      }
    } catch (error: any) {
      console.error("Backend Error Details:", error.response?.data);
      const errorMsg = error.response?.data?.message || error.response?.data?.error || "กรุณาลองใหม่อีกครั้ง";
      alert(`เกิดข้อผิดพลาด: ${errorMsg}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="register-container">
      <div className="register-wrapper">
        <div className="register-card">
          <div className="register-header">
            <h1 className="register-title">ลงทะเบียนนิสิต</h1>
            <p className="register-subtitle">
              กรุณากรอกข้อมูลส่วนตัวให้ครบถ้วนเพื่อเข้าใช้งานระบบ
            </p>
          </div>

          <form onSubmit={handleSubmit} className="register-form">
            
            {/* ชื่อ - นามสกุล */}
            <div className="form-row">
              <div className="form-group">
                <label htmlFor="first_name" className="form-label">
                  ชื่อ <span className="required">*</span>
                </label>
                <input
                  type="text"
                  id="first_name"
                  name="first_name"
                  required
                  value={formData.first_name}
                  onChange={handleInputChange}
                  className="form-input"
                  placeholder="ชื่อ"
                />
              </div>

              <div className="form-group">
                <label htmlFor="last_name" className="form-label">
                  นามสกุล <span className="required">*</span>
                </label>
                <input
                  type="text"
                  id="last_name"
                  name="last_name"
                  required
                  value={formData.last_name}
                  onChange={handleInputChange}
                  className="form-input"
                  placeholder="นามสกุล"
                />
              </div>
            </div>

            {/*  รหัสนิสิต - เบอร์โทรศัพท์ */}
            <div className="form-row">
              <div className="form-group">
                <label htmlFor="nisit_id" className="form-label">
                  รหัสนิสิต <span className="required">*</span>
                </label>
                <input
                  type="text"
                  id="nisit_id"
                  name="nisit_id"
                  required
                  value={formData.nisit_id}
                  onChange={handleInputChange}
                  className="form-input"
                  placeholder="เช่น 6410001111"
                  pattern="[0-9]{10}"
                  maxLength={10}
                />
              </div>

              <div className="form-group">
                <label htmlFor="phone" className="form-label">
                  เบอร์โทรศัพท์ <span className="required">*</span>
                </label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  required
                  value={formData.phone}
                  onChange={handleInputChange}
                  className="form-input"
                  placeholder="0812345678"
                  pattern="[0-9]{10}"
                  maxLength={10}
                />
              </div>
            </div>

            {/* วิทยาเขต - คณะ */}
            <div className="form-row">
              <div className="form-group">
                <label htmlFor="campus_id" className="form-label">
                  วิทยาเขต <span className="required">*</span>
                </label>
                <select
                  id="campus_id"
                  name="campus_id"
                  required
                  value={formData.campus_id}
                  onChange={handleInputChange}
                  className="form-select"
                >
                  <option value="">-- เลือกวิทยาเขต --</option>
                  {campuses.map((campus) => (
                    <option key={campus.ID} value={campus.ID}>
                      {campus.Name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label htmlFor="faculty_id" className="form-label">
                  คณะ <span className="required">*</span>
                </label>
                <select
                  id="faculty_id"
                  name="faculty_id"
                  required
                  value={formData.faculty_id}
                  onChange={handleInputChange}
                  disabled={!formData.campus_id}
                  className="form-select"
                >
                  <option value="">
                    {formData.campus_id ? "-- เลือกคณะ --" : "-- กรุณาเลือกวิทยาเขตก่อน --"}
                  </option>
                  {faculties.map((faculty) => (
                    <option key={faculty.ID} value={faculty.ID}>
                      {faculty.Name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* ภาควิชา  */}
            <div style={{ display: "flex", justifyContent: "center" }}>
              <div className="form-group" style={{ width: "100%", maxWidth: "340px" }}>
                <label htmlFor="department_id" className="form-label">
                  ภาควิชา / สาขา <span className="required">*</span>
                </label>
                <select
                  id="department_id"
                  name="department_id"
                  required
                  value={formData.department_id}
                  onChange={handleInputChange}
                  disabled={!formData.faculty_id}
                  className="form-select"
                >
                  <option value="">
                    {formData.faculty_id ? "-- เลือกภาควิชา --" : "-- กรุณาเลือกคณะก่อน --"}
                  </option>
                  {departments.map((dept) => (
                    <option key={dept.ID} value={dept.ID}>
                      {dept.Name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* ปุ่มส่ง */}
            <div className="form-submit mt-2">
              <button
                type="submit"
                disabled={isLoading}
                className="submit-button"
              >
                {isLoading ? (
                  <span className="loading-text">
                    <svg className="spinner" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="spinner-circle" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="spinner-path" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    กำลังส่งข้อมูล...
                  </span>
                ) : (
                  "ลงทะเบียน"
                )}
              </button>
            </div>
            
          </form>
        </div>
      </div>
    </div>
  );
}

export default function RegisterPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <RegisterFormContent />
    </Suspense>
  );
}