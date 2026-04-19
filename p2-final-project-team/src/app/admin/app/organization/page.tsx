'use client'
import { useState, useEffect } from 'react'
import {
  Building2, GraduationCap, BookOpen, Plus, X, Check,
  Loader2, ChevronRight, Pencil, Trash2, Search, Filter
} from 'lucide-react'
import api from '../../../../lib/axios'
import '../../styles/users.css'

//  TYPES 

interface Campus {
  ID: string
  Name: string
}

interface Faculty {
  ID: string
  Name: string
  CampusID: string
  Campus?: Campus
}

interface Department {
  ID: string
  Name: string
  FacultyID: string
  Faculty?: Faculty
}

type ActiveTab = 'campus' | 'faculty' | 'department'

//  COMPONENT 

export default function OrganizationPage() {
  const [activeTab, setActiveTab] = useState<ActiveTab>('campus')
  const [searchQuery, setSearchQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [refreshKey, setRefreshKey] = useState(0)

  //  DATA STATES 
  const [campuses, setCampuses] = useState<Campus[]>([])
  const [faculties, setFaculties] = useState<Faculty[]>([])
  const [departments, setDepartments] = useState<Department[]>([])

  //  FILTER STATES 
  const [filterCampusId, setFilterCampusId] = useState<string>('')
  const [filterFacultyId, setFilterFacultyId] = useState<string>('')
  const [filterFacultiesList, setFilterFacultiesList] = useState<Faculty[]>([]) 

  // MODAL STATES 
  const [showModal, setShowModal] = useState(false)
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [form, setForm] = useState({ name: '', campus_id: '', faculty_id: '' })

  // Dropdown data for Modal Form
  const [formFaculties, setFormFaculties] = useState<Faculty[]>([])

  const tabs = [
    { key: 'campus' as ActiveTab, label: 'วิทยาเขต', icon: Building2, color: '#3b82f6' },
    { key: 'faculty' as ActiveTab, label: 'คณะ', icon: GraduationCap, color: '#8b5cf6' },
    { key: 'department' as ActiveTab, label: 'ภาควิชา', icon: BookOpen, color: '#f59e0b' },
  ]

  //  FETCH LOGIC 

  useEffect(() => {
    setLoading(true)
    api.get('/api/v1/organization/campuses')
      .then(res => setCampuses(res.data?.data || res.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [refreshKey])

  useEffect(() => {
    if (activeTab === 'faculty' && filterCampusId) {
      setLoading(true)
      api.get(`/api/v1/organization/campuses/${filterCampusId}/faculties`)
        .then(res => setFaculties(res.data?.data || res.data || []))
        .catch(console.error)
        .finally(() => setLoading(false))
    } else if (activeTab === 'faculty' && !filterCampusId) {
      setFaculties([]) 
    }
  }, [activeTab, filterCampusId, refreshKey])

  useEffect(() => {
    if (activeTab === 'department' && filterCampusId) {
      api.get(`/api/v1/organization/campuses/${filterCampusId}/faculties`)
        .then(res => setFilterFacultiesList(res.data?.data || res.data || []))
        .catch(console.error)
    } else {
      setFilterFacultiesList([])
    }
  }, [activeTab, filterCampusId])

  useEffect(() => {
    if (activeTab === 'department' && filterFacultyId) {
      setLoading(true)
      api.get(`/api/v1/organization/faculties/${filterFacultyId}/departments`)
        .then(res => setDepartments(res.data?.data || res.data || []))
        .catch(console.error)
        .finally(() => setLoading(false))
    } else if (activeTab === 'department' && !filterFacultyId) {
      setDepartments([])
    }
  }, [activeTab, filterFacultyId, refreshKey])

  useEffect(() => {
    if (form.campus_id && activeTab === 'department' && showModal) {
      api.get(`/api/v1/organization/campuses/${form.campus_id}/faculties`)
        .then(res => setFormFaculties(res.data?.data || res.data || []))
        .catch(console.error)
    } else {
      setFormFaculties([])
    }
  }, [form.campus_id, activeTab, showModal])


  const handleTabChange = (tab: ActiveTab) => {
    setActiveTab(tab)
    setSearchQuery('')
    // รีเซ็ตตัวกรองเมื่อเปลี่ยนแท็บ
    if (tab === 'campus') {
      setFilterCampusId('')
      setFilterFacultyId('')
    }
  }

  const openCreateModal = () => {
    setModalMode('create')
    setEditingId(null)
    setForm({ 
      name: '', 
      campus_id: filterCampusId || '', 
      faculty_id: filterFacultyId || '' 
    })
    setShowModal(true)
  }


  const handleSubmit = async () => {
    if (!form.name.trim()) return alert('Please enter a name.')
    if (activeTab === 'faculty' && !form.campus_id) return alert('Please select a campus.')
    if (activeTab === 'department' && (!form.campus_id || !form.faculty_id)) return alert('Please select campus and faculty.')

    setIsSubmitting(true)
    try {
      let endpoint = ''
      const payload: any = { name: form.name }

      if (activeTab === 'campus') {
        endpoint = modalMode === 'create' ? '/api/v1/admin/campuses' : `/api/v1/admin/campuses/${editingId}`
      } else if (activeTab === 'faculty') {
        payload.campus_id = form.campus_id 
        endpoint = modalMode === 'create' ? '/api/v1/admin/faculties' : `/api/v1/admin/faculties/${editingId}`
      } else if (activeTab === 'department') {
        payload.faculty_id = form.faculty_id 
        endpoint = modalMode === 'create' ? '/api/v1/admin/departments' : `/api/v1/admin/departments/${editingId}`
      }

      if (modalMode === 'create') {
        await api.post(endpoint, payload)
      } else {
        await api.put(endpoint, payload)
      }

      setShowModal(false)
      
      if (activeTab === 'faculty') setFilterCampusId(form.campus_id)
      if (activeTab === 'department') {
        setFilterCampusId(form.campus_id)
        setFilterFacultyId(form.faculty_id)
      }

      setRefreshKey(prev => prev + 1)
    } catch (err: any) {
      alert(`Error: ${err.response?.data?.message || 'Something went wrong'}`)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Delete "${name}"? This action cannot be undone.`)) return
    try {
      let endpoint = ''
      if (activeTab === 'campus') endpoint = `/api/v1/admin/campuses/${id}`
      else if (activeTab === 'faculty') endpoint = `/api/v1/admin/faculties/${id}`
      else if (activeTab === 'department') endpoint = `/api/v1/admin/departments/${id}`
      
      await api.delete(endpoint)
      setRefreshKey(prev => prev + 1)
    } catch (err: any) {
      alert(`Failed to delete: ${err.response?.data?.message || 'Error'}`)
    }
  }

  //  FILTERING DATA FOR TABLE 

  const filteredCampuses = campuses.filter(c => c.Name?.toLowerCase().includes(searchQuery.toLowerCase()))
  const filteredFaculties = faculties.filter(f => f.Name?.toLowerCase().includes(searchQuery.toLowerCase()))
  const filteredDepartments = departments.filter(d => d.Name?.toLowerCase().includes(searchQuery.toLowerCase()))

  //  RENDER 

  return (
    <div className="users-container">

      {/* ── Page Header ── */}
      <div className="users-header">
        <div className="header-content">
          <h1 className="page-title">จัดการวิทยาเขต คณะ ภาควิชา</h1>
        </div>
        <div className="header-actions">
          <button className="btn btn-primary" onClick={openCreateModal}>
            <span className="btn-icon"><Plus size={18} /></span>
            เพิ่ม{tabs.find(t => t.key === activeTab)?.label}
          </button>
        </div>
      </div>

      {/* ── Tabs ── */}
      <div className="tabs-container">
        {tabs.map((tab) => {
          const TabIcon = tab.icon
          return (
            <button
              key={tab.key}
              className={`tab ${activeTab === tab.key ? 'active' : ''}`}
              onClick={() => handleTabChange(tab.key)}
              style={{ display: 'flex', alignItems: 'center', gap: 8 }}
            >
              <TabIcon size={15} />
              {tab.label}
            </button>
          )
        })}
      </div>

      {/*  Table Controls */}
      <div className="table-controls" style={{ display: 'flex', gap: '16px', flexWrap: 'wrap' }}>
        
        {/* Search */}
        <div className="search-box" style={{ minWidth: '250px' }}>
          <span className="search-icon"><Search size={18} /></span>
          <input
            type="text"
            placeholder={`Search ${activeTab}...`}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
        </div>

        {/* กรองวิทยาเขต  */}
        {(activeTab === 'faculty' || activeTab === 'department') && (
          <div className="search-box" style={{ maxWidth: '250px' }}>
            <span className="search-icon" style={{ color: '#22c55e' }}><Filter size={18} /></span>
            <select
              value={filterCampusId}
              onChange={(e) => {
                setFilterCampusId(e.target.value)
                setFilterFacultyId('') 
              }}
              className="search-input"
              style={{ paddingLeft: '40px', cursor: 'pointer' }}
            >
              <option value="">-- เลือกวิทยาเขต (กรองข้อมูล) --</option>
              {campuses.map(c => <option key={c.ID} value={c.ID}>{c.Name}</option>)}
            </select>
          </div>
        )}

        {/* ตัวกรองคณะ  */}
        {activeTab === 'department' && (
          <div className="search-box" style={{ maxWidth: '250px' }}>
            <span className="search-icon" style={{ color: '#22c55e' }}><Filter size={18} /></span>
            <select
              value={filterFacultyId}
              onChange={(e) => setFilterFacultyId(e.target.value)}
              className="search-input"
              disabled={!filterCampusId} // ปิดการใช้งานถ้ายังไม่ได้เลือกวิทยาเขต
              style={{ paddingLeft: '40px', cursor: filterCampusId ? 'pointer' : 'not-allowed', opacity: !filterCampusId ? 0.6 : 1 }}
            >
              <option value="">-- เลือกคณะ (กรองข้อมูล) --</option>
              {filterFacultiesList.map(f => <option key={f.ID} value={f.ID}>{f.Name}</option>)}
            </select>
          </div>
        )}

      </div>

      {/*  Table  */}
      <div className="table-container">
        <table className="users-table">
          <thead>
            <tr>
              {activeTab === 'campus' && (
                <>
                  <th>ชื่อวิทยาเขต</th>
                  
                </>
              )}
              {activeTab === 'faculty' && (
                <>
                  <th>ชื่อคณะ</th>
                  <th>สังกัดวิทยาเขต</th>
                  
                </>
              )}
              {activeTab === 'department' && (
                <>
                  <th>ชื่อภาควิชา</th>
                  <th>สังกัดคณะ</th>
                  
                </>
              )}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            
            {loading ? (
              <tr>
                <td colSpan={4} className="py-20 text-center">
                  <Loader2 className="w-8 h-8 animate-spin mx-auto" style={{ color: '#22c55e' }} />
                  <p style={{ marginTop: '10px', color: '#64748b' }}>Loading data...</p>
                </td>
              </tr>
            ) : (

              /*  Campus Rows  */
              activeTab === 'campus' ? (
                filteredCampuses.length === 0 ? (
                  <tr><td colSpan={2} className="py-20 text-center" style={{ color: '#64748b' }}>ไม่พบวิทยาเขต</td></tr>
                ) : filteredCampuses.map((campus) => (
                  <tr key={campus.ID}>
                    <td>
                      <div className="user-details">
                        <div className="stat-icon-wrapper" style={{ width: 40, height: 40, background: 'rgba(59,130,246,0.1)', borderRadius: 10, flexShrink: 0 }}>
                          <Building2 size={18} color="#3b82f6" />
                        </div>
                        <div className="user-info">
                          <div className="user-name">{campus.Name}</div>
                          <div className="user-email">ID: {campus.ID}</div>
                        </div>
                      </div>
                    </td>
                    <td style={{ textAlign: 'center' }}>
                     
                    </td>
                  </tr>
                ))

              /* Faculty Rows  */
              ) : activeTab === 'faculty' ? (
                !filterCampusId ? (
                   <tr><td colSpan={3} className="py-20 text-center" style={{ color: '#f59e0b', fontWeight: 600 }}>โปรดเลือกวิทยาเขตจากตัวกรองด้านบนก่อน</td></tr>
                ) : filteredFaculties.length === 0 ? (
                  <tr><td colSpan={3} className="py-20 text-center" style={{ color: '#64748b' }}>ไม่พบคณะในวิทยาเขตนี้</td></tr>
                ) : filteredFaculties.map((faculty) => (
                  <tr key={faculty.ID}>
                    <td>
                      <div className="user-details">
                        <div className="stat-icon-wrapper" style={{ width: 40, height: 40, background: 'rgba(139,92,246,0.1)', borderRadius: 10, flexShrink: 0 }}>
                          <GraduationCap size={18} color="#8b5cf6" />
                        </div>
                        <div className="user-info">
                          <div className="user-name">{faculty.Name}</div>
                          <div className="user-email">ID: {faculty.ID}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <span className="role-badge" style={{ backgroundColor: 'rgba(59,130,246,0.1)', color: '#3b82f6' }}>
                        {campuses.find(c => c.ID === filterCampusId)?.Name || '-'}
                      </span>
                    </td>
                    <td style={{ textAlign: 'center' }}>
                    
                    </td>
                  </tr>
                ))

              /*  Department Rows  */
              ) : (
                !filterFacultyId ? (
                   <tr><td colSpan={3} className="py-20 text-center" style={{ color: '#f59e0b', fontWeight: 600 }}>โปรดเลือกวิทยาเขตและคณะจากตัวกรองด้านบนก่อน</td></tr>
                ) : filteredDepartments.length === 0 ? (
                  <tr><td colSpan={4} className="py-20 text-center" style={{ color: '#64748b' }}>ไม่พบภาควิชาในคณะนี้</td></tr>
                ) : filteredDepartments.map((dept) => (
                  <tr key={dept.ID}>
                    <td>
                      <div className="user-details">
                        <div className="stat-icon-wrapper" style={{ width: 40, height: 40, background: 'rgba(245,158,11,0.1)', borderRadius: 10, flexShrink: 0 }}>
                          <BookOpen size={18} color="#f59e0b" />
                        </div>
                        <div className="user-info">
                          <div className="user-name">{dept.Name}</div>
                          <div className="user-email">ID: {dept.ID}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <span className="role-badge" style={{ backgroundColor: 'rgba(139,92,246,0.1)', color: '#8b5cf6' }}>
                        {filterFacultiesList.find(f => f.ID === filterFacultyId)?.Name || '-'}
                      </span>
                    </td>
                    <td style={{ textAlign: 'center' }}>
                    
                    </td>
                  </tr>
                ))
              )
            )}
          </tbody>
        </table>
      </div>

      {/*  FORM MODAL  */}
      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2 className="modal-title">
                {modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'} {tabs.find(t => t.key === activeTab)?.label}
              </h2>
              <button className="modal-close" onClick={() => setShowModal(false)}>
                <X size={20} />
              </button>
            </div>

            <div className="profile-content">
              <div className="flex flex-col gap-5">

                {/* Campus Dropdown for Modal */}
                {(activeTab === 'faculty' || activeTab === 'department') && (
                  <div>
                    <label className="form-label">
                      วิทยาเขต <span style={{ color: '#ef4444' }}>*</span>
                    </label>
                    <select
                      value={form.campus_id}
                      onChange={(e) => setForm({ ...form, campus_id: e.target.value, faculty_id: '' })}
                      className="form-select"
                      style={{ width: '100%', cursor: 'pointer' }}
                    >
                      <option value="">-- Select Campus --</option>
                      {campuses.map(c => (
                        <option key={c.ID} value={c.ID}>{c.Name}</option>
                      ))}
                    </select>
                  </div>
                )}

                {/* Faculty Dropdown for Modal */}
                {activeTab === 'department' && (
                  <div>
                    <label className="form-label">
                      คณะ <span style={{ color: '#ef4444' }}>*</span>
                    </label>
                    <select
                      value={form.faculty_id}
                      onChange={(e) => setForm({ ...form, faculty_id: e.target.value })}
                      disabled={!form.campus_id}
                      className="form-select"
                      style={{ width: '100%', opacity: !form.campus_id ? 0.6 : 1, cursor: form.campus_id ? 'pointer' : 'not-allowed' }}
                    >
                      <option value="">
                        {form.campus_id ? '-- Select Faculty --' : 'Select Campus First'}
                      </option>
                      {formFaculties.map(f => (
                        <option key={f.ID} value={f.ID}>{f.Name}</option>
                      ))}
                    </select>
                  </div>
                )}

                {/* Name Input */}
                <div>
                  <label className="form-label">
                    {tabs.find(t => t.key === activeTab)?.label}  <span style={{ color: '#ef4444' }}>*</span>
                  </label>
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    placeholder={`Enter name...`}
                    className="form-input"
                    style={{ width: '100%' }}
                  />
                </div>

                {/* Actions */}
                <div className="form-actions">
                  <button className="btn btn-secondary" onClick={() => setShowModal(false)}>
                    Cancel
                  </button>
                  <button className="btn btn-primary" onClick={handleSubmit} disabled={isSubmitting}
                    style={{ opacity: isSubmitting ? 0.6 : 1 }}>
                    {isSubmitting ? <Loader2 size={16} className="animate-spin" /> : <Check size={16} />}
                    {isSubmitting ? 'Saving...' : modalMode === 'create' ? 'Create' : 'Save Changes'}
                  </button>
                </div>

              </div>
            </div>
          </div>
        </div>
      )}

    </div>
  )
}