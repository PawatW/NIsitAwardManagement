"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { 
  FileText, Download, CheckCircle, XCircle, Loader2, Eye, X, Award, Mail, BookOpen,
  ChevronLeft, Info, GraduationCap, Phone, Calendar, Users, MapPin 
} from "lucide-react";
import api from "../../../../lib/axios"; 
import { useAuth } from "../../../../context/AuthContext";


interface UserBasicResponse {
  id: string;
  first_name: string;
  last_name: string;
  nisit_id: string;
  phone_number?: string;
  campus_name?: string;     
  faculty_name?: string;
  department_name?: string;
}

interface AwardCategoryResponse {
  id: number;
  name: string;
  description: string;
  form_structure?: any[]; 
}

interface RequestDocumentResponse {
  id: string;
  original_name: string;
  file_url: string;
  size: number;
}

interface RequestResponse {
  id: string;
  student_id: string;
  award_id: number;
  academic_term_id: number;
  additional_data: any; 
  snapshot_study_year?: number;    
  snapshot_gpa?: number;
  snapshot_advisor?: string;       
  snapshot_phone?: string;
  snapshot_date_of_birth?: string; 
  snapshot_address?: string;       
  current_status: string;
  created_at: string;
  student?: UserBasicResponse;
  award?: AwardCategoryResponse;
  documents?: RequestDocumentResponse[]| null;
}

const FilePreviewModal = ({ url, name, onClose }: { url: string; name: string; onClose: () => void }) => {
  const isImage = /\.(jpg|jpeg|png|webp)$/i.test(name);
  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-300">
      <div className="relative w-full max-w-5xl h-[85vh] bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col">
        <div className="p-4 border-b flex justify-between items-center bg-white">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-green-50 rounded-lg text-green-600"><FileText size={20} /></div>
            <h3 className="font-bold text-gray-800 truncate max-w-xs md:max-w-md">{name}</h3>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded-full text-gray-400"><X size={24} /></button>
        </div>
        <div className="flex-1 bg-gray-50 flex items-center justify-center p-2">
          {isImage ? <img src={url} alt={name} className="max-w-full max-h-full object-contain" /> : 
          <iframe src={`${url}#view=FitH`} className="w-full h-full border-none rounded-xl bg-white" title={name} />}
        </div>
      </div>
    </div>
  );
};

export default function ApplicationReviewPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth(); 
  const applicationId = params.id as string; 

  const [data, setData] = useState<RequestResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);
  const [previewFile, setPreviewFile] = useState<{ url: string; name: string } | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response = await api.get(`/api/v1/requests/${applicationId}`);
        setData(response.data);
      } catch (error) {
        console.error("Error fetching detail:", error);
        alert("Could not load application details.");
        router.push("/approval");
      } finally {
        setLoading(false);
      }
    };
    if (applicationId) fetchData();
  }, [applicationId, router]);

  const handleApprove = async () => {
    if (!data || !user) return;
    if (!confirm(`Confirm approval?`)) return;

    try {
      setProcessing(true);
      await api.post(`/api/v1/requests/${applicationId}/approve`, { remark: "Approved" }); 
      alert("อนุมัติเรียบร้อยแล้ว! ");
      router.push("/approval"); 
    } catch (error: any) {
      alert(`Approval Failed: ${error.response?.data?.message || "Unknown error"}`);
    } finally {
      setProcessing(false);
    }
  };

  const handleReject = async () => {
    const reason = prompt("Please provide a reason for rejection (Required):");
    if (!reason) return; 

    try {
      setProcessing(true);
      await api.post(`/api/v1/requests/${applicationId}/reject`, { remark: reason });
      alert("ปฏิเสธเรียบร้อยแล้ว! ");
      router.push("/approval");
    } catch (error: any) {
      alert(`Reject Failed: ${error.response?.data?.message || "Unknown error"}`);
    } finally {
      setProcessing(false);
    }
  };

  const translateStatus = (status: string) => {
    const s = status.toUpperCase();
    if (s === "SUBMITTED") return "ส่งใบสมัครแล้ว";
    if (s.includes("PENDING")) return "รอดำเนินการ";
    if (s === "APPROVED") return "อนุมัติแล้ว";
    if (s.includes("REJECTED")) return "ปฏิเสธ/ไม่ผ่านอนุมัติ";
    if (s === "AWARD_CHANGED") return "เปลี่ยนประเภทรางวัล";
    return status.replace(/_/g, ' '); 
  };

  const getFieldLabel = (fieldId: string) => {
    if (!data?.award?.form_structure || !Array.isArray(data.award.form_structure)) {
        return fieldId.replace(/_/g, ' '); 
    }
    const matchedField = data.award.form_structure.find((field: any) => field.id === fieldId);
    return matchedField ? matchedField.label : fieldId.replace(/_/g, ' ');
  };

  if (loading) return <div className="flex h-[60vh] items-center justify-center"><Loader2 className="animate-spin text-green-600" size={40} /></div>;
  if (!data) return null;

  const studentName = data.student ? `${data.student.first_name} ${data.student.last_name}` : "Unknown";
  const studentId = data.student?.nisit_id || "-";
  const studentContact = data.snapshot_phone || data.student?.phone_number || "-"; 
  const categoryName = data.award?.name || "-";
  
  const rawStatus = data.current_status ? data.current_status.toUpperCase() : "PENDING";
  const isPending = rawStatus === "SUBMITTED" || rawStatus.includes("PENDING") || rawStatus === "AWARD_CHANGED";
  
  const specificAnswers = data.additional_data?.items || data.additional_data || {};
  const statementContent = data.additional_data?.statement_content || "No statement provided.";

  return (
    <div className="max-w-6xl mx-auto space-y-6 pb-20">
      {/* Top Bar */}
      <div className="flex items-center justify-between">
        <button onClick={() => router.back()} className="flex items-center gap-2 text-gray-500 hover:text-green-600 transition-colors font-medium">
          <ChevronLeft size={20} /> ย้อนกลับ
        </button>
        <div className="flex items-center gap-3">
            <span className="px-4 py-1.5 bg-green-50 text-green-700 border border-green-100 rounded-full text-xs font-bold uppercase tracking-wider">
                {categoryName}
            </span>
            <div className={`px-4 py-1.5 border rounded-full text-xs font-bold tracking-wider ${
               rawStatus === "APPROVED" ? "bg-green-50 text-green-600 border-green-200" :
               rawStatus.includes("REJECT") ? "bg-red-50 text-red-600 border-red-200" :
               "bg-yellow-50 text-yellow-600 border-yellow-200"
            }`}>
                {translateStatus(rawStatus)}
            </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          
 
          <div className="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm">
            {/* Header: รูป + ชื่อ + ข้อมูลหลัก */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center gap-6 pb-6 border-b border-gray-100">
          
              <div className="flex-1">
                <h1 className="text-2xl font-black text-gray-900 leading-tight">{studentName}</h1>
                <div className="flex flex-wrap items-center gap-3 mt-2">
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 bg-green-50 text-green-700 rounded-lg text-sm font-bold border border-green-100">
                    <Award size={14}/> {studentId}
                  </span>
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 bg-indigo-50 text-indigo-700 rounded-lg text-sm font-bold border border-indigo-100">
                    เกรดเฉลี่ย (GPA): {typeof data.snapshot_gpa === "number" ? data.snapshot_gpa.toFixed(2) : "-"}
                  </span>
                </div>
              </div>
            </div>

            {/* Grid ข้อมูลส่วนตัวอื่นๆ */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-y-6 gap-x-8 pt-6">
              
              {/* คณะ / สาขา */}
              <div>
                <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                  <GraduationCap size={14} className="text-green-500" /> คณะ / ภาควิชา
                </p>
                <p className="text-sm font-bold text-gray-800 leading-tight">
                  {data.student?.faculty_name || '-'} <br/>
                  <span className="text-xs font-medium text-gray-500">{data.student?.department_name || '-'}</span>
                </p>
              </div>

              {/* วิทยาเขต & ชั้นปี */}
              <div className="flex gap-8">
                <div>
                  <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                    <MapPin size={14} className="text-green-500" /> วิทยาเขต
                  </p>
                  <p className="text-sm font-bold text-gray-800">{data.student?.campus_name || '-'}</p>
                </div>
                <div>
                  <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                    <Calendar size={14} className="text-green-500" /> ชั้นปีที่
                  </p>
                  <p className="text-sm font-bold text-gray-800">{data.snapshot_study_year || '-'}</p>
                </div>
              </div>

              {/* อาจารย์ที่ปรึกษา */}
              <div>
                <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                  <Users size={14} className="text-green-500" /> อาจารย์ที่ปรึกษา
                </p>
                <p className="text-sm font-bold text-gray-800">{data.snapshot_advisor || '-'}</p>
              </div>

              {/* เบอร์โทรศัพท์ */}
              <div>
                <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                  <Phone size={14} className="text-green-500" /> เบอร์โทรศัพท์
                </p>
                <p className="text-sm font-bold text-gray-800">{studentContact}</p>
              </div>

              {/* วันเกิด */}
              <div>
                <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                  <Calendar size={14} className="text-green-500" /> วันเกิด
                </p>
                <p className="text-sm font-bold text-gray-800">
                  {data.snapshot_date_of_birth ? new Date(data.snapshot_date_of_birth).toLocaleDateString('en-GB') : '-'}
                </p>
              </div>

              {/* วันที่ส่งคำร้อง */}
              <div>
                <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-1.5">
                  <Calendar size={14} className="text-green-500" /> วันที่ส่งคำร้อง
                </p>
                <p className="text-sm font-bold text-gray-800">
                  {new Date(data.created_at).toLocaleString('en-GB', { 
                      day: '2-digit', month: 'short', year: 'numeric', 
                      hour: '2-digit', minute: '2-digit' 
                  })} น.
                </p>
              </div>

              {/* ที่อยู่ป */}
              {data.snapshot_address && (
                <div className="col-span-1 md:col-span-2">
                  <p className="text-[11px] font-extrabold text-gray-400 uppercase tracking-widest flex items-center gap-1.5 mb-2">
                    <MapPin size={14} className="text-green-500" /> ที่อยู่ปัจจุบัน
                  </p>
                  <div className="bg-gray-50 p-4 rounded-xl border border-gray-100">
                    <p className="text-sm font-medium text-gray-700 leading-relaxed">
                      {data.snapshot_address}
                    </p>
                  </div>
                </div>
              )}

            </div>
          </div>

          {/* Details Section */}
          <div className="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm space-y-6">
            <h3 className="text-lg font-bold text-gray-800 flex items-center gap-2">
              <BookOpen size={22} className="text-green-600" /> ข้อมูลรายละเอียดเฉพาะหมวดหมู่
            </h3>
            
            <div className="space-y-4">
              {Array.isArray(specificAnswers) ? (
                specificAnswers.map((item: any, idx: number) => (
                  <div key={idx} className="p-6 bg-gray-50 rounded-2xl border border-gray-100">
                    <p className="text-[10px] font-black text-green-600 uppercase mb-3 tracking-widest">Entry #{idx + 1}</p>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      {Object.entries(item).map(([k, v]: any) => {
                        if(k === "statement_content") return null;
                        return (
                          <div key={k}>
                            <p className="text-[10px] text-gray-400 font-black uppercase tracking-widest">{getFieldLabel(k)}</p>
                            <p className="text-gray-800 font-bold text-sm">{String(v)}</p>
                          </div>
                        )
                      })}
                    </div>
                  </div>
                ))
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {Object.entries(specificAnswers || {}).map(([k, v]: any) => {
                     if(k === "statement_content" || k === "items") return null;
                     return (
                      <div key={k} className="p-4 bg-gray-50 rounded-2xl border border-transparent">
                        <p className="text-[10px] text-gray-400 font-black uppercase tracking-widest">{getFieldLabel(k)}</p>
                        <p className="text-gray-800 font-bold mt-1">{String(v)}</p>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          </div>

          {/* Statement Section */}
          <div className="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm space-y-4">
            <h3 className="text-lg font-bold text-gray-800 flex items-center gap-2">
              <BookOpen size={22} className="text-green-600" /> คำชี้แจงจากนิสิต / Statement
            </h3>
            <div className="p-6 bg-green-50/30 rounded-2xl border border-green-50">
               <p className="text-gray-700 leading-relaxed whitespace-pre-line text-sm font-medium">
                 {statementContent}
               </p>
            </div>
          </div>
        </div>
        
        {/* Right Column - Actions & Attachments */}
        <div className="space-y-6">
          <div className="bg-white p-6 rounded-3xl border border-gray-100 shadow-lg shadow-green-50/50 sticky top-24">
            <h3 className="text-lg font-bold text-gray-800 mb-4">การดำเนินการ</h3>
            
            {isPending ? (
              <>

                <div className="space-y-3">
                    <button 
                        onClick={handleApprove}
                        disabled={processing}
                        className="w-full py-3 bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white rounded-xl font-bold shadow-md shadow-green-200 transition-all flex items-center justify-center gap-2"
                    >
                        {processing ? <Loader2 className="animate-spin" size={18}/> : <CheckCircle size={18}/>}
                        อนุมัติ / Approve
                    </button>

                    <button 
                        onClick={handleReject}
                        disabled={processing}
                        className="w-full py-3 bg-white border-2 border-red-100 text-red-600 hover:bg-red-50 disabled:bg-gray-50 rounded-xl font-bold transition-all flex items-center justify-center gap-2"
                    >
                        <XCircle size={18}/>
                        ปฏิเสธ / Reject
                    </button>
                </div>
              </>
            ) : (
              <div className="p-4 bg-gray-50 rounded-xl border border-gray-200 text-center">
                <p className="text-sm font-bold text-gray-500 uppercase tracking-widest">ดำเนินการเรียบร้อยแล้ว</p>
                <p className="text-xs text-gray-400 mt-1">ใบคำร้องนี้ได้รับการดำเนินการเรียบร้อยแล้ว</p>
              </div>
            )}
          </div>

          <div className="bg-white p-8 rounded-3xl border border-gray-100 shadow-sm">
            <h3 className="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2 font-mono">
              <FileText size={22} className="text-green-600" /> เอกสารที่เกี่ยวข้อง
            </h3>
            <div className="space-y-4">
              {data.documents && data.documents.length > 0 ? (
                data.documents.map((doc, i) => (
                  <div key={i} className="group p-4 bg-white border border-gray-100 rounded-2xl hover:border-green-500 hover:shadow-md transition-all flex items-center justify-between">
                    <div className="flex items-center gap-3 overflow-hidden">
                      <div className="w-10 h-10 bg-green-50 text-green-600 rounded-xl flex items-center justify-center shrink-0"><FileText size={20} /></div>
                      <div className="min-w-0">
                        <p className="text-sm font-bold text-gray-800 truncate">{doc.original_name || "Document"}</p>
                        <p className="text-[10px] text-gray-400 font-bold">{typeof doc.size === "number"? `${(doc.size / 1024).toFixed(2)} KB`: "-"}</p>
                      </div>
                    </div>
                    <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button onClick={() => setPreviewFile({ url: doc.file_url, name: doc.original_name })} className="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg"><Eye size={18} /></button>
                      <a href={doc.file_url} download className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg"><Download size={18} /></a>
                    </div>
                  </div>
                ))
              ) : <p className="text-gray-400 text-sm text-center py-4 italic">No documents attached.</p>}
            </div>
          </div>
        </div>
      </div>

      {previewFile && <FilePreviewModal url={previewFile.url} name={previewFile.name} onClose={() => setPreviewFile(null)} />}
    </div>
  );
}