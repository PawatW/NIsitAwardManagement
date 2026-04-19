'use client'

import { useState, useEffect, useMemo } from 'react'
import { 
  Search, Settings as SettingsIcon, Menu, Edit, History, User as UserIcon, 
  UserPlus, Check, X, ClipboardList, Users as UsersIcon, Briefcase, Loader2, Power
} from 'lucide-react'
import api from '../../../../lib/axios'
import { useAuth } from '../../../../context/AuthContext'
import '../../styles/users.css'

interface User {
  id: string
  name: string
  email: string
  role: string
  roleThai: string
  roleColor: string
  department: string
  status: string
  statusColor: string
  isActive: boolean
  avatar: string
  raw: any 
}

interface OrganizationItem {
  ID: string;
  Name: string;
}

interface NewUserData {
  first_name: string
  last_name: string
  email: string
  phone_number: string
  role: string
  campus_id: string
  faculty_id: string
  department_id: string
}

export default function Users() {
  const { user: currentUser } = useAuth()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [refreshKey, setRefreshKey] = useState(0) 
  const [togglingId, setTogglingId] = useState<string | null>(null)

  const [activeTab, setActiveTab] = useState<string>('ผู้ใช้งานทั้งหมด')
  const [searchQuery, setSearchQuery] = useState<string>('')
  
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]) 
  const [selectAll, setSelectAll] = useState<boolean>(false)
  
  const [showAddUserModal, setShowAddUserModal] = useState<boolean>(false)
  const [showProfileModal, setShowProfileModal] = useState<boolean>(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)

  const [campuses, setCampuses] = useState<OrganizationItem[]>([])
  const [faculties, setFaculties] = useState<OrganizationItem[]>([])
  const [departments, setDepartments] = useState<OrganizationItem[]>([])

  const [newUserData, setNewUserData] = useState<NewUserData>({
    first_name: '',
    last_name: '',
    email: '',
    phone_number: '',
    role: 'HEAD_OF_DEPARTMENT', 
    campus_id: '',
    faculty_id: '',
    department_id: ''
  })

  const tabs = ['ผู้ใช้งานทั้งหมด', 'นักศึกษา', 'คณะ', 'สาขา', 'คณะกรรมการ']

  const getRoleThaiName = (role: string) => {
    if (role.includes('ADMIN')) return 'ผู้ดูแลระบบ'
    if (role.includes('STUDENT')) return 'นิสิต'
    if (role.includes('COMMITTEE')) return 'คณะกรรมการ'
    if (role.includes('VICE_DEAN')) return 'รองคณบดี'
    if (role.includes('DEAN')) return 'คณบดี'
    if (role.includes('HEAD')) return 'หัวหน้าภาควิชา'
    return role.replace(/_/g, ' ')
  }

  const getRoleColor = (role: string) => {
    if (role.includes('ADMIN')) return '#ef4444' 
    if (role.includes('STUDENT')) return '#3b82f6' 
    if (role.includes('COMMITTEE')) return '#a855f7' 
    if (role.includes('DEAN') || role.includes('HEAD')) return '#f59e0b' 
    return '#64748b' 
  }

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setLoading(true)
        const res = await api.get('/api/v1/admin/users')
        const rawData = res.data.data || res.data || []

        const mappedUsers: User[] = rawData.map((u: any) => {
          const roleRaw = u.Role || u.role || 'UNKNOWN'
          const roleStr = roleRaw.toUpperCase()
          
          const deptName = u.Department?.Name || u.department?.name || 
                           u.Faculty?.Name || u.faculty?.name || 
                           u.Campus?.Name || u.campus?.name || '-'

          const isActive = u.is_active !== undefined ? u.is_active : (u.IsActive !== undefined ? u.IsActive : true)

          return {
            id: u.UserID || u.id,
            name: `${u.FirstName || u.first_name || ''} ${u.LastName || u.last_name || ''}`.trim() || 'ไม่ระบุชื่อ',
            email: u.Email || u.email || '-',
            role: roleStr,
            roleThai: getRoleThaiName(roleStr),
            roleColor: getRoleColor(roleStr),
            department: deptName,
            isActive: isActive,
            status: isActive ? 'ใช้งานปกติ' : 'ระงับการใช้งาน',
            statusColor: isActive ? '#22c55e' : '#94a3b8',
            avatar: u.ProfileURL || u.profile_url || `https://ui-avatars.com/api/?name=${u.FirstName || u.first_name}+${u.LastName || u.last_name}&background=random&color=fff`,
            raw: u
          }
        })
        setUsers(mappedUsers)
      } catch (error) {
        console.error("Failed to fetch users:", error)
      } finally {
        setLoading(false)
      }
    }

    if (currentUser?.role === 'ADMIN') {
        fetchUsers()
    } else {
        setLoading(false) 
    }
  }, [currentUser, refreshKey])

  useEffect(() => {
    if (showAddUserModal) {
      api.get("/api/v1/organization/campuses")
         .then(res => setCampuses(res.data.data || res.data || []))
         .catch(err => console.error(err))
    }
  }, [showAddUserModal])

  useEffect(() => {
    if (newUserData.campus_id) {
      api.get(`/api/v1/organization/campuses/${newUserData.campus_id}/faculties`)
         .then(res => setFaculties(res.data.data || res.data || []))
         .catch(err => console.error(err))
    } else {
      setFaculties([])
    }
  }, [newUserData.campus_id])

  useEffect(() => {
    if (newUserData.faculty_id) {
      api.get(`/api/v1/organization/faculties/${newUserData.faculty_id}/departments`)
         .then(res => setDepartments(res.data.data || res.data || []))
         .catch(err => console.error(err))
    } else {
      setDepartments([])
    }
  }, [newUserData.faculty_id])

  const handleNewUserInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setNewUserData(prev => {
      const updated = { ...prev, [name]: value }
      if (name === "campus_id") {
        updated.faculty_id = "";
        updated.department_id = "";
      }
      if (name === "faculty_id") {
        updated.department_id = "";
      }
      if (name === "role") {
        if (value === 'ADMIN') {
          updated.campus_id = "";
          updated.faculty_id = "";
          updated.department_id = "";
        } else if (value === 'COMMITTEE_CHAIR') {
          updated.faculty_id = "";
          updated.department_id = "";
        } else if (value === 'DEAN' || value === 'VICE_DEAN') {
          updated.department_id = "";
        }
      }
      return updated;
    })
  }

  const needsCampus = newUserData.role !== 'ADMIN';
  const needsFaculty = ['HEAD_OF_DEPARTMENT', 'VICE_DEAN', 'DEAN'].includes(newUserData.role);
  const needsDepartment = ['HEAD_OF_DEPARTMENT'].includes(newUserData.role);

  const handleAddUserSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsSubmitting(true)
    
    try {
        const payload = { ...newUserData }

        if (!needsCampus) delete (payload as any).campus_id;
        if (!needsFaculty) delete (payload as any).faculty_id;
        if (!needsDepartment) delete (payload as any).department_id;
        
        await api.post('/api/v1/admin/users', payload)
        
        alert('สร้างผู้ใช้งานสำเร็จ!')
        setShowAddUserModal(false)
        setRefreshKey(prev => prev + 1) 

        setNewUserData({
          first_name: '', last_name: '', email: '', phone_number: '', 
          role: 'HEAD_OF_DEPARTMENT', campus_id: '', faculty_id: '', department_id: ''
        })
    } catch (error: any) {
        console.error(error)
        alert(`เกิดข้อผิดพลาด: ${error.response?.data?.message || "กรุณาตรวจสอบข้อมูลอีกครั้ง"}`)
    } finally {
        setIsSubmitting(false)
    }
  }
  
  const handleToggleUserStatus = async (user: User) => {
    const actionText = user.isActive ? "ระงับการใช้งาน" : "เปิดการใช้งาน";
    if (!confirm(`คุณต้องการ ${actionText} ของคุณ ${user.name} ใช่หรือไม่?`)) return;

    try {
      setTogglingId(user.id);
      await api.patch(`/api/v1/admin/users/${user.id}/status`, { is_active: !user.isActive });
      setRefreshKey(prev => prev + 1); 
    } catch (error: any) {
      console.error(error);
      alert(`ไม่สามารถเปลี่ยนสถานะได้: ${error.response?.data?.message || "เกิดข้อผิดพลาด"}`);
    } finally {
      setTogglingId(null);
    }
  };

  const filteredUsers = useMemo(() => {
    return users.filter(u => {
      const matchesSearch = u.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                            u.email.toLowerCase().includes(searchQuery.toLowerCase())
      let matchesTab = true
      if (activeTab === 'นักศึกษา') matchesTab = u.role.includes('STUDENT')
      else if (activeTab === 'คณะ') matchesTab = u.role.includes('DEAN') || u.role.includes('VICE_DEAN')
      else if (activeTab === 'สาขา') matchesTab = u.role.includes('HEAD')
      else if (activeTab === 'คณะกรรมการ') matchesTab = u.role.includes('COMMITTEE')

      return matchesSearch && matchesTab
    })
  }, [users, searchQuery, activeTab])

  const statsData = [
    {
      label: 'ผู้ใช้งานทั้งหมด', value: users.length.toString(), change: 'ผู้ใช้งานในระบบ', changeType: 'neutral',
      IconComponent: UsersIcon, iconBg: 'rgba(34, 197, 94, 0.15)'
    },
    {
      label: 'นักศึกษา', value: users.filter(u => u.role.includes('STUDENT')).length.toString(), change: 'นักศึกษาที่ลงทะเบียน', changeType: 'neutral',
      IconComponent: ClipboardList, iconBg: 'rgba(59, 130, 246, 0.15)' 
    },
    {
      label: 'บุคลากร', value: users.filter(u => u.role.includes('STAFF') || u.role.includes('COMMITTEE') || u.role.includes('DEAN') || u.role.includes('HEAD')).length.toString(), change: 'ผู้บริหารและเจ้าหน้าที่', changeType: 'neutral',
      IconComponent: Briefcase, iconBg: 'rgba(245, 158, 11, 0.15)' 
    }
  ]

  const handleSelectAll = () => {
    if (selectAll) setSelectedUsers([])
    else setSelectedUsers(filteredUsers.map(user => user.id))
    setSelectAll(!selectAll)
  }

  const handleSelectUser = (userId: string) => {
    if (selectedUsers.includes(userId)) {
      setSelectedUsers(selectedUsers.filter(id => id !== userId))
      setSelectAll(false)
    } else {
      const newSelected = [...selectedUsers, userId]
      setSelectedUsers(newSelected)
      if (newSelected.length === filteredUsers.length) setSelectAll(true)
    }
  }

  return (
    <div className="users-container">
      {/* Page Header */}
      <div className="users-header">
        <div className="header-content">
          <h1 className="page-title">จัดการผู้ใช้งาน</h1>
        </div>
        <div className="header-actions">
          <button className="btn btn-primary" onClick={() => setShowAddUserModal(true)}>
            <span className="btn-icon"><UserPlus size={18} /></span>
            เพิ่มผู้ใช้งาน
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="stats-grid">
        {statsData.map((stat, index) => {
          const StatIcon = stat.IconComponent
          const iconColors = ['#22c55e', '#3b82f6', '#f59e0b']
          return (
            <div key={index} className="stat-card" style={{ animationDelay: `${index * 0.1}s` }}>
              <div className="stat-icon-wrapper" style={{ background: stat.iconBg }}>
                <span className="stat-icon">
                  <StatIcon size={24} color={iconColors[index]} strokeWidth={2} />
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

      {/* Search and Filters */}
      <div className="table-controls">
        <div className="search-box">
          <span className="search-icon"><Search size={18} /></span>
          <input
            type="text"
            placeholder="ค้นหาผู้ใช้งาน..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs-container">
        {tabs.map((tab) => (
          <button key={tab} className={`tab ${activeTab === tab ? 'active' : ''}`} onClick={() => setActiveTab(tab)}>
            {tab}
          </button>
        ))}
      </div>

      {/* Users Table */}
      <div className="table-container">
        <table className="users-table">
          <thead>
            <tr>
              <th className="checkbox-col">
                <input type="checkbox" checked={selectAll} onChange={handleSelectAll} />
              </th>
              <th>รายละเอียดผู้ใช้งาน</th>
              <th>ตำแหน่ง</th>
              <th>สังกัด</th>
              <th>สถานะ</th>
              <th style={{ textAlign: 'right' }}>จัดการ</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {loading ? (
                <tr><td colSpan={6} className="py-20 text-center"><Loader2 className="w-8 h-8 animate-spin mx-auto text-green-600" /></td></tr>
            ) : filteredUsers.length === 0 ? (
                <tr><td colSpan={6} className="py-20 text-center text-gray-500">ไม่พบข้อมูลผู้ใช้งาน</td></tr>
            ) : filteredUsers.map((user) => (
              <tr key={user.id} className={!user.isActive ? 'opacity-60 bg-gray-50' : ''}>
                <td className="checkbox-col">
                  <input type="checkbox" checked={selectedUsers.includes(user.id)} onChange={() => handleSelectUser(user.id)} />
                </td>
                <td>
                  <div className="user-details">
                    <img src={user.avatar} alt={user.name} className={`user-avatar ${!user.isActive ? 'grayscale' : ''}`} />
                    <div className="user-info">
                      <div className="user-name">{user.name}</div>
                      <div className="user-email">{user.email}</div>
                    </div>
                  </div>
                </td>
                <td>
                  <span className="role-badge" style={{ backgroundColor: `${user.roleColor}15`, color: user.roleColor }}>
                    {user.roleThai}
                  </span>
                </td>
                <td className="department-col">{user.department}</td>
                <td>
                  <span className="status-badge" style={{ backgroundColor: `${user.statusColor}15`, color: user.statusColor }}>
                    ● {user.status}
                  </span>
                </td>
                <td>
                  <div className="action-buttons" style={{ justifyContent: 'flex-end', gap: '8px' }}>
                    <button className="action-btn profile-btn" title="ดูข้อมูลส่วนตัว" onClick={() => { setSelectedUser(user); setShowProfileModal(true); }}>
                      <UserIcon size={18} />
                    </button>
                    
            
                    <button
                      type="button"
                      onClick={() => handleToggleUserStatus(user)}
                      disabled={togglingId === user.id}
                      title={user.isActive ? "ระงับการใช้งาน" : "เปิดการใช้งาน"}
                      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors disabled:opacity-50 focus:outline-none ${
                        user.isActive ? "bg-green-500" : "bg-gray-300"
                      }`}
                      style={{ cursor: togglingId === user.id ? 'wait' : 'pointer' }}
                    >
                      <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform flex items-center justify-center ${
                          user.isActive ? "translate-x-6" : "translate-x-1"
                        }`}
                      >
                        {togglingId === user.id && <Loader2 size={12} className="animate-spin text-gray-500" />}
                      </span>
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {showAddUserModal && (
        <div className="modal-overlay" onClick={() => setShowAddUserModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
             <div className="modal-header">
               <h2 className="modal-title">เพิ่มผู้ใช้งานใหม่</h2>
               <button className="modal-close" onClick={() => setShowAddUserModal(false)}><X size={20}/></button>
             </div>
             
             <form onSubmit={handleAddUserSubmit} className="user-form">
                
                <div className="form-row">
                  <div className="form-group">
                    <label className="form-label">ชื่อ <span className="text-red-500">*</span></label>
                    <input type="text" name="first_name" value={newUserData.first_name} onChange={handleNewUserInputChange} className="form-input" required />
                  </div>
                  <div className="form-group">
                    <label className="form-label">นามสกุล <span className="text-red-500">*</span></label>
                    <input type="text" name="last_name" value={newUserData.last_name} onChange={handleNewUserInputChange} className="form-input" required />
                  </div>
                </div>

                <div className="form-group">
                  <label className="form-label">อีเมล <span className="text-red-500">*</span></label>
                  <input type="email" name="email" value={newUserData.email} onChange={handleNewUserInputChange} className="form-input" required />
                </div>

                <div className="form-group">
                  <label className="form-label">เบอร์โทรศัพท์ <span className="text-red-500">*</span></label>
                  <input type="tel" name="phone_number" value={newUserData.phone_number} onChange={handleNewUserInputChange} className="form-input" required />
                </div>

                <div className="form-group">
                  <label className="form-label">ตำแหน่งในระบบ <span className="text-red-500">*</span></label>
                  {/* ลบ className.form-select ออกเพื่อไม่ให้ไอคอนลูกศรมาทับด้านซ้ายตาม CSS 原ฉบับ */}
                  <select name="role" value={newUserData.role} onChange={handleNewUserInputChange} className="form-input" required>
                    <option value="HEAD_OF_DEPARTMENT">หัวหน้าภาควิชา (HOD)</option>
                    <option value="VICE_DEAN">รองคณบดี (Vice Dean)</option>
                    <option value="DEAN">คณบดี (Dean)</option>
                    <option value="COMMITTEE_CHAIR">คณะกรรมการ (Committee)</option>
                    <option value="ADMIN">ผู้ดูแลระบบ (Admin)</option>
                  </select>
                </div>

                {needsCampus && (
                  <div className="form-group">
                    <label className="form-label">วิทยาเขต <span className="text-red-500">*</span></label>
                    <select name="campus_id" value={newUserData.campus_id} onChange={handleNewUserInputChange} className="form-input" required>
                      <option value="">-- เลือกวิทยาเขต --</option>
                      {campuses.map(c => <option key={c.ID} value={c.ID}>{c.Name}</option>)}
                    </select>
                  </div>
                )}

                {needsFaculty && (
                  <div className="form-group">
                    <label className="form-label">คณะ <span className="text-red-500">*</span></label>
                    <select name="faculty_id" value={newUserData.faculty_id} onChange={handleNewUserInputChange} className="form-input" required disabled={!newUserData.campus_id} style={{ backgroundColor: !newUserData.campus_id ? '#f1f5f9' : 'white' }}>
                      <option value="">{newUserData.campus_id ? "-- เลือกคณะ --" : "กรุณาเลือกวิทยาเขตก่อน"}</option>
                      {faculties.map(f => <option key={f.ID} value={f.ID}>{f.Name}</option>)}
                    </select>
                  </div>
                )}
       
                {needsDepartment && (
                  <div className="form-group">
                    <label className="form-label">ภาควิชา/สาขา <span className="text-red-500">*</span></label>
                    <select name="department_id" value={newUserData.department_id} onChange={handleNewUserInputChange} className="form-input" disabled={!newUserData.faculty_id} style={{ backgroundColor: !newUserData.faculty_id ? '#f1f5f9' : 'white' }}>
                      <option value="">{newUserData.faculty_id ? "-- เลือกภาควิชา/สาขา --" : "กรุณาเลือกคณะก่อน"}</option>
                      {departments.map(d => <option key={d.ID} value={d.ID}>{d.Name}</option>)}
                    </select>
                  </div>
                )}

                <div className="form-actions">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowAddUserModal(false)}>ยกเลิก</button>
                  <button type="submit" disabled={isSubmitting} className="btn btn-primary" style={{ opacity: isSubmitting ? 0.7 : 1 }}>
                    {isSubmitting ? <Loader2 size={16} className="animate-spin" /> : <Check size={16}/>}
                    {isSubmitting ? "กำลังบันทึก..." : "สร้างผู้ใช้งาน"}
                  </button>
                </div>
             </form>
          </div>
        </div>
      )}

      {/* Profile Modal */}
      {showProfileModal && selectedUser && (
        <div className="modal-overlay" onClick={() => setShowProfileModal(false)}>
          <div className="modal-content modal-large" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">ข้อมูลผู้ใช้งาน</h2>
              <button className="modal-close" onClick={() => setShowProfileModal(false)}><X size={20} /></button>
            </div>
            <div className="profile-content">
              <div className="profile-header">
                <div className="profile-avatar">
                  <img src={selectedUser.avatar} alt={selectedUser.name} className={!selectedUser.isActive ? 'grayscale' : ''} />
                </div>
                <div className="profile-info">
                  <h3 className="profile-name">{selectedUser.name}</h3>
                  <p className="profile-email">{selectedUser.email}</p>
                  <span className="status-badge mt-2 inline-block" style={{ backgroundColor: `${selectedUser.statusColor}15`, color: selectedUser.statusColor }}>
                    ● {selectedUser.status}
                  </span>
                </div>
              </div>
              <div className="profile-details">
                <h4 className="section-title">ข้อมูลสังกัด</h4>
                <div className="details-grid">
                  <div className="detail-item"><span className="detail-label">ตำแหน่ง</span><span className="detail-value">{selectedUser.roleThai}</span></div>
                  <div className="detail-item"><span className="detail-label">วิทยาเขต</span><span className="detail-value">{selectedUser.raw?.Campus?.Name || selectedUser.raw?.campus?.name || '-'}</span></div>
                  <div className="detail-item"><span className="detail-label">คณะ</span><span className="detail-value">{selectedUser.raw?.Faculty?.Name || selectedUser.raw?.faculty?.name || '-'}</span></div>
                  <div className="detail-item"><span className="detail-label">ภาควิชา/สาขา</span><span className="detail-value">{selectedUser.department}</span></div>
                  <div className="detail-item"><span className="detail-label">รหัสนิสิต</span><span className="detail-value">{selectedUser.raw?.NisitID || selectedUser.raw?.nisit_id || '-'}</span></div>
                  <div className="detail-item"><span className="detail-label">เบอร์โทรศัพท์</span><span className="detail-value">{selectedUser.raw?.PhoneNumber || selectedUser.raw?.phone_number || '-'}</span></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

    </div>
  )
}