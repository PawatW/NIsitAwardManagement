"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "../../../context/AuthContext";
import api from "../../../lib/axios";
import { 
  CheckCircle2, XCircle, Clock, FileText, User, Award, 
  Calendar, Phone, GraduationCap, TrendingUp, UserCheck, Loader2, AlertCircle, ChevronLeft, MapPin
} from "lucide-react";

interface ApplicationDetail {
  id: string;
  current_status: string;
  created_at: string;
  snapshot_study_year: number;
  snapshot_gpa: number;
  snapshot_advisor: string;
  snapshot_phone: string;
  snapshot_date_of_birth?: string; 
  snapshot_address?: string;
  additional_data: any;
  award?: {
    name: string;
    description: string;
    form_structure: any[];
  };
  academic_term?: {
    semester: string;
    academic_year: number;
  };
  status_history?: any[];
  documents?: any[]; 
}

export default function MyApplicationPage() {
  const { user } = useAuth();
  const [appData, setAppData] = useState<ApplicationDetail | null>(null);
  const [loading, setLoading] = useState(true);

  const steps = [
    { label: "ส่งแล้ว" }, 
    { label: "หัวหน้าภาควิชา" }, 
    { label: "รองคณบดี" },
    { label: "คณบดี" }, 
    { label: "คณะกรรมการ" },
    { label: "ประธาน" },
  ];

  const mapStatusToStep = (status: string): number => {
    const s = status.toUpperCase();

    if (s.includes("REJECTED_BY_HOD")) return 1;
    if (s.includes("REJECTED_BY_VICE_DEAN")) return 2;
    if (s.includes("REJECTED_BY_DEAN")) return 3;
    if (s.includes("REJECTED_BY_COMMITTEE")) return 4;

    // กรณีปกติ (Pending)
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
      
      default: return status.replace(/_/g, ' '); 
    }
  };

  useEffect(() => {
    const fetchMyApplication = async () => {
      try {
        setLoading(true);
        const listRes = await api.get("/api/v1/requests/my");
        const requestsList = listRes.data?.data || listRes.data || [];
        
        if (requestsList.length > 0) {
          const requestId = requestsList[0].id;
          const detailRes = await api.get(`/api/v1/requests/${requestId}`);
          const detailData = detailRes.data.data || detailRes.data; // เผื่อติดชั้น data
          setAppData(detailData); 
        }
      } catch (error) {
        console.error("Failed to fetch application details:", error);
      } finally {
        setLoading(false);
      }
    };

    if (user) fetchMyApplication();
  }, [user]);

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-[70vh]">
        <Loader2 className="w-10 h-10 animate-spin text-green-600 mb-4" />
        <p className="text-gray-500 font-bold uppercase tracking-widest text-xs">Loading Application Data...</p>
      </div>
    );
  }

  if (!appData) {
    return (
      <div className="max-w-4xl mx-auto mt-20 p-16 bg-white border-2 border-dashed border-gray-200 rounded-[3rem] text-center">
        <FileText className="text-gray-300 mx-auto mb-6" size={64} />
        <h3 className="text-2xl font-black text-gray-800">ยังไม่มีใบสมัครที่ส่ง</h3>
        <p className="text-gray-500 mt-2 mb-8 font-medium">คุณยังไม่ได้ส่งใบสมัครใด ๆ สำหรับเทอมการศึกษาปัจจุบัน</p>
        <Link href="/dashboard" className="px-8 py-3 bg-green-600 text-white rounded-xl font-bold hover:bg-green-700 transition-colors">
          ไปที่หน้าหลัก
        </Link>
      </div>
    );
  }

  const isRejected = appData.current_status.toUpperCase().includes("REJECT");
  const isAwardChanged = appData.current_status.toUpperCase() === "AWARD_CHANGED";
  const currentStep = mapStatusToStep(appData.current_status);

  return (
    <div className="max-w-5xl mx-auto pb-20 px-4 space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      {/* HEADER */}
      <div className="flex items-center gap-4 pt-6 pb-2">
        <Link href="/dashboard" className="w-10 h-10 bg-gray-100 text-gray-500 rounded-full flex items-center justify-center hover:bg-green-100 hover:text-green-700 transition-colors">
          <ChevronLeft size={20} />
        </Link>
        <div>
          <h1 className="text-3xl font-black text-gray-900 leading-tight">รายละเอียดใบสมัคร</h1>
          
        </div>
      </div>

      {/* STATUS CARD & TIMELINE  */}
      <section className={`bg-white rounded-[2.5rem] border-2 shadow-xl p-8 md:p-10 ${isRejected ? 'border-red-100 shadow-red-50' : 'border-gray-50 shadow-green-50'}`}>
        <div className="flex flex-col md:flex-row gap-8 items-center mb-10">
          <div className={`w-24 h-24 rounded-3xl flex items-center justify-center shrink-0 shadow-inner ${isRejected ? 'bg-red-50 text-red-500' : 'bg-green-50 text-green-600'}`}>
            <Award size={44} />
          </div>
          <div className="flex-1 text-center md:text-left">
            <h3 className="text-2xl font-black text-gray-900">{appData.award?.name || "Award Application"}</h3>
            <div className="flex items-center justify-center md:justify-start gap-2 mt-2">
              {isRejected ? <XCircle className="text-red-500" size={20} /> : <CheckCircle2 className="text-green-500 animate-pulse" size={20} />}

              <span className={`font-bold uppercase tracking-wide text-sm ${isRejected ? "text-red-600" : isAwardChanged ? "text-orange-600" : "text-green-700"}`}>
                {translateStatus(appData.current_status)}
              </span>
            </div>
            <p className="text-gray-400 text-xs mt-2 font-medium flex items-center justify-center md:justify-start gap-1">
              <Calendar size={14}/> ส่งเมื่อ: {new Date(appData.created_at).toLocaleDateString()}
            </p>
          </div>
        </div>

        {/* Visual Stepper */}
        <div className="relative pt-2 pb-6 px-4">
            <div className="absolute left-0 top-6 w-full h-1.5 bg-gray-100 -z-10 rounded-full" />
            <div 
                className={`absolute left-0 top-6 h-1.5 -z-10 transition-all duration-1000 rounded-full ${isRejected ? 'bg-red-400' : 'bg-green-500 shadow-lg shadow-green-100'}`}
                style={{ width: `${(currentStep / (steps.length - 1)) * 100}%` }}
            />
            <div className="flex items-center justify-between w-full">
                {steps.map((s, i) => (
                    <div key={i} className="flex flex-col items-center">
                        <div className={`w-12 h-12 rounded-[1.2rem] border-4 bg-white flex items-center justify-center transition-all duration-500
                            ${i < currentStep ? 'border-green-500 bg-green-500 text-white' : 
                              i === currentStep ? (isRejected ? 'border-red-500 shadow-xl shadow-red-100' : 'border-green-500 shadow-xl shadow-green-100') : 
                              'border-gray-100'}`}>
                            {i < currentStep ? <CheckCircle2 size={20} /> : 
                              i === currentStep ? (isRejected ? <XCircle size={24} className="text-red-500" /> : <div className="w-3 h-3 rounded-full bg-green-500 animate-pulse" />) : 
                              <div className="w-2.5 h-2.5 rounded-full bg-gray-200" />}
                        </div>
                        <span className={`text-[12px] mt-4 font-black uppercase tracking-tighter text-center max-w-[60px] leading-tight ${i === currentStep ? "text-gray-900" : "text-gray-400"}`}>
                            {s.label}
                        </span>
                    </div>
                ))}
            </div>
        </div>
      </section>

      {/* REJECTION REASON */}
      {(isRejected || isAwardChanged) && appData.status_history && appData.status_history.length > 0 && (
        <div className={`p-6 rounded-2xl border flex gap-4 ${isRejected ? 'bg-red-50 border-red-100' : 'bg-orange-50 border-orange-100'}`}>
          <AlertCircle className={`shrink-0 ${isRejected ? 'text-red-500' : 'text-orange-500'}`} size={24} />
          <div>
            {/* ปรับแก้ข้อความเป็นภาษาไทยด้วยเลยเพื่อความกลมกลืน */}
            <h4 className={`font-bold uppercase text-xs tracking-widest mb-1 ${isRejected ? 'text-red-800' : 'text-orange-800'}`}>
               {isAwardChanged ? "หมายเหตุจากคณะกรรมการ (การเปลี่ยนรางวัล)" : "เหตุผลที่ไม่ผ่านการอนุมัติ"}
            </h4>
            <p className={`text-sm font-medium ${isRejected ? 'text-red-700' : 'text-orange-700'}`}>
              {appData.status_history[appData.status_history.length - 1]?.remark || "ไม่มีการระบุหมายเหตุเพิ่มเติม"}
            </p>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* LEFT COLUMN: Snapshot & Docs */}
        <div className="lg:col-span-1 space-y-8">
          
          <div className="bg-white p-8 rounded-[2rem] border border-gray-100 shadow-sm">
            <h3 className="font-black text-gray-800 text-lg mb-6 flex items-center gap-2 border-b border-gray-100 pb-4">
               <User size={20} className="text-green-600"/> ข้อมูลส่วนตัว
            </h3>
            <div className="space-y-5">
              <div className="flex items-center gap-3 text-sm">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400"><GraduationCap size={16}/></div>
                <div><p className="text-[10px] font-bold text-gray-400 uppercase">ชั้นปีที่กำลังศึกษา</p><p className="font-bold text-gray-800">ปีที่ {appData.snapshot_study_year}</p></div>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400"><TrendingUp size={16}/></div>
                <div><p className="text-[10px] font-bold text-gray-400 uppercase">GPA สะสม</p><p className="font-bold text-gray-800">{appData.snapshot_gpa?.toFixed(2)}</p></div>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400"><UserCheck size={16}/></div>
                <div><p className="text-[10px] font-bold text-gray-400 uppercase">อาจารย์ที่ปรึกษา</p><p className="font-bold text-gray-800">{appData.snapshot_advisor}</p></div>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400"><Phone size={16}/></div>
                <div><p className="text-[10px] font-bold text-gray-400 uppercase">เบอร์โทรศัพท์</p><p className="font-bold text-gray-800">{appData.snapshot_phone || "-"}</p></div>
              </div>
              <div className="flex items-center gap-3 text-sm">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400"><Calendar size={16}/></div>
                <div>
                  <p className="text-[10px] font-bold text-gray-400 uppercase">วันเกิด</p>
                  <p className="font-bold text-gray-800">
                    {appData.snapshot_date_of_birth ? new Date(appData.snapshot_date_of_birth).toLocaleDateString('th-TH') : "-"}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3 text-sm pt-4 border-t border-gray-50">
                <div className="w-8 h-8 rounded-lg bg-gray-50 flex items-center justify-center text-gray-400 shrink-0"><MapPin size={16}/></div>
                <div>
                  <p className="text-[10px] font-bold text-gray-400 uppercase">ที่อยู่ปัจจุบัน</p>
                  <p className="font-bold text-gray-800 leading-relaxed">
                    {appData.snapshot_address || "-"}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white p-8 rounded-[2rem] border border-gray-100 shadow-sm">
            <h3 className="font-black text-gray-800 text-lg mb-6 flex items-center gap-2 border-b border-gray-100 pb-4">
               <FileText size={20} className="text-green-600"/> เอกสารประกอบ
            </h3>
            {appData.documents && appData.documents.length > 0 ? (
              <div className="space-y-3">
                {appData.documents.map((doc, idx) => (
                  <a key={idx} href={doc.file_url} target="_blank" rel="noopener noreferrer" className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl border border-gray-100 hover:border-green-200 transition-colors cursor-pointer group">
                    <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center text-red-500 shadow-sm"><FileText size={14}/></div>
                    <div className="overflow-hidden flex-1">

                      <p className="text-xs font-bold text-gray-700 truncate group-hover:text-green-700 transition-colors">{doc.original_name || "Document"}</p>
                      <p className="text-[10px] text-gray-400">{typeof doc.size === 'number' ? (doc.size / 1024).toFixed(2) : "-"} KB</p>
                    </div>
                  </a>
                ))}
              </div>
            ) : (
              <p className="text-xs font-medium text-gray-400 italic">ไม่มีเอกสารแนบ</p>
            )}
          </div>
        </div>

        {/* RIGHT COLUMN: Specific Answers & Statement */}
        <div className="lg:col-span-2 space-y-8">
          <div className="bg-white p-8 rounded-[2rem] border border-gray-100 shadow-sm h-full">
             <h3 className="font-black text-gray-800 text-lg mb-6 border-b border-gray-100 pb-4">รายละเอียดคำตอบเพิ่มเติม</h3>
             

             <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                {Object.entries(appData.additional_data || {}).map(([key, value]) => {
                  if (key === "statement_content" || key === "items") return null; 
                  
                  const fieldSchema = appData.award?.form_structure?.find((f: any) => f.id === key);
                  const label = fieldSchema ? fieldSchema.label : key.replace(/_/g, ' ');

                  return (
                    <div key={key} className="bg-gray-50 p-4 rounded-2xl border border-gray-100">
                      <p className="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-1">{label}</p>
                      <p className="text-gray-900 font-bold text-sm">{String(value)}</p>
                    </div>
                  );
                })}
             </div>

   
             {appData.additional_data?.items && Array.isArray(appData.additional_data.items) && (
               <div className="space-y-4 mb-8">
                 <p className="text-[10px] font-black text-gray-400 uppercase tracking-widest">ข้อมูลผลงานเพิ่มเติม (ถ้ามี)</p>
                 {appData.additional_data.items.map((item: any, idx: number) => (
                    <div key={idx} className="bg-gray-50 p-4 rounded-2xl border border-gray-100 grid grid-cols-1 md:grid-cols-2 gap-4">
                      {Object.entries(item).map(([k, v]) => (
                        <div key={k}>
                          <span className="text-[10px] text-gray-400 uppercase font-bold block">{k}</span>
                          <span className="text-sm font-bold text-gray-800 block">{String(v)}</span>
                        </div>
                      ))}
                    </div>
                 ))}
               </div>
             )}

             {/* Render Statement */}
             {appData.additional_data?.statement_content && (
               <div className="mt-8 pt-8 border-t border-gray-100">
                 <h4 className="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-3">เรียงความ / เหตุผลที่สมควรได้รับรางวัล</h4>
                 <div className="bg-blue-50/50 p-6 rounded-2xl border border-blue-100/50">
                    <p className="text-gray-700 text-sm leading-relaxed font-medium whitespace-pre-wrap">
                      {appData.additional_data.statement_content}
                    </p>
                 </div>
               </div>
             )}
          </div>
        </div>

      </div>
    </div>
  );
}