"use client";

import { useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { Search, Check, X, Loader2, RefreshCcw, Award } from "lucide-react";
import api from "../../../lib/axios"; 
import { useAuth } from "../../../context/AuthContext";

interface ApplicationUI {
  id: string;
  name: string;
  studentId: string;
  department: string;
  category: string;
  submissionDate: string;
  status: string; 
  rawStatus: string; 
}


const roleTranslationMap: Record<string, string> = {
  "ADMIN": "ผู้ดูแลระบบ",
  "STUDENT": "นิสิต",
  "COMMITTEE_CHAIR": "คณะกรรมการ",
  "DEAN": "คณบดี",
  "VICE_DEAN": "รองคณบดี",
  "HEAD_OF_DEPARTMENT": "หัวหน้าภาควิชา"
};

export default function ApprovalPage() {

  
  const getThaiRole = (englishRole: string) => {
    return roleTranslationMap[englishRole.toUpperCase()] || englishRole.replace(/_/g, ' ');
  };

  const router = useRouter();
  const { user } = useAuth(); 
  
  const [applications, setApplications] = useState<ApplicationUI[]>([]);
  const [latestTerm, setLatestTerm] = useState<any>(null); 
  const [loading, setLoading] = useState(true);
  const [processingId, setProcessingId] = useState<string | null>(null);
  
  const [searchTerm, setSearchTerm] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("All Categories");
  const [statusFilter, setStatusFilter] = useState("All Status");
  const [refreshKey, setRefreshKey] = useState(0);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        
        let currentTermId = null;
        try {
          const termRes = await api.get("/api/v1/academic-terms/latest");
          const termData = termRes.data?.data || termRes.data;
          setLatestTerm(termData);
          currentTermId = termData?.id;
        } catch (termErr) {
          console.warn("Could not fetch latest term:", termErr);
        }

       
        const requestUrl = currentTermId 
          ? `/api/v1/requests?academic_term_id=${currentTermId}` 
          : "/api/v1/requests";

        const response = await api.get(requestUrl); 
        const rawData = response?.data?.data?.data || response?.data?.data || response?.data;
        
        const requestsData = Array.isArray(rawData) ? rawData : (rawData?.items || []);

        const mappedData: ApplicationUI[] = requestsData.map((item: any) => ({
          id: item.id,
          name: item.student ? `${item.student.first_name} ${item.student.last_name}` : "Unknown Student",
          studentId: item.student?.nisit_id || "-",
          department: item.student?.department_name || "-",
          category: item.award_name || "-",
          submissionDate: item.created_at
          ? new Date(item.created_at).toLocaleDateString("en-US", {
              month: "short",
              day: "numeric",
              year: "numeric"
            })
          : "-",
          status: mapStatusForDisplay(item.current_status),
          rawStatus: item.current_status
        }));

        setApplications(mappedData);
      } catch (error) {
        console.error("Error fetching data:", error);
        setApplications([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [refreshKey]);

  const handleQuickApprove = async (id: string) => {
    if (!user) return alert("Session expired.");
    if (!confirm(`ยืนยันการอนุมัติ?`)) return;

    try {
      setProcessingId(id);
      await api.post(`/api/v1/requests/${id}/approve`, { remark: "Approved via quick action" }); 
      setRefreshKey(k => k + 1);
    } catch (error: any) {
      alert(`Approval Failed: ${error.response?.data?.message || "Error"}`);
    } finally {
      setProcessingId(null);
    }
  };

  const handleQuickReject = async (id: string) => {
    const reason = prompt("กรุณาระบุเหตุผลในการปฏิเสธ:");
    if (!reason) return; 

    try {
      setProcessingId(id);
      await api.post(`/api/v1/requests/${id}/reject`, { remark: reason });
      setRefreshKey(k => k + 1);
    } catch (error: any) {
      alert(`Reject Failed: ${error.response?.data?.message || "Error"}`);
    } finally {
      setProcessingId(null);
    }
  };

  const mapStatusForDisplay = (backendStatus: string) => {
    const s = backendStatus ? backendStatus.toUpperCase() : "";
    if (!s || s === "SUBMITTED" || s.includes("PENDING")) return "Pending";
    if (s === "APPROVED") return "Approved";
    if (s.includes("REJECTED")) return "Rejected";
    if (s === "AWARD_CHANGED") return "Award_Changed"; 
    return "Pending";
  };

  const getAvatarColor = (name: string) => {
    const colors = ["bg-gray-900", "bg-indigo-900", "bg-blue-900", "bg-purple-900", "bg-green-900", "bg-red-900"];
    const index = name.charCodeAt(0) % colors.length;
    return colors[index];
  };

  const uniqueCategories = useMemo(() => {
    const categories = applications
      .map(app => app.category)
      .filter(cat => cat && cat !== "-");
    return Array.from(new Set(categories));
  }, [applications]);

  const filteredData = useMemo(() => {
    return applications.filter(item => {
      const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase()) || item.studentId.includes(searchTerm);
      const matchesCategory = categoryFilter === "All Categories" || item.category === categoryFilter;
      const matchesStatus = statusFilter === "All Status" || item.status === statusFilter;
      return matchesSearch && matchesCategory && matchesStatus;
    });
  }, [applications, searchTerm, categoryFilter, statusFilter]);

   const StatusBadge = ({ status }: { status: string }) => {
    const statusConfig: Record<string, { label: string; style: string }> = {
      Award_Changed: { 
        label: "เปลี่ยนประเภทรางวัล", 
        style: "bg-orange-50 text-orange-600 border border-orange-200" 
      },
      Pending: { 
        label: "รอดำเนินการ", 
        style: "bg-yellow-50 text-yellow-600 border border-yellow-200" 
      },
      Approved: { 
        label: "อนุมัติแล้ว", 
        style: "bg-green-50 text-green-600 border border-green-200" 
      },
      Rejected: { 
        label: "ไม่ผ่านการอนุมัติ", 
        style: "bg-red-50 text-red-600 border border-red-200" 
      },
    };

    const displayLabel = statusConfig[status]?.label || status.replace(/_/g, ' ');
    const displayStyle = statusConfig[status]?.style || "bg-gray-50 text-gray-600 border border-gray-200";

    return (
      <span className={`px-3 py-1 rounded-full text-xs font-bold ${displayStyle}`}>
        {displayLabel}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50/50 p-8 space-y-8 font-sans">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">รายการคำร้อง</h1>
          {/* 🔥 อัปเดตส่วนแสดงข้อมูลเทอมล่าสุด */}
          <div className="flex items-center gap-3 mt-2 text-sm text-gray-500">
            {latestTerm && (
              <span className="px-3 py-1 bg-blue-50 text-blue-700 rounded-lg font-semibold border border-blue-100">
                ภาคการศึกษา {latestTerm.semester === 'first' ? 'ต้น' : latestTerm.semester === 'second' ? 'ปลาย' : latestTerm.semester}/{latestTerm.academic_year}
              </span>
            )}
            <p>ตำแหน่ง: <span className="font-bold text-green-600">{getThaiRole(user?.role ?? "")}</span></p>
          </div>
        </div>
        <button onClick={() => setRefreshKey(k => k + 1)} className="flex items-center gap-2 text-sm text-gray-500 hover:text-green-600 transition-colors">
          <RefreshCcw size={16} className={loading ? "animate-spin" : ""} /> รีเฟรชข้อมูล
        </button>
      </div>

      <div className="bg-white rounded-3xl border border-gray-100 shadow-sm p-6">
        <div className="flex flex-col md:flex-row gap-4 mb-8">
          <div className="flex-1">
            <label className="text-xs font-semibold text-gray-500 mb-1.5 block ml-1">ค้นหา</label>
            <div className="relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
              <input 
                type="text" placeholder="ค้นหาชื่อนิสิตหรือรหัสนิสิต"
                className="w-full pl-12 pr-4 py-3 bg-gray-100 rounded-xl border-none focus:ring-2 focus:ring-gray-200 outline-none text-sm"
                value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>
          
          <div className="w-full md:w-48">
             <label className="text-xs font-semibold text-gray-500 mb-1.5 block ml-1">หมวดหมู่</label>
             <select 
                className="w-full px-4 py-3 bg-gray-100 rounded-xl border-none text-sm cursor-pointer outline-none" 
                value={categoryFilter} 
                onChange={(e) => setCategoryFilter(e.target.value)}
             >
               <option value="All Categories">ทั้งหมด</option>
               {uniqueCategories.map((cat, index) => (
                 <option key={index} value={cat}>{cat}</option>
               ))}
             </select>
          </div>

          <div className="w-full md:w-48">
             <label className="text-xs font-semibold text-gray-500 mb-1.5 block ml-1">สถานะ</label>
             <select className="w-full px-4 py-3 bg-gray-100 rounded-xl border-none text-sm cursor-pointer outline-none" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
               <option value="All Status">ทั้งหมด</option>
               <option value="Pending">รอดำเนินการ</option>
               <option value="Rejected">ปฏิเสธ</option>
             </select>
          </div>
        </div>

        <div className="overflow-x-auto min-h-[300px]">
          <table className="w-full">
            <thead>
              <tr className="text-left border-b border-gray-100">
                <th className="pb-4 pl-4 text-xs font-bold text-gray-400 uppercase tracking-wider">ผู้สมัคร</th>
                <th className="pb-4 text-xs font-bold text-gray-400 uppercase tracking-wider">สาขาวิชา</th>
                <th className="pb-4 text-xs font-bold text-gray-400 uppercase tracking-wider">หมวดหมู่</th>
                <th className="pb-4 text-xs font-bold text-gray-400 uppercase tracking-wider">วันที่ส่ง</th>
                <th className="pb-4 text-xs font-bold text-gray-400 uppercase tracking-wider text-center">สถานะ</th>
                <th className="pb-4 pr-4 text-xs font-bold text-gray-400 uppercase tracking-wider text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50">
              {loading ? (
                 <tr><td colSpan={6} className="py-20 text-center"><Loader2 className="w-8 h-8 animate-spin mx-auto text-green-600" /></td></tr>
              ) : filteredData.length === 0 ? (
                 <tr><td colSpan={6} className="py-20 text-center text-gray-500">ไม่พบข้อมูลคำร้องในภาคการศึกษานี้</td></tr>
              ) : filteredData.map((item) => (
                <tr key={item.id} className="group hover:bg-gray-50/50 transition-colors">
                  <td className="py-4 pl-4">
                    <div className="flex items-center gap-4 cursor-pointer" onClick={() => router.push(`/approval/${item.id}`)}>
                    
                      <div>
                        <p className="font-bold text-gray-900 text-sm group-hover:text-green-700">{item.name}</p>
                        <p className="text-xs text-gray-400 mt-0.5">ID: {item.studentId}</p>
                      </div>
                    </div>
                  </td>
                  <td className="py-4 text-sm text-gray-500 font-medium">{item.department}</td>
                  <td className="py-4 text-sm text-gray-500">{item.category}</td>
                  <td className="py-4 text-sm text-gray-500">{item.submissionDate}</td>
                  <td className="py-4 text-center"><StatusBadge status={item.status} /></td>
                  <td className="py-4 pr-4 text-right">
                    <div className="flex items-center justify-end gap-2">
                      {["Pending", "Award_Changed"].includes(item.status) ? (
                        <>
                          <button 
                            onClick={() => handleQuickApprove(item.id)} disabled={!!processingId}
                            className="w-8 h-8 rounded-full bg-green-100 text-green-600 flex items-center justify-center hover:bg-green-200 transition-colors disabled:opacity-50"
                            title="อนุมัติ"
                          >
                            {processingId === item.id ? <Loader2 size={14} className="animate-spin" /> : <Check size={16} />}
                          </button>
                          <button 
                            onClick={() => handleQuickReject(item.id)} disabled={!!processingId}
                            className="w-8 h-8 rounded-full bg-gray-100 text-gray-400 flex items-center justify-center hover:bg-red-100 hover:text-red-500 transition-colors disabled:opacity-50"
                            title="ปฏิเสธ"
                          >
                            <X size={16} />
                          </button>
                        </>
                      ) : null}
                      <button onClick={() => router.push(`/approval/${item.id}`)} className="px-4 py-1.5 rounded-lg border border-gray-200 text-xs font-semibold text-gray-600 hover:bg-gray-50">รายละเอียด</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}