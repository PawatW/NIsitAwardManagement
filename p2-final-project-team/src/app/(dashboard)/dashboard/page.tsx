"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "../../../context/AuthContext"; 
import api from "../../../lib/axios"; 
import { 
  BookOpen, CheckCircle2, Loader2, XCircle, Info, 
  Award, ChevronRight, CalendarX2
} from "lucide-react";

interface AwardCategory {
  id: number;
  name: string;
  description: string;
  form_structure: any[]; 
  is_active: boolean;
}

interface ApplicationData {
  id: string;
  category: string;
  submittedDate: string;
  status: string;
  currentStep: number;
  isRejected: boolean;
  isAwardChanged: boolean;
}

export default function StudentDashboard() {
  const { user, loading: authLoading } = useAuth(); 
  
  const [categories, setCategories] = useState<AwardCategory[]>([]); 
  const [application, setApplication] = useState<ApplicationData | null>(null);
  
  // 🔥 State สำหรับเช็คว่ามีเทอมเปิดอยู่ไหม
  const [hasOpenTerm, setHasOpenTerm] = useState<boolean>(false);
  const [loading, setLoading] = useState(true); 

  const steps = [
    { label: "ส่งแล้ว" }, 
    { label: "หัวหน้าภาควิชา" }, 
    { label: "รองคณบดี" },
    { label: "คณบดี" }, 
    { label: "คณะกรรมการ" },
    { label: "อนุมัติแล้ว" },
  ];

  const mapStatusToStep = (status: string): number => {
    const s = status.toUpperCase();

    if (s.includes("REJECTED_BY_HOD")) return 1;
    if (s.includes("REJECTED_BY_VICE_DEAN")) return 2;
    if (s.includes("REJECTED_BY_DEAN")) return 3;
    if (s.includes("REJECTED_BY_COMMITTEE")) return 4;

    if (s === "SUBMITTED") return 0;
    if (s === "PENDING_HOD") return 1;
    if (s === "PENDING_VICE_DEAN") return 2;
    if (s === "PENDING_DEAN") return 3;
    if (s === "PENDING_COMMITTEE" || s === "AWARD_CHANGED") return 4;
    
    if (s === "APPROVED_PENDING_PRESIDENT" || s === "APPROVED") return 5;

    return 0;
  };

  const translateStatus = (status: string): string => {
    const s = status.toUpperCase();
    
    switch (s) {
      case "SUBMITTED": return "ส่งใบสมัครแล้ว รอการตรวจสอบ";
      case "PENDING_HOD": return "รอหัวหน้าภาควิชาพิจารณา";
      case "PENDING_VICE_DEAN": return "รอรองคณบดีพิจารณา";
      case "PENDING_DEAN": return "รอคณบดีพิจารณา";
      case "PENDING_COMMITTEE": return "รอคณะกรรมการพิจารณา";
      case "AWARD_CHANGED": return "ถูกเปลี่ยนประเภทรางวัล (รอพิจารณาใหม่)";
      case "APPROVED_PENDING_PRESIDENT": return "ผ่านการคัดเลือก รอประธานอนุมัติ";
      case "APPROVED": return "อนุมัติแล้ว (ได้รับรางวัล)";
      
      case "REJECTED_BY_HOD": return "ไม่ผ่านการอนุมัติ (โดยหัวหน้าภาควิชา)";
      case "REJECTED_BY_VICE_DEAN": return "ไม่ผ่านการอนุมัติ (โดยรองคณบดี)";
      case "REJECTED_BY_DEAN": return "ไม่ผ่านการอนุมัติ (โดยคณบดี)";
      case "REJECTED_BY_COMMITTEE": return "ไม่ผ่านการอนุมัติ (โดยคณะกรรมการ)";
      
      default: return status.replace(/_/g, ' '); // ถ้าไม่ตรงเงื่อนไขเลย ให้คืนค่าเดิมที่เอาขีดล่างออก
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      if (!user) return;
      try {
        setLoading(true);
        
        const [termRes, catRes, appRes] = await Promise.all([
          api.get("/api/v1/academic-terms"),
          api.get("/api/v1/award-categories"), 
          api.get("/api/v1/requests/my")
        ]);

        const termsData = termRes.data.data || termRes.data || [];
        const isOpen = termsData.some((term: any) => term.is_open === true);
        setHasOpenTerm(isOpen);

        setCategories((catRes.data.data || catRes.data || []).filter((c: AwardCategory) => c.is_active));

        const requestsList = appRes?.data?.data || appRes?.data || [];

        if (requestsList && requestsList.length > 0) {
          const b = requestsList[0]; 
          const rawStatus = b.current_status || "SUBMITTED";
          
          setApplication({
            id: b.id.toString(),
            category: b.award_name || "Award Application", 
            submittedDate: new Date(b.created_at).toLocaleDateString("en-US", {
              year: 'numeric', month: 'short', day: 'numeric'
            }),
            status: rawStatus, 
            currentStep: mapStatusToStep(rawStatus),
            isRejected: rawStatus.toUpperCase().includes("REJECT"),
            isAwardChanged: rawStatus.toUpperCase() === "AWARD_CHANGED"
          });
        }
      } catch (error) {
        console.error("Data fetch error:", error);
      } finally {
        setLoading(false);
      }
    };

    if (!authLoading && user) fetchData();
  }, [user, authLoading]);

  if (authLoading || loading) {
    return (
      <div className="flex flex-col items-center justify-center h-[60vh]">
        <Loader2 className="w-10 h-10 animate-spin text-green-600 mb-4" />
        <p className="text-gray-500 font-bold uppercase tracking-widest text-xs">Synchronizing with Server</p>
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-20 px-4 max-w-7xl mx-auto">
      <header>
        <h1 className="text-3xl font-black text-gray-900 leading-tight">Overview</h1>
        <p className="text-gray-500 font-medium italic">ยินดีต้อนรับ, {user?.email}</p>
      </header>

      {/* Progress Tracking Section  */}
      <section>
        <h2 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">ติดตามความคืบหน้า</h2>
        {application ? (
          // เปลี่ยนสี Card ตามสถานะ (Reject = แดง, Change = ส้ม, ปกติ = เขียว)
          <div className={`bg-white rounded-[2.5rem] border-2 shadow-xl p-8 md:p-10 flex flex-col gap-10 
            ${application.isRejected ? 'border-red-100 shadow-red-50' : 
              application.isAwardChanged ? 'border-orange-100 shadow-orange-50' : 
              'border-gray-50 shadow-green-50'}`}
          >
            
            <div className="flex flex-col md:flex-row gap-8 items-center">
                <div className={`w-24 h-24 rounded-3xl flex items-center justify-center shrink-0 shadow-inner 
                    ${application.isRejected ? 'bg-red-50 text-red-500' : 
                      application.isAwardChanged ? 'bg-orange-50 text-orange-500' : 
                      'bg-green-50 text-green-600'}`}
                >
                    <Award size={44} />
                </div>
                <div className="flex-1 text-center md:text-left">
                    <h3 className="text-2xl font-black text-gray-900">{application.category}</h3>
                    <div className="flex items-center justify-center md:justify-start gap-2 mt-2">
                        {application.isRejected ? <XCircle className="text-red-500" size={20} /> : 
                         application.isAwardChanged ? <Info className="text-orange-500" size={20} /> : 
                         <CheckCircle2 className="text-green-500 animate-pulse" size={20} />}
                        
                        <span className={`font-bold uppercase tracking-wide text-sm 
                            ${application.isRejected ? "text-red-600" : 
                              application.isAwardChanged ? "text-orange-600" : 
                              "text-green-700"}`}
                        >
                            {/* 🔥 แก้ไขบรรทัดนี้ ให้เรียกใช้งานฟังก์ชัน translateStatus */}
                            {translateStatus(application.status)}
                        </span>
                    </div>
                    <p className="text-gray-400 text-xs mt-2 font-medium flex items-center justify-center md:justify-start gap-1">
                        ส่งเมื่อ: {application.submittedDate}
                    </p>
                </div>
                <div className="text-right hidden md:block">
                    <p className="text-[10px] text-gray-400 font-black uppercase tracking-widest">ID ใบคำร้อง</p>
                    <p className="font-bold text-gray-800 text-lg">#{application.id.split('-')[0]}</p>
                </div>
            </div>

            <div className="relative pt-2 pb-6 px-4">
                <div className="absolute left-0 top-6 w-full h-1 bg-gray-100 -z-10 rounded-full" />
                <div 
                    className={`absolute left-0 top-6 h-1 -z-10 transition-all duration-1000 rounded-full 
                        ${application.isRejected ? 'bg-red-400' : 
                          application.isAwardChanged ? 'bg-orange-400' : 
                          'bg-green-500 shadow-lg shadow-green-100'}`}
                    style={{ width: `${(application.currentStep / (steps.length - 1)) * 100}%` }}
                />
                <div className="flex items-center justify-between w-full">
                    {steps.map((s, i) => (
                        <div key={i} className="flex flex-col items-center">
                            <div className={`w-10 h-10 rounded-2xl border-4 bg-white flex items-center justify-center transition-all duration-500
                                ${i < application.currentStep ? (application.isAwardChanged ? 'border-orange-500 bg-orange-500 text-white' : 'border-green-500 bg-green-500 text-white') : 
                                  i === application.currentStep ? (application.isRejected ? 'border-red-500 shadow-lg shadow-red-100' : application.isAwardChanged ? 'border-orange-500 shadow-lg shadow-orange-100' : 'border-green-500 shadow-lg shadow-green-100') : 
                                  'border-gray-100'}`}>
                                {i < application.currentStep ? <CheckCircle2 size={16} /> : 
                                 i === application.currentStep ? (
                                     application.isRejected ? <XCircle size={18} className="text-red-500" /> : 
                                     application.isAwardChanged ? <div className="w-2.5 h-2.5 rounded-full bg-orange-500 animate-pulse" /> :
                                     <div className="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
                                 ) : 
                                 <div className="w-2 h-2 rounded-full bg-gray-200" />}
                            </div>
                            <span className={`text-[12px] mt-3 font-black uppercase tracking-tighter ${i === application.currentStep ? "text-gray-900" : "text-gray-400"}`}>
                                {s.label}
                            </span>
                        </div>
                    ))}
                </div>
            </div>
            
        

          </div>
        ) : (
          <div className="p-16 bg-white border-2 border-dashed border-gray-200 rounded-[2.5rem] text-center group hover:border-green-300 transition-all">
            <div className="w-20 h-20 bg-gray-50 rounded-3xl flex items-center justify-center mx-auto mb-6 group-hover:bg-green-50 transition-colors">
                <BookOpen className="text-gray-300 group-hover:text-green-500 transition-colors" size={40} />
            </div>
            <h3 className="text-xl font-black text-gray-800">ยังไม่มีการส่งคำร้อง</h3>
            <p className="text-gray-500 mt-2">เลือกหมวดหมู่ด้านล่างเพื่อเริ่มต้น</p>
          </div>
        )}
      </section>

      {/* Available Categories Section */}
      <section>
        <div className="flex items-center justify-between mb-8">
            <h2 className="text-xl font-bold text-gray-800 uppercase tracking-widest">หมวดหมู่การส่งคำร้อง</h2>
        </div>
        

        {!hasOpenTerm ? (
          <div className="flex flex-col items-center justify-center p-16 bg-white border border-gray-200 rounded-[2.5rem] text-center shadow-sm">
            <div className="w-20 h-20 bg-amber-50 rounded-full flex items-center justify-center text-amber-500 mb-6">
              <CalendarX2 size={40} />
            </div>
            <h3 className="text-2xl font-black text-gray-800 mb-2">System Closed</h3>
            <p className="text-gray-500 font-medium">
              ขณะนี้ไม่อยู่ในช่วงเวลาเปิดรับสมัคร กรุณาติดตามประกาศจากทางมหาวิทยาลัย
            </p>
          </div>
        ) : (
    
          <>
            {application && (
              <div className="mb-10 p-6 bg-blue-50 border border-blue-100 rounded-[2rem] flex gap-5 items-center text-blue-800 shadow-sm">
                <Info size={24} className="shrink-0" />
                <p className="text-sm font-bold leading-relaxed">
                    คุณได้ส่งคำร้องในเทอมนี้แล้ว สามารถติดตามคำร้องเพื่อดูความคืบหน้าและผลการพิจารณา 
                </p>
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
              {categories.map((cat) => (
                <div 
                  key={cat.id} 
                  className={`group relative bg-white rounded-[2.5rem] border border-gray-100 p-10 flex flex-col h-full transition-all duration-500
                    ${!application ? 'hover:shadow-2xl hover:-translate-y-3 hover:border-green-200' : 'opacity-40 grayscale pointer-events-none'}`}
                >
                  <div className="flex justify-between items-start mb-8">
                    <div className="bg-green-50 p-4 rounded-2xl text-green-600 shadow-inner group-hover:bg-green-600 group-hover:text-white transition-all duration-500">
                      <Award size={32} />
                    </div>
                    {!application && <ChevronRight className="text-gray-200 group-hover:text-green-500 transition-colors" size={28} />}
                  </div>

                  <h3 className="font-black text-gray-900 text-xl mb-4 leading-tight uppercase tracking-tight">{cat.name}</h3>
                  <p className="text-sm text-gray-500 mb-12 flex-1 leading-relaxed font-medium">
                    {cat.description || "Submit your achievements for evaluation by the committee."}
                  </p>
                  
                  {application ? (
                    <div className="w-full py-4 bg-gray-50 text-gray-400 rounded-2xl text-center text-[10px] font-black tracking-widest border border-gray-100 uppercase">
                      Already Applied
                    </div>
                  ) : (
                    <Link 
                      href={`/application/form?award_id=${cat.id}`} 
                      className="w-full py-4.5 bg-gray-900 text-white rounded-[1.5rem] text-xs font-black tracking-[0.2em] hover:bg-green-600 text-center transition-all shadow-xl shadow-gray-200 uppercase active:scale-95 flex items-center justify-center h-12"
                    >
                      Apply Now
                    </Link>
                  )}
                </div>
              ))}
            </div>
          </>
        )}
      </section>
    </div>
  );
}