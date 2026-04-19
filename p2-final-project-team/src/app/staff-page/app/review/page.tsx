"use client";

import { useRouter } from 'next/navigation';
import { useState, useEffect } from "react";
import Sidebar from "../../components/Sidebar"; 
import api from "../../../../lib/axios"; 
import { Eye, Edit3, CheckCircle, XCircle, Loader2, Search, Filter, Info, ArrowRight, X, FileText, Download } from "lucide-react"; // 🔥 เพิ่ม FileText, Download

interface Application {
  id: string; 
  student: {
    id: string;
    first_name: string;
    last_name: string;
    nisit_id: string;
    faculty_name: string;
    department_name: string;
  };
  award_id: number; 
  award_name: string;
  current_status: string; 
  created_at: string;
  additional_data?: any; 
}

interface AwardCategory {
  id: number;
  name: string;
  form_structure: any; 
}

export default function ReviewPage() {
  const FilePreviewModal = ({ url, name, onClose }: { url: string; name: string; onClose: () => void }) => {
  const isImage = /\.(jpg|jpeg|png|webp)$/i.test(name);
  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-in fade-in duration-300">
      <div className="relative w-full max-w-5xl h-[85vh] bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col">
        <div className="p-4 border-b flex justify-between items-center bg-white">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-50 rounded-lg text-blue-600"><FileText size={20} /></div>
            <h3 className="font-bold text-gray-800 truncate max-w-xs md:max-w-md">{name}</h3>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded-full text-gray-400"><X size={24} /></button>
        </div>
        <div className="flex-1 bg-gray-50 flex items-center justify-center p-2">
          {isImage 
            ? <img src={url} alt={name} className="max-w-full max-h-full object-contain" /> 
            : <iframe src={`${url}#view=FitH`} className="w-full h-full border-none rounded-xl bg-white" title={name} />
          }
        </div>
      </div>
    </div>
  );
};
  const [previewFile, setPreviewFile] = useState<{ url: string; name: string } | null>(null);
  const router = useRouter();
  const [applications, setApplications] = useState<Application[]>([]);
  const [categories, setCategories] = useState<AwardCategory[]>([]);
  
  const [searchQuery, setSearchQuery] = useState("");
  const [filterStatus, setFilterStatus] = useState<string>("all");
  const [showFilterMenu, setShowFilterMenu] = useState(false);
  const [loading, setLoading] = useState(true);

  const [showChangeModal, setShowChangeModal] = useState(false);
  const [selectedAppId, setSelectedAppId] = useState<string>("");
  const [selectedAppDetail, setSelectedAppDetail] = useState<any>(null); 
  const [newCategoryId, setNewCategoryId] = useState<number | "">("");
  const [changeRemark, setChangeRemark] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [newFormStructure, setNewFormStructure] = useState<any[]>([]);
  const [newAdditionalData, setNewAdditionalData] = useState<Record<string, any>>({});
  const [latestTerm, setLatestTerm] = useState<any>(null);

  const translateStatus = (status: string) => {
    const s = status.toUpperCase();
    if (s === "SUBMITTED") return "ส่งใบสมัครแล้ว";
    if (s === "APPROVED_PENDING_PRESIDENT") return "อนุมัติแล้ว";
    if (s.includes("PENDING")) return "รอดำเนินการ";
    if (s === "APPROVED") return "อนุมัติแล้ว";
    if (s.includes("REJECTED")) return "ปฏิเสธ / ไม่ผ่านอนุมัติ";
    if (s === "AWARD_CHANGED") return "เปลี่ยนประเภทรางวัล";
    return status.replace(/_/g, ' '); 
  };

  useEffect(() => {
    loadApplications();
    loadCategories();
  }, []);

  async function loadApplications() {
    try {
      setLoading(true);

      let termId = null;
      try {
        const termRes = await api.get("/api/v1/academic-terms/latest");
        const termData = termRes.data?.data || termRes.data;
        setLatestTerm(termData);
        termId = termData?.id;
      } catch (e) {
        console.warn("Could not fetch latest term:", e);
      }

      const url = termId
        ? `/api/v1/requests?academic_term_id=${termId}`
        : `/api/v1/requests`;

      const response = await api.get(url);
      const appsData = response.data.data?.data || response.data.data || [];

      setApplications(appsData.sort((a: any, b: any) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      ));
    } catch (error) {
      console.error('Failed to load applications:', error);
      alert('ไม่สามารถโหลดข้อมูลใบสมัครได้');
    } finally {
      setLoading(false);
    }
  }

  async function loadCategories() {
    try {
      const res = await api.get('/api/v1/award-categories');
      const catsData = res.data.data || res.data || [];
      
      const parsedCats = catsData.map((c: any) => {
         if (typeof c.form_structure === 'string') {
             try { c.form_structure = JSON.parse(c.form_structure); } catch (e) {}
         }
         return c;
      });
      setCategories(parsedCats.filter((c: any) => c.is_active));
    } catch (error) {
      console.error('Failed to load categories:', error);
    }
  }

  const handleApprove = async (id: string) => {
    const app = applications.find(a => a.id === id);
    if (!window.confirm(`ยืนยันการอนุมัติใบสมัครของ ${app?.student?.first_name} ${app?.student?.last_name}?`)) return;

    try {
      await api.post(`/api/v1/requests/${id}/approve`, { remark: 'อนุมัติผ่านระบบ' });
      alert("อนุมัติสำเร็จ!");
      loadApplications();
    } catch (error: any) {
      alert(`อนุมัติไม่สำเร็จ: ${error.response?.data?.message || 'Error'}`);
    }
  };

  const handleReject = async (id: string) => {
    const app = applications.find(a => a.id === id);
    const reason = window.prompt(`ปฏิเสธใบสมัครของ ${app?.student?.first_name}\n\nกรุณาระบุเหตุผล:`);
    
    if (reason !== null && reason.trim()) {
      try {
        await api.post(`/api/v1/requests/${id}/reject`, { remark: reason });
        alert('ปฏิเสธสำเร็จ!');
        loadApplications();
      } catch (error: any) {
        alert(`ปฏิเสธไม่สำเร็จ: ${error.response?.data?.message || 'Error'}`);
      }
    }
  };

  const openChangeCategoryModal = async (id: string) => {
    setSelectedAppId(id);
    setNewCategoryId("");
    setChangeRemark("");
    setNewFormStructure([]);
    setNewAdditionalData({});
    setSelectedAppDetail(null);
    setShowChangeModal(true);

    try {
        const res = await api.get(`/api/v1/requests/${id}`);
        let fullData = res.data.data || res.data;
        if (typeof fullData.additional_data === 'string') {
            try { fullData.additional_data = JSON.parse(fullData.additional_data); } catch(e){}
        }
        if (fullData.award && typeof fullData.award.form_structure === 'string') {
            try { fullData.award.form_structure = JSON.parse(fullData.award.form_structure); } catch(e){}
        }
        setSelectedAppDetail(fullData);
    } catch (err) {
        console.error("Cannot load full request detail", err);
    }
  };

  const handleCategoryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
      const val = e.target.value;
      setNewCategoryId(val === "" ? "" : Number(val));
      
      if (!val) {
          setNewFormStructure([]);
          setNewAdditionalData({});
          return;
      }

      const foundCat = categories.find(c => c.id === Number(val));
      if (foundCat && Array.isArray(foundCat.form_structure)) {
          setNewFormStructure(foundCat.form_structure);
          const initialData: Record<string, any> = {};
          foundCat.form_structure.forEach((field: any) => {
              initialData[field.id] = "";
          });
          setNewAdditionalData(initialData);
      } else {
          setNewFormStructure([]);
          setNewAdditionalData({});
      }
  };

  const handleNewFormDataChange = (fieldId: string, val: string) => {
      setNewAdditionalData(prev => ({ ...prev, [fieldId]: val }));
  };

  const handleChangeCategorySubmit = async () => {
    if (!newCategoryId) return alert("กรุณาเลือกประเภทรางวัลใหม่");
    if (!changeRemark.trim()) return alert("กรุณาระบุเหตุผลในการเปลี่ยนแปลง");

    for (const field of newFormStructure) {
        if (field.required && !newAdditionalData[field.id]) {
            return alert(`กรุณากรอกข้อมูล: ${field.label}`);
        }
    }

    try {
      setIsSubmitting(true);
      await api.post(`/api/v1/requests/${selectedAppId}/change-category`, { 
          new_category_id: Number(newCategoryId), 
          remark: changeRemark,
          additional_data: newAdditionalData 
      });
      
      alert('เปลี่ยนประเภทรางวัลสำเร็จ!');
      setShowChangeModal(false);
      loadApplications();
    } catch (error: any) {
      alert(`เปลี่ยนประเภทรางวัลไม่สำเร็จ: ${error.response?.data?.message || 'Error'}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  const filteredApplications = applications.filter((app) => {
    if(!app.student) return false; 
    
    const fullName = `${app.student.first_name} ${app.student.last_name}`.toLowerCase();
    const nisitId = app.student.nisit_id || "";
    const matchSearch = fullName.includes(searchQuery.toLowerCase()) || nisitId.includes(searchQuery);
    
    let matchStatus = true;
    const currentStatus = (app.current_status || "").toUpperCase();
    
    if (filterStatus === "pending") {
      matchStatus = currentStatus.includes("PENDING");
    } else if (filterStatus === "approved") {
      matchStatus = currentStatus === "APPROVED";
    } else if (filterStatus === "rejected") {
      matchStatus = currentStatus.includes("REJECTED");
    }

    return matchSearch && matchStatus;
  });

  if (loading) {
    return (
      <div className="flex min-h-screen bg-slate-50 font-sans">
        <Sidebar activePage="review" />
        <main className="flex-1 flex flex-col items-center justify-center gap-4">
            <Loader2 className="animate-spin text-blue-600" size={40} />
            <p className="text-slate-500 font-medium">กำลังโหลดข้อมูลใบสมัคร...</p>
        </main>
      </div>
    );
  }

  const getOldFieldLabel = (fieldId: string) => {
      if (!selectedAppDetail?.award?.form_structure || !Array.isArray(selectedAppDetail.award.form_structure)) {
          return fieldId.replace(/_/g, ' '); 
      }
      const matchedField = selectedAppDetail.award.form_structure.find((f: any) => f.id === fieldId);
      return matchedField ? matchedField.label : fieldId.replace(/_/g, ' ');
  };

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans text-slate-800">
      <Sidebar activePage="review" />

      <main className="flex-1 p-6 lg:p-10 ml-[260px] overflow-y-auto">
      <header className="mb-8">
        <h1 className="text-3xl lg:text-4xl font-extrabold text-slate-900 tracking-tight">
          ตรวจสอบใบสมัคร
        </h1>
        {latestTerm && (
          <p className="mt-2 text-sm text-slate-500">
            ภาคการศึกษา{" "}
            <span className="font-bold text-blue-600 bg-blue-50 border border-blue-100 px-2.5 py-0.5 rounded-lg">
              {latestTerm.semester === "first" ? "ต้น" : latestTerm.semester === "second" ? "ปลาย" : latestTerm.semester}
              /{latestTerm.academic_year}
            </span>
          </p>
        )}
      </header>

        <div className="bg-white rounded-[24px] border border-slate-200 shadow-sm overflow-hidden">
          <div className="p-6 border-b border-slate-100 bg-slate-50/50 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <h2 className="text-xl font-bold text-slate-900">รายการใบสมัครทั้งหมด</h2>
            
            <div className="flex items-center gap-3 w-full md:w-auto">
              <div className="relative flex-1 md:w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                <input
                  type="text"
                  placeholder="ค้นหาชื่อ หรือ รหัสนิสิต..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-10 pr-4 py-2 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm transition-all"
                />
              </div>

              <div className="relative">
                <button 
                  className="flex items-center gap-2 bg-white border border-slate-200 hover:bg-slate-50 px-4 py-2 rounded-xl text-sm font-bold text-slate-700 transition-colors"
                  onClick={() => setShowFilterMenu(!showFilterMenu)}
                >
                  <Filter size={16} /> 
                  {filterStatus === 'all' ? 'ทั้งหมด' : filterStatus === 'pending' ? 'รอตรวจสอบ' : filterStatus === 'approved' ? 'อนุมัติแล้ว' : 'ปฏิเสธแล้ว'}
                </button>
                
                {showFilterMenu && (
                  <div className="absolute right-0 mt-2 w-48 bg-white border border-slate-100 shadow-xl rounded-2xl overflow-hidden z-10">
                    <div className="px-4 py-2 bg-slate-50 border-b border-slate-100 text-xs font-bold text-slate-400 uppercase tracking-widest">กรองตามสถานะ</div>
                    <div className="flex flex-col">
                      <button className={`px-4 py-2.5 text-left text-sm font-medium hover:bg-blue-50 transition-colors ${filterStatus === 'all' ? 'text-blue-600 bg-blue-50/50' : 'text-slate-700'}`} onClick={() => { setFilterStatus('all'); setShowFilterMenu(false); }}>ทั้งหมด</button>
                      <button className={`px-4 py-2.5 text-left text-sm font-medium hover:bg-blue-50 transition-colors ${filterStatus === 'pending' ? 'text-blue-600 bg-blue-50/50' : 'text-slate-700'}`} onClick={() => { setFilterStatus('pending'); setShowFilterMenu(false); }}>รอตรวจสอบ</button>
                      <button className={`px-4 py-2.5 text-left text-sm font-medium hover:bg-blue-50 transition-colors ${filterStatus === 'approved' ? 'text-blue-600 bg-blue-50/50' : 'text-slate-700'}`} onClick={() => { setFilterStatus('approved'); setShowFilterMenu(false); }}>อนุมัติแล้ว</button>
                      <button className={`px-4 py-2.5 text-left text-sm font-medium hover:bg-blue-50 transition-colors ${filterStatus === 'rejected' ? 'text-blue-600 bg-blue-50/50' : 'text-slate-700'}`} onClick={() => { setFilterStatus('rejected'); setShowFilterMenu(false); }}>ปฏิเสธแล้ว</button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse min-w-[800px]">
              <thead>
                <tr className="bg-white border-b border-slate-200 text-xs uppercase tracking-widest text-slate-400">
                  <th className="p-5 font-bold">ลำดับ</th>
                  <th className="p-5 font-bold">ข้อมูลนิสิต</th>
                  <th className="p-5 font-bold">คณะ</th>
                  <th className="p-5 font-bold">ประเภทรางวัล</th>
                  <th className="p-5 font-bold text-center">สถานะ</th>
                  <th className="p-5 font-bold text-center">จัดการ</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 text-sm">
                {filteredApplications.length === 0 ? (
                    <tr>
                        <td colSpan={6} className="text-center p-12 text-slate-400 bg-slate-50/50">
                          ไม่พบข้อมูลใบสมัครที่ตรงกับการค้นหา
                        </td>
                    </tr>
                ) : filteredApplications.map((app, index) => {
                  
                  // ใหม่
                  const currentStatus = app.current_status?.toUpperCase();
                  const isPendingPresident = currentStatus === "APPROVED_PENDING_PRESIDENT";
                  const isPending = currentStatus?.includes("PENDING") && !isPendingPresident;
                  const isApproved = currentStatus === "APPROVED" || currentStatus === "AWARD_APPROVED";
                  const isRejected = currentStatus?.includes("REJECTED");

                  return (
                    <tr key={app.id} className="hover:bg-slate-50 transition-colors group">
                        <td className="p-5 text-slate-500 font-medium">{index + 1}</td>
                        <td className="p-5">
                          <p className="font-bold text-slate-900">{app.student?.first_name} {app.student?.last_name}</p>
                          <p className="text-xs text-slate-500 font-mono mt-0.5">{app.student?.nisit_id || '-'}</p>
                        </td>
                        <td className="p-5 text-slate-600 font-medium">
                          {app.student?.faculty_name || '-'}
                        </td>
                    <td className="p-5">
                      <div className="flex items-center gap-2">
                        <span className="bg-blue-50 text-blue-700 px-3 py-1 rounded-full font-bold text-xs border border-blue-100">
                          {app.award_name}
                        </span>
                        
   
                        {!isRejected && !isPendingPresident && (
                          <button 
                              onClick={() => openChangeCategoryModal(app.id)}
                              title="เปลี่ยนประเภทรางวัล"
                              className="flex items-center gap-1.5 whitespace-nowrap text-xs font-bold text-slate-400 hover:text-blue-600 bg-white hover:bg-blue-50 p-1.5 rounded-lg border border-transparent hover:border-blue-200 transition-all opacity-0 group-hover:opacity-100"
                          >
                              <Edit3 size={14} /> เปลี่ยน
                          </button>
                        )}
                      </div>
                    </td>
          
                      <td className="p-5 text-center">
                        {isPending && (
                            <span className="inline-flex items-center gap-1.5 bg-amber-50 text-amber-700 border border-amber-200 px-3 py-1 rounded-full text-xs font-bold">
                              <span className="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span>
                              {translateStatus(app.current_status)}
                            </span>
                        )}
                        {isPendingPresident && (
                            <span className="inline-flex items-center gap-1.5 bg-green-50 text-green-700 border border-green-200 px-3 py-1 rounded-full text-xs font-bold">
                              <span className="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
                              {translateStatus(app.current_status)}
                            </span>
                        )}
                        {isApproved && (
                            <span className="inline-flex items-center gap-1 bg-green-50 text-green-700 border border-green-200 px-3 py-1 rounded-full text-xs font-bold">
                              <CheckCircle size={12}/> {translateStatus(app.current_status)}
                            </span>
                        )}
                        {isRejected && (
                            <span className="inline-flex items-center gap-1 bg-red-50 text-red-700 border border-red-200 px-3 py-1 rounded-full text-xs font-bold">
                              <XCircle size={12}/> {translateStatus(app.current_status)}
                            </span>
                        )}
                        {(!isPending && !isPendingPresident && !isApproved && !isRejected) && (
                            <span className="inline-flex bg-slate-100 text-slate-600 px-3 py-1 rounded-full text-xs font-bold">
                              {translateStatus(app.current_status)}
                            </span>
                        )}
                      </td>
                        <td className="p-5">
                          <div className="flex items-center justify-center gap-2">
                              <button 
                                  onClick={() => router.push(`/staff-page/app/review/${app.id}`)}
                                  className="p-2 bg-white border border-slate-200 hover:border-blue-300 hover:bg-blue-50 text-slate-500 hover:text-blue-600 rounded-xl transition-all shadow-sm"
                                  title="ดูรายละเอียดคำร้อง"
                              >
                                  <Eye size={16} />
                              </button>
                              
                              <button
                                  onClick={() => handleApprove(app.id)}
                                  disabled={!isPending || isPendingPresident}
                                  className="flex items-center gap-1.5 px-3 py-2 bg-white border border-slate-200 hover:border-green-300 hover:bg-green-50 text-slate-600 hover:text-green-700 rounded-xl text-xs font-bold transition-all disabled:opacity-40 disabled:cursor-not-allowed shadow-sm"
                              >
                                  <CheckCircle size={14} /> อนุมัติ
                              </button>
                              
                              <button
                                  onClick={() => handleReject(app.id)}
                                  disabled={!isPending || isPendingPresident}
                                  className="flex items-center gap-1.5 px-3 py-2 bg-white border border-slate-200 hover:border-red-300 hover:bg-red-50 text-slate-600 hover:text-red-700 rounded-xl text-xs font-bold transition-all disabled:opacity-40 disabled:cursor-not-allowed shadow-sm"
                              >
                                  <XCircle size={14} /> ปฏิเสธ
                              </button>
                          </div>
                        </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        </div>
      </main>

      {/* MODAL  */}
      {showChangeModal && (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center z-[60] p-4">
          <div className="bg-white w-full max-w-6xl rounded-[24px] overflow-hidden shadow-2xl flex flex-col max-h-[90vh]">
            
            <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-white shrink-0">
              <h3 className="text-xl font-extrabold text-slate-900 flex items-center gap-2">
                <Edit3 className="text-blue-600" /> เปลี่ยนประเภทรางวัล
              </h3>
              <button onClick={() => setShowChangeModal(false)} className="p-2 text-slate-400 hover:text-red-500 bg-slate-50 hover:bg-red-50 rounded-full transition-colors"><X size={20}/></button>
            </div>
            
            <div className="flex-1 overflow-y-auto p-6 bg-slate-50/50">
                {!selectedAppDetail ? (
                    <div className="flex justify-center items-center py-20">
                        <Loader2 className="animate-spin text-blue-500" size={40} />
                    </div>
                ) : (
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                        
                        <div className="space-y-4">
                            <h4 className="text-sm font-bold text-slate-400 uppercase tracking-widest flex items-center gap-2 border-b border-slate-200 pb-2">
                                <Info size={16}/> ข้อมูลคำร้องเดิม
                            </h4>
                            
                            <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm">
                                <p className="text-[10px] font-bold text-blue-500 uppercase tracking-widest mb-1">รางวัลที่สมัครไว้</p>
                                <p className="font-bold text-slate-800 text-lg">{selectedAppDetail?.award?.name || '-'}</p>
                            </div>

                            <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm space-y-4">
                                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">ข้อมูลฟอร์มเดิม</p>
                                
                                {selectedAppDetail?.additional_data && Object.keys(selectedAppDetail.additional_data).length > 0 ? (
                                    <div className="space-y-3">
                                        {Object.entries(selectedAppDetail.additional_data).map(([k, v]: any) => {
                                            if (k === 'items' || k === 'statement_content') return null; 
                                            return (
                                                <div key={k} className="bg-slate-50 p-3 rounded-xl border border-slate-100">
                                                    <p className="text-[10px] text-slate-400 font-bold mb-1">{getOldFieldLabel(k)}</p>
                                                    <p className="text-sm font-medium text-slate-800 break-words">{String(v)}</p>
                                                </div>
                                            )
                                        })}
                                    </div>
                                ) : (
                                    <p className="text-sm text-slate-400 italic text-center py-4">ไม่มีข้อมูลฟอร์มเพิ่มเติม</p>
                                )}
                            </div>

                            <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm space-y-4">
                                <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest">เอกสารแนบเดิม</p>
                                {selectedAppDetail?.documents && selectedAppDetail.documents.length > 0 ? (
                                    <div className="space-y-3">
                                        {selectedAppDetail.documents.map((doc: any, i: number) => (
                                            <div key={i} className="flex items-center justify-between p-3 bg-slate-50 border border-slate-100 rounded-xl">
                                                <div className="flex items-center gap-3 overflow-hidden">
                                                    <div className="w-8 h-8 bg-purple-100 text-purple-600 rounded-lg flex items-center justify-center shrink-0">
                                                        <FileText size={16} />
                                                    </div>
                                                    <div className="min-w-0">
                                                        <p className="text-xs font-bold text-slate-800 truncate">{doc.original_name || "Document"}</p>
                                                        <p className="text-[9px] text-slate-400 font-bold">{typeof doc.size === "number" ? `${(doc.size / 1024).toFixed(2)} KB` : "-"}</p>
                                                    </div>
                                                </div>
                                                <div className="flex gap-1 shrink-0">
                                                    <button
                                                        onClick={() => setPreviewFile({ url: doc.file_url, name: doc.original_name })}
                                                        className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg"
                                                        title="ดูเอกสาร"
                                                    >
                                                        <Eye size={14} />
                                                    </button>
                                                    <a href={doc.file_url} download className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg" title="ดาวน์โหลด">
                                                        <Download size={14} />
                                                    </a>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="text-sm text-slate-400 italic text-center py-4">ไม่มีเอกสารแนบ</p>
                                )}
                            </div>

                        </div>

                        <div className="hidden lg:flex absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-10 h-10 bg-blue-100 rounded-full items-center justify-center text-blue-600 shadow-sm border-2 border-white z-10 mt-6">
                            <ArrowRight size={20} />
                        </div>

                        <div className="space-y-4 relative">
                            <h4 className="text-sm font-bold text-blue-500 uppercase tracking-widest flex items-center gap-2 border-b border-blue-200 pb-2">
                                <Edit3 size={16}/> แก้ไขเป็นรางวัลใหม่
                            </h4>

                            <div className="bg-white p-5 rounded-2xl border border-blue-200 shadow-md">
                                <label className="block text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-2">เลือกประเภทรางวัลที่เหมาะสม <span className="text-red-500">*</span></label>
                                <select 
                                    value={newCategoryId} 
                                    onChange={handleCategoryChange}
                                    className="w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-slate-50 text-slate-900 font-medium cursor-pointer"
                                >
                                    <option value="">-- กรุณาเลือกรางวัลใหม่ --</option>
                                    
                                    {categories
                                      .filter(c => c.id !== selectedAppDetail?.award?.id)
                                      .map(c => (
                                        <option key={c.id} value={c.id}>{c.name}</option>
                                    ))}
                                </select>
                            </div>

                            {newFormStructure.length > 0 && (
                                <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm space-y-4 animate-in fade-in slide-in-from-bottom-4">
                                    <p className="text-[10px] font-bold text-slate-400 uppercase tracking-widest border-b border-slate-100 pb-2">กรอกข้อมูลเพิ่มเติมให้เข้ากับรางวัลใหม่</p>
                                    
                                    {newFormStructure.map((field) => (
                                        <div key={field.id}>
                                            <label className="block text-xs font-bold text-slate-700 mb-1.5">
                                                {field.label} {field.required && <span className="text-red-500">*</span>}
                                            </label>

                                            {field.type === "textarea" ? (
                                                <textarea 
                                                    value={newAdditionalData[field.id] || ""}
                                                    onChange={(e) => handleNewFormDataChange(field.id, e.target.value)}
                                                    required={field.required}
                                                    rows={3}
                                                    className="w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none text-sm resize-none"
                                                />
                                            ) : field.type === "select" ? (
                                                <select 
                                                    value={newAdditionalData[field.id] || ""}
                                                    onChange={(e) => handleNewFormDataChange(field.id, e.target.value)}
                                                    required={field.required}
                                                    className="w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none text-sm bg-white"
                                                >
                                                    <option value="">-- กรุณาเลือก --</option>
                                                    {field.options?.map((opt: string) => <option key={opt} value={opt}>{opt}</option>)}
                                                </select>
                                            ) : (
                                                <input 
                                                    type={field.type}
                                                    value={newAdditionalData[field.id] || ""}
                                                    onChange={(e) => handleNewFormDataChange(field.id, e.target.value)}
                                                    required={field.required}
                                                    className="w-full px-3 py-2 border border-slate-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none text-sm"
                                                />
                                            )}
                                        </div>
                                    ))}
                                </div>
                            )}

                            <div className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm">
                                <label className="block text-[11px] font-bold text-slate-500 uppercase tracking-widest mb-2">เหตุผลในการเปลี่ยนแปลง <span className="text-red-500">*</span></label>
                                <textarea 
                                    value={changeRemark}
                                    onChange={(e) => setChangeRemark(e.target.value)}
                                    placeholder="เช่น ผลงานมีความสอดคล้องกับประเภทนี้มากกว่า..."
                                    rows={2}
                                    className="w-full px-4 py-3 rounded-xl border border-slate-200 outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-slate-50 text-slate-900 font-medium resize-none"
                                />
                            </div>

                        </div>

                    </div>
                )}
            </div>

            <div className="p-6 border-t border-slate-100 bg-white shrink-0 flex justify-end gap-3">
                <button 
                    onClick={() => setShowChangeModal(false)}
                    className="px-6 py-3 rounded-xl font-bold bg-slate-100 text-slate-600 hover:bg-slate-200 transition-colors"
                >
                    ยกเลิก
                </button>
                <button 
                    onClick={handleChangeCategorySubmit}
                    disabled={isSubmitting || !newCategoryId || !selectedAppDetail}
                    className="px-8 py-3 rounded-xl font-bold bg-blue-600 text-white hover:bg-blue-700 flex items-center gap-2 disabled:opacity-50 transition-all shadow-md shadow-blue-200"
                >
                    {isSubmitting ? <Loader2 size={18} className="animate-spin" /> : <CheckCircle size={18} />}
                    บันทึกการเปลี่ยนรางวัล
                </button>
            </div>

          </div>
        </div>
      )}
      {previewFile && (
        <FilePreviewModal 
          url={previewFile.url} 
          name={previewFile.name} 
          onClose={() => setPreviewFile(null)} 
        />
      )}
    </div>
  );
}