'use client'
import { useState, useEffect, useRef } from 'react'
import { X, Check, Calendar, AlertTriangle, Trophy, Loader2, Plus, FileText, GraduationCap, Power } from 'lucide-react'
import api from '../../../lib/axios'
import '../styles/users.css' 

interface AcademicTerm {
  id: number
  semester: string
  academic_year: number
  start_date: string
  end_date: string
  is_open: boolean
}

interface AwardCategory {
  id: number
  name: string
  description: string
  form_structure: any[]
  is_active: boolean
  is_repeatable?: boolean
}

export default function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [refreshKey, setRefreshKey] = useState(0)
  const [togglingId, setTogglingId] = useState<number | null>(null)

  // ── States ──
  const [terms, setTerms] = useState<AcademicTerm[]>([])
  const [categories, setCategories] = useState<AwardCategory[]>([])
  const [requests, setRequests] = useState<any[]>([]) 
  const [stats, setStats] = useState({ totalRequests: 0, pendingRequests: 0, activeCategories: 0 })

  // ── Refs สำหรับการทำ Scroll ──
  const requestsRef = useRef<HTMLDivElement>(null)
  const termsRef = useRef<HTMLDivElement>(null)
  const categoriesRef = useRef<HTMLDivElement>(null)

  // Modal State
  const [showTermModal, setShowTermModal] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [newTerm, setNewTerm] = useState({
    semester: 'first', 
    academic_year: new Date().getFullYear() , 
    start_date: '',
    end_date: '',
    is_open: false
  })

  const formatDisplayDate = (dateString: string): string => {
    if (!dateString) return '—'
    const date = new Date(dateString)
    return date.toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: '2-digit' }) // ปรับให้แสดงวันที่แบบไทย
  }

  
  const scrollToSection = (ref: React.RefObject<HTMLDivElement | null>) => {
    if (ref.current) {
      ref.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        const [termsRes, catsRes, reqsRes] = await Promise.all([
          api.get('/api/v1/academic-terms'),
          api.get('/api/v1/award-categories'),
          api.get('/api/v1/requests')
        ])

        const termsData = termsRes.data.data || termsRes.data || []
        const rawCatsData = catsRes.data.data || catsRes.data || []
        
        let reqsData = reqsRes.data.data || reqsRes.data;
        if (!Array.isArray(reqsData)) {
          reqsData = []; 
        }

        const catsData = rawCatsData.map((cat: any) => ({
          ...cat,
          form_structure: typeof cat.form_structure === 'string'
            ? JSON.parse(cat.form_structure)
            : (cat.form_structure || [])
        }))

        setTerms(Array.isArray(termsData) ? termsData.sort((a: any, b: any) => b.id - a.id) : [])
        setCategories(catsData)
        
        setRequests([...reqsData].sort((a: any, b: any) => b.id - a.id))

        const activeCats = catsData.filter((c: AwardCategory) => c.is_active).length
        const pendingReqs = reqsData.filter((r: any) => r.current_status?.includes('PENDING')).length

        setStats({ totalRequests: reqsData.length, pendingRequests: pendingReqs, activeCategories: activeCats })
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [refreshKey])

  const handleCreateTerm = async () => {
    if (!newTerm.semester || !newTerm.academic_year || !newTerm.start_date || !newTerm.end_date) {
      return alert('กรุณากรอกข้อมูลให้ครบถ้วน')
    }

    try {
      setIsSubmitting(true)
      const payload = {
        academic_year: Number(newTerm.academic_year),
        semester: newTerm.semester, 
        start_date: newTerm.start_date,
        end_date: newTerm.end_date,    
        is_open: newTerm.is_open
      }

      await api.post('/api/v1/academic-terms', payload)

      alert('เพิ่มภาคการศึกษาใหม่สำเร็จ!')
      setShowTermModal(false)
      setNewTerm({
        semester: 'first',
        academic_year: new Date().getFullYear(),
        start_date: '',
        end_date: '',
        is_open: false
      })
      setRefreshKey(prev => prev + 1)

    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.response?.data?.message || 'เกิดข้อผิดพลาด'
      alert(`ไม่สามารถเพิ่มภาคการศึกษาได้:\n${errorMsg}`)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleToggleTermStatus = async (term: AcademicTerm, newStatus: boolean) => {
    const actionText = newStatus ? "เปิด" : "ปิด"
    
    const semesterDisplay = term.semester === 'first' ? 'ภาคต้น' : term.semester === 'second' ? 'ภาคปลาย' : term.semester;
    
    if (!confirm(`คุณต้องการ ${actionText} ภาคการศึกษา ${semesterDisplay}/${term.academic_year} ใช่หรือไม่?`)) return
    try {
      setTogglingId(term.id)
      await api.patch(`/api/v1/academic-terms/${term.id}`, { is_open: newStatus })
      setRefreshKey(prev => prev + 1)
    } catch (error: any) {
      alert(`เกิดข้อผิดพลาด: ${error.response?.data?.error || error.response?.data?.message || 'ทำรายการไม่สำเร็จ'}`)
    } finally {
      setTogglingId(null)
    }
  }

  const handleToggleCategory = async (category: AwardCategory) => {
    setTogglingId(category.id)
    try {
      const payload = {
        name: category.name,
        description: category.description,
        is_repeatable: category.is_repeatable || false,
        form_structure: category.form_structure,
        is_active: !category.is_active
      }
      await api.put(`/api/v1/award-categories/${category.id}`, payload)
      setRefreshKey(prev => prev + 1)
    } catch (error: any) {
      alert(`ไม่สามารถเปลี่ยนสถานะประเภทรางวัลได้: ${error.response?.data?.message || 'เกิดข้อผิดพลาด'}`)
    } finally {
      setTogglingId(null)
    }
  }

  const statsData = [
    {
      label: 'ใบคำร้องทั้งหมด', value: requests.length.toString(), change: 'คำร้องในระบบทั้งหมด', changeType: 'neutral',
      IconComponent: FileText, iconBg: 'rgba(34, 197, 94, 0.15)', iconColor: '#22c55e',
      targetRef: requestsRef
    },
    {
      label: 'ภาคการศึกษา ', value: terms.length.toString(), change: `เปิดใช้งาน ${terms.filter(t => t.is_open).length} ภาค`, changeType: 'neutral',
      IconComponent: Calendar, iconBg: 'rgba(245, 158, 11, 0.15)', iconColor: '#f59e0b',
      targetRef: termsRef
    },
    {
      label: 'ประเภทรางวัล ', value: categories.length.toString(), change: `เปิดใช้งาน ${categories.filter(c => c.is_active).length} ประเภท`, changeType: 'neutral',
      IconComponent: Trophy, iconBg: 'rgba(59, 130, 246, 0.15)', iconColor: '#3b82f6',
      targetRef: categoriesRef
    }
  ]

  if (loading) {
    return (
      <div className="users-container" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 12 }}>
          <Loader2 className="w-8 h-8 animate-spin" style={{ color: '#22c55e' }} />
          <p style={{ fontSize: 14, color: '#64748b' }}>กำลังโหลดแดชบอร์ด...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="users-container">

      {/* Page Header  */}
      <div className="users-header">
        <div className="header-content">
          <h1 className="page-title">แดชบอร์ดระบบ (System Dashboard)</h1>
        </div>
        <div className="header-actions">
          <button className="btn btn-primary" onClick={() => setShowTermModal(true)}>
            <span className="btn-icon"><Plus size={18} /></span>
            เพิ่มภาคการศึกษาใหม่
          </button>
        </div>
      </div>

      {/*  Status Cards  */}
      <div className="stats-grid">
        {statsData.map((stat, index) => {
          const StatIcon = stat.IconComponent
          return (
            <div 
              key={index} 
              className="stat-card" 
              style={{ animationDelay: `${index * 0.1}s`, cursor: 'pointer', transition: 'transform 0.2s' }}
              onClick={() => scrollToSection(stat.targetRef)} // 🔥 กดแล้ววิ่งไปที่ตาราง
              onMouseEnter={(e) => e.currentTarget.style.transform = 'translateY(-4px)'}
              onMouseLeave={(e) => e.currentTarget.style.transform = 'translateY(0)'}
            >
              <div className="stat-icon-wrapper" style={{ background: stat.iconBg }}>
                <span className="stat-icon">
                  <StatIcon size={24} color={stat.iconColor} strokeWidth={2} />
                </span>
              </div>
              <div className="stat-content">
                <div className="stat-label">{stat.label}</div>
                <div className="stat-value">{stat.value}</div>
                <div className={`stat-change ${stat.changeType}`}>{stat.change}</div>
              </div>
            </div>
          )
        })}
      </div>

      {/* ── Requests Table ── */}
      <div ref={requestsRef} className="users-header" style={{ marginTop: 32, marginBottom: 16, scrollMarginTop: '80px' }}>
        <div className="header-content" style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <FileText size={20} style={{ color: '#64748b' }} />
          <h2 style={{ fontSize: 18, fontWeight: 700, color: '#0f172a' }}>ใบคำร้อง (Requests)</h2>
        </div>
      </div>

      <div className="table-container">
        <table className="users-table">
          <thead>
            <tr>
              <th>รหัสอ้างอิง</th>
              <th>ผู้สมัคร</th>
              <th>ประเภทรางวัล</th>
              <th style={{ textAlign: 'center' }}>สถานะ</th>
              <th style={{ textAlign: 'center' }}>วันที่สมัคร</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {requests.length === 0 ? (
              <tr>
                <td colSpan={5} className="py-20 text-center" style={{ color: '#64748b' }}>
                  ยังไม่มีข้อมูลใบคำร้อง
                </td>
              </tr>
            ) : requests.map((req) => {

              const studentName = req.student ? `${req.student.first_name || ''} ${req.student.last_name || ''}` : 'ไม่ระบุชื่อผู้สมัคร';
              const studentIdOrEmail = req.student?.nisit_id || req.student?.email || '-';
              const awardName = req.award_name || req.award?.name || 'ไม่ระบุประเภทรางวัล';
              
              // ฟังก์ชันแปลสถานะ
              const translateStatus = (status: string) => {
                const s = status ? status.toUpperCase() : "";
                if (s === "SUBMITTED") return "ส่งใบสมัครแล้ว";
                if (s === "APPROVED_PENDING_PRESIDENT") return "อนุมัติ";
                if (s.includes("PENDING")) return "รอดำเนินการ";
                if (s === "APPROVED") return "อนุมัติแล้ว";
                if (s.includes("REJECTED")) return "ไม่ผ่านการอนุมัติ";
                if (s === "AWARD_CHANGED") return "เปลี่ยนประเภทรางวัล";
                return status ? status.replace(/_/g, ' ') : "รอดำเนินการ"; 
              };

              const statusDisplay = translateStatus(req.current_status);

              return (
                <tr key={req.id}>
                  <td>
                    <span style={{ fontWeight: 600, color: '#0f172a' }}>REQ-{String(req.id).padStart(4, '0')}</span>
                  </td>
                  <td>
                    <div className="user-name font-bold text-gray-800 text-sm">
                      {studentName}
                    </div>
                    <div className="user-email text-xs text-gray-500 mt-0.5">
                      {studentIdOrEmail}
                    </div>
                  </td>
                  <td className="department-col text-sm text-gray-600 font-medium">
                    {awardName}
                  </td>
                  <td style={{ textAlign: 'center' }}>
                    <span className={`status-badge px-3 py-1 rounded-full text-xs font-bold ${
                      statusDisplay === "อนุมัติแล้ว" ? "bg-green-50 text-green-600" :
                      statusDisplay.includes("ไม่ผ่าน") ? "bg-red-50 text-red-600" :
                      statusDisplay.includes("เปลี่ยน") ? "bg-orange-50 text-orange-600" :
                      "bg-yellow-50 text-yellow-600"
                    }`}>
                      {statusDisplay}
                    </span>
                  </td>
                  <td style={{ textAlign: 'center' }}>{formatDisplayDate(req.created_at)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/*  Academic Terms Table  */}
      <div ref={termsRef} className="users-header" style={{ marginTop: 32, marginBottom: 16, scrollMarginTop: '80px' }}>
        <div className="header-content" style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Calendar size={20} style={{ color: '#64748b' }} />
          <h2 style={{ fontSize: 18, fontWeight: 700, color: '#0f172a' }}>ภาคการศึกษา (Terms)</h2>
        </div>
      </div>

      <div className="table-container">
        <table className="users-table">
          <thead>
            <tr>
              <th>ภาค / ปี</th>
              <th>วันเริ่มต้น</th>
              <th>วันสิ้นสุด</th>
              <th style={{ textAlign: 'center' }}>สถานะ</th>
              <th style={{ textAlign: 'center' }}>จัดการ</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {terms.length === 0 ? (
              <tr>
                <td colSpan={5} className="py-20 text-center" style={{ color: '#64748b' }}>
                  ยังไม่มีข้อมูลภาคการศึกษา
                </td>
              </tr>
            ) : terms.map((term) => (
              <tr key={term.id}>
                <td>
                  <div className="user-details">
                    <div
                      className="stat-icon-wrapper"
                      style={{
                        width: 40, height: 40,
                        background: term.is_open ? 'rgba(34,197,94,0.12)' : '#f1f5f9',
                        borderRadius: 10, flexShrink: 0
                      }}
                    >
                      <GraduationCap size={18} color={term.is_open ? '#22c55e' : '#64748b'} />
                    </div>
                    <div className="user-info">

                      <div className="user-name">
                        {term.semester === 'first' ? 'ภาคต้น' : term.semester === 'second' ? 'ภาคปลาย' : term.semester}
                      </div>
                      <div className="user-email">ปีการศึกษา {term.academic_year}</div>
                    </div>
                  </div>
                </td>
                <td className="department-col">{formatDisplayDate(term.start_date)}</td>
                <td className="department-col">{formatDisplayDate(term.end_date)}</td>
                <td style={{ textAlign: 'center' }}>
                  <span
                    className="status-badge"
                    style={{
                      backgroundColor: term.is_open ? 'rgba(34,197,94,0.12)' : '#f1f5f9',
                      color: term.is_open ? '#22c55e' : '#64748b'
                    }}
                  >
                    ● {term.is_open ? 'เปิด' : 'ปิด'}
                  </span>
                </td>
                <td style={{ textAlign: 'center' }}>
                  <div className="action-buttons" style={{ justifyContent: 'center' }}>
                    <button
                      className="action-btn"
                      title={term.is_open ? 'ปิดภาคการศึกษา' : 'เปิดภาคการศึกษา'}
                      disabled={togglingId === term.id}
                      onClick={() => handleToggleTermStatus(term, !term.is_open)}
                      style={{
                        color: term.is_open ? '#ef4444' : '#22c55e',
                        borderColor: term.is_open ? '#fecaca' : '#bbf7d0',
                        width: 'auto', padding: '0 12px', gap: 6, fontSize: 12, fontWeight: 600
                      }}
                    >
                      {togglingId === term.id
                        ? <Loader2 size={14} className="animate-spin" />
                        : <Power size={14} />
                      }
                      {term.is_open ? 'ปิดเทอม' : 'เปิดเทอม'}
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* ── Award Categories Table ── */}
      <div ref={categoriesRef} className="users-header" style={{ marginTop: 32, marginBottom: 16, scrollMarginTop: '80px' }}>
        <div className="header-content" style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Trophy size={20} style={{ color: '#64748b' }} />
          <h2 style={{ fontSize: 18, fontWeight: 700, color: '#0f172a' }}>ประเภทรางวัล (Categories)</h2>
        </div>
      </div>

      <div className="table-container">
        <table className="users-table">
          <thead>
            <tr>
              <th>ชื่อรางวัล</th>
              <th>คำอธิบาย</th>
              <th style={{ textAlign: 'center' }}>สถานะ</th>
              <th style={{ textAlign: 'center' }}>เปิด/ปิด</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {categories.length === 0 ? (
              <tr>
                <td colSpan={4} className="py-20 text-center" style={{ color: '#64748b' }}>
                  ยังไม่มีข้อมูลประเภทรางวัล
                </td>
              </tr>
            ) : categories.map((cat) => (
              <tr key={cat.id}>
                <td>
                  <div className="user-name">{cat.name}</div>
                </td>
                <td className="department-col" style={{ maxWidth: 300 }}>
                  <span style={{ display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                    {cat.description}
                  </span>
                </td>
                <td style={{ textAlign: 'center' }}>
                  <span
                    className="status-badge"
                    style={{
                      backgroundColor: cat.is_active ? 'rgba(34,197,94,0.12)' : '#f1f5f9',
                      color: cat.is_active ? '#22c55e' : '#64748b'
                    }}
                  >
                    ● {cat.is_active ? 'ใช้งาน' : 'ระงับการใช้งาน'}
                  </span>
                </td>
                <td style={{ textAlign: 'center' }}>
                  <div className="action-buttons" style={{ justifyContent: 'center' }}>
                    <button
                      onClick={() => handleToggleCategory(cat)}
                      disabled={togglingId === cat.id}
                      style={{
                        position: 'relative', display: 'inline-flex', height: 24, width: 44,
                        alignItems: 'center', borderRadius: 9999, border: 'none', cursor: 'pointer',
                        backgroundColor: cat.is_active ? '#22c55e' : '#cbd5e1',
                        opacity: togglingId === cat.id ? 0.5 : 1,
                        transition: 'background-color 0.2s'
                      }}
                    >
                      <span style={{
                        display: 'inline-block', width: 16, height: 16, borderRadius: '50%',
                        backgroundColor: 'white',
                        transform: cat.is_active ? 'translateX(24px)' : 'translateX(4px)',
                        transition: 'transform 0.2s'
                      }} />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* ── Term Creation Modal ── */}
      {showTermModal && (
        <div className="modal-overlay" onClick={() => setShowTermModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">เพิ่มภาคการศึกษาใหม่</h2>
              <button className="modal-close" onClick={() => setShowTermModal(false)}>
                <X size={20} />
              </button>
            </div>

   
            <form 
              className="user-form" 
              onSubmit={(e) => {
                e.preventDefault();
                handleCreateTerm();
              }}
            >
              
              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">ภาคการศึกษา <span className="text-red-500">*</span></label>
                  <select
                    value={newTerm.semester}
                    onChange={(e) => setNewTerm({ ...newTerm, semester: e.target.value })}
                    className="form-input"
                    required
                  >
                    <option value="first">ภาคต้น</option>
                    <option value="second">ภาคปลาย</option>
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">ปีการศึกษา <span className="text-red-500">*</span></label>
                  <input
                    type="number"
                    value={newTerm.academic_year}
                    onChange={(e) => setNewTerm({ ...newTerm, academic_year: Number(e.target.value) })}
                    className="form-input"
                    required
                  />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">วันเริ่มต้น <span className="text-red-500">*</span></label>
                  <input
                    type="date"
                    value={newTerm.start_date}
                    onChange={(e) => setNewTerm({ ...newTerm, start_date: e.target.value })}
                    className="form-input"
                    required
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">วันสิ้นสุด <span className="text-red-500">*</span></label>
                  <input
                    type="date"
                    value={newTerm.end_date}
                    onChange={(e) => setNewTerm({ ...newTerm, end_date: e.target.value })}
                    className="form-input"
                    required
                  />
                </div>
              </div>

              {/* Toggle เปิดใช้งานทันที */}
              <div className="form-group" style={{ marginTop: '8px' }}>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '16px', backgroundColor: '#f8fafc', borderRadius: '12px', border: '1px solid #e2e8f0' }}>
                  <div>
                    <p style={{ fontSize: '14px', fontWeight: 600, color: '#0f172a', marginBottom: '4px' }}>เปิดใช้งานทันที</p>
                    <p style={{ fontSize: '12px', color: '#64748b' }}>หากเปิด นิสิตจะสามารถเริ่มส่งคำร้องในเทอมนี้ได้ทันที</p>
                  </div>
                  <button
                    type="button"
                    onClick={() => setNewTerm({ ...newTerm, is_open: !newTerm.is_open })}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ${newTerm.is_open ? "bg-green-500" : "bg-gray-300"}`}
                  >
                    <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${newTerm.is_open ? "translate-x-6" : "translate-x-1"}`} />
                  </button>
                </div>
              </div>

              {/* ปุ่ม Actions */}
              <div className="form-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary" 
                  onClick={() => setShowTermModal(false)}
                >
                  ยกเลิก
                </button>
                <button 
                  type="submit" 
                  disabled={isSubmitting} 
                  className="btn btn-primary" 
                  style={{ opacity: isSubmitting ? 0.7 : 1 }}
                >
                  {isSubmitting ? <Loader2 size={16} className="animate-spin" /> : <Check size={16}/>}
                  {isSubmitting ? "กำลังสร้าง..." : "ยืนยันการเพิ่ม"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  )
}