'use client'

import { useState, useEffect } from 'react'
import { 
  Trophy, GraduationCap, Plus, Check, X, FileText, AlignLeft, Hash, List, 
  Edit, Trash2, Loader2
} from 'lucide-react'
import api from '../../../../lib/axios'
import '../../styles/awards.css'

interface FormField {
  id: string          
  field_key: string  
  label: string
  type: string
  required: boolean
  placeholder?: string
  options?: string[]
}

interface AwardCategory {
  id: number
  name: string
  description: string
  form_structure: any[]
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export default function Awards() {
  const [awards, setAwards] = useState<AwardCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [refreshKey, setRefreshKey] = useState(0)

  const [showModal, setShowModal] = useState<boolean>(false)
  const [isEditing, setIsEditing] = useState<boolean>(false)
  const [editingId, setEditingId] = useState<number | null>(null)

  const [formData, setFormData] = useState({ name: '', description: '', is_active: true })
  const [formFields, setFormFields] = useState<FormField[]>([])

  const fieldTypes = [
    { value: 'text',     label: 'Short Text',      IconComponent: FileText  },
    { value: 'textarea', label: 'Long Text',         IconComponent: AlignLeft },
    { value: 'number',   label: 'Number',            IconComponent: Hash      },
    { value: 'select',   label: 'Dropdown Select',   IconComponent: List      },
  ]

  // ── Fetch ──
  useEffect(() => {
    const fetchAwards = async () => {
      try {
        setLoading(true)
        const res = await api.get('/api/v1/award-categories')
        const rawData = res.data.data || res.data || []
        const mappedData = rawData.map((item: any) => ({
          ...item,
          form_structure: typeof item.form_structure === 'string'
            ? JSON.parse(item.form_structure)
            : (item.form_structure || [])
        }))
        setAwards(mappedData)
      } catch (error) {
        console.error('Failed to fetch awards:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchAwards()
  }, [refreshKey])

  // ── Form handlers ──
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target
    const checked = (e.target as HTMLInputElement).checked
    setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }))
  }

  const addField = () => {
    setFormFields([...formFields, {
      id: `react_key_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      field_key: '', 
      label: '', 
      type: 'text', 
      required: true, 
      placeholder: '', 
      options: []
    }])
  }

  const removeField = (id: string) => setFormFields(formFields.filter(f => f.id !== id))

  const updateField = (id: string, key: keyof FormField, value: any) => {
    setFormFields(formFields.map(f => f.id === id ? { ...f, [key]: value } : f))
  }

  const addOption = (fieldId: string) => {
    setFormFields(formFields.map(f => {
      if (f.id !== fieldId) return f
      const opts = f.options || []
      return { ...f, options: [...opts, `Option ${opts.length + 1}`] }
    }))
  }

  const updateOption = (fieldId: string, idx: number, val: string) => {
    setFormFields(formFields.map(f => {
      if (f.id !== fieldId || !f.options) return f
      const opts = [...f.options]; opts[idx] = val
      return { ...f, options: opts }
    }))
  }

  const removeOption = (fieldId: string, idx: number) => {
    setFormFields(formFields.map(f => {
      if (f.id !== fieldId || !f.options) return f
      return { ...f, options: f.options.filter((_, i) => i !== idx) }
    }))
  }

  // ── Submit ──
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    
    const invalidField = formFields.find(f => f.label.trim() === '' || f.field_key.trim() === '')
    if (invalidField) {
        return alert("กรุณากรอก Field ID และ ชื่อคำถาม ให้ครบทุกข้อ")
    }

    const fieldKeys = formFields.map(f => f.field_key.trim())
    const uniqueKeys = new Set(fieldKeys)
    if (uniqueKeys.size !== fieldKeys.length) {
        return alert("มี Field ID ซ้ำกัน กรุณาตั้งชื่อ Field ID ให้ไม่ซ้ำกัน")
    }

    setIsSubmitting(true)

    // นำ field_key ที่ Admin พิมพ์ มาใช้เป็น id ส่งให้ Backend
    const finalStructure = formFields.map(f => {
      const fieldObj: any = { 
          id: f.field_key.trim(), // 🔥 ส่งตัวนี้ให้ Backend
          label: f.label.trim(), 
          type: f.type, 
          required: true 
      }
      if (f.type === 'select' && f.options?.length) {
        fieldObj.options = f.options.filter(o => o.trim() !== '')
      }
      return fieldObj
    })

    const payload = {
      name: formData.name,
      description: formData.description,
      is_active: formData.is_active,
      form_structure: finalStructure
    }

    try {
      if (isEditing && editingId) {
        await api.put(`/api/v1/award-categories/${editingId}`, payload)
        alert('อัปเดตประเภทรางวัลสำเร็จ!')
      } else {
        await api.post('/api/v1/award-categories', payload)
        alert('สร้างประเภทรางวัลสำเร็จ!')
      }
      resetForm()
      setShowModal(false)
      setRefreshKey(prev => prev + 1)
    } catch (error: any) {
      alert(`Error: ${error.response?.data?.message || 'Operation failed.'}`)
    } finally {
      setIsSubmitting(false)
    }
  }


  const openAddModal = () => { resetForm(); setIsEditing(false); setShowModal(true) }

  const openEditModal = (award: AwardCategory) => {
    setFormData({ name: award.name, description: award.description, is_active: award.is_active ?? true })
    const structure = Array.isArray(award.form_structure)
      ? award.form_structure
      : (award.form_structure as any)?.fields || []
    
    setFormFields(structure.map((f: any) => ({
      ...f,
      id: `react_key_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      field_key: f.id, // 🔥 เอา id จาก Backend มาใส่ใน field_key ให้ Admin เห็นและแก้ได้
      options: f.options || []
    })))
    setEditingId(award.id)
    setIsEditing(true)
    setShowModal(true)
  }

  const resetForm = () => {
    setFormData({ name: '', description: '', is_active: true })
    setFormFields([])
    setEditingId(null)
  }

  const statsData = [
    { label: 'ปริมาณประเภทรางวัลทั้งหมด',  value: awards.length.toString(),                          IconComponent: Trophy,        iconBg: 'rgba(59,130,246,0.15)',  iconColor: '#3b82f6' },
    { label: 'ประเภทรางวัลที่เปิดใช้งาน', value: awards.filter(a => a.is_active).length.toString(), IconComponent: GraduationCap, iconBg: 'rgba(16,185,129,0.15)',  iconColor: '#10b981' },
  ]

  return (
    <div className="awards-container">

      {/* ── Page Header ── */}
      <div className="awards-header">
        <div className="header-content">
          <h1 className="page-title">จัดการประเภทรางวัล</h1>
        </div>
        <button className="btn btn-primary" onClick={openAddModal}>
          <span className="btn-icon"><Plus size={18} /></span>
          เพิ่มประเภทรางวัล
        </button>
      </div>

      {/* ── Stats ── */}
      <div className="stats-grid" style={{ gridTemplateColumns: 'repeat(auto-fit,minmax(250px,1fr))', marginBottom: 30 }}>
        {statsData.map((stat, i) => {
          const Icon = stat.IconComponent
          return (
            <div key={i} className="stat-card" style={{ animationDelay: `${i * 0.1}s` }}>
              <div className="stat-icon-wrapper" style={{ background: stat.iconBg }}>
                <span className="stat-icon"><Icon size={24} color={stat.iconColor} strokeWidth={2} /></span>
              </div>
              <div className="stat-content">
                <div className="stat-label">{stat.label}</div>
                <div className="stat-value">{stat.value}</div>
              </div>
            </div>
          )
        })}
      </div>

    {/* ── Awards Table ── */}
      <div className="awards-table-container">
        <div className="table-header">
          <div className="table-title-section">
            <h2 className="table-title">ชื่อและรายละเอียดรางวัล</h2>
          </div>
          <div className="table-columns">
            <div className="table-column">สถานะ</div>
            <div className="table-column text-right">Actions</div>
          </div>
        </div>

        <div className="awards-list">
          {loading ? (
            <div className="flex justify-center items-center py-20">
              <Loader2 className="animate-spin text-green-600" size={32} />
            </div>
          ) : awards.length === 0 ? (
            <div className="text-center py-20 text-gray-500">No award categories found.</div>
          ) : awards.map((award, index) => (
            <div key={award.id} className="award-item" style={{ animationDelay: `${index * 0.05}s` }}>
              <div className="award-main">
                <div className="award-icon-wrapper" style={{ background: 'rgba(59,130,246,0.15)' }}>
                  <span className="award-icon"><Trophy size={24} color="#3b82f6" strokeWidth={2} /></span>
                </div>
                <div className="award-info">
                  <h3 className="award-name">{award.name}</h3>
                  <p className="award-description text-sm text-gray-500">{award.description}</p>
                </div>
              </div>
              <div className="award-details">
                <div className="award-detail-item">
                  <span className={`px-3 py-1 rounded-full text-xs font-bold ${award.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                    {award.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
                
   
                <div className="award-detail-item flex justify-end">
                  <div className="action-buttons flex gap-2"> 
                    <button className="action-btn edit-btn" title="Edit" onClick={() => openEditModal(award)}><Edit size={18} /></button>
                  
                  </div>
                </div>
                
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* MODAL  */}
      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal-content modal-large" onClick={e => e.stopPropagation()}>

            <div className="modal-header">
              <h2 className="modal-title">
                {isEditing ? 'แก้ไขประเภทรางวัล' : 'เพิ่มประเภทรางวัลใหม่'}
              </h2>
              <button className="modal-close" onClick={() => setShowModal(false)}>
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit}>
              <div className="profile-content">

                {/* ── Basic Info ── */}
                <div className="profile-details">
                  <h4 className="section-title">ข้อมูลพื้นฐาน</h4>

                  <div className="form-group">
                    <label className="form-label">
                      ชื่อประเภทรางวัล <span style={{ color: '#ef4444' }}>*</span>
                    </label>
                    <input
                      type="text" name="name" value={formData.name}
                      onChange={handleInputChange} required
                      placeholder="เช่น รางวัลนิสิตดีเด่นด้านวิชาการ"
                      className="form-input"
                    />
                  </div>

                  <div className="form-group">
                    <label className="form-label">
                      คำอธิบาย <span style={{ color: '#ef4444' }}>*</span>
                    </label>
                    <textarea
                      name="description" value={formData.description}
                      onChange={handleInputChange} required rows={3}
                      placeholder="อธิบายเกี่ยวกับเกณฑ์การสมัคร..."
                      className="form-textarea"
                    />
                  </div>

                  {/* Active toggle */}
                  <div style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                    padding: '16px 20px', background: '#f8fafc',
                    borderRadius: 12, border: '1px solid #e2e8f0'
                  }}>
                    <div>
                      <p style={{ fontSize: 14, fontWeight: 600, color: '#0f172a', margin: 0 }}>สถานะการใช้งาน</p>
                      <p style={{ fontSize: 12, color: '#64748b', margin: '4px 0 0' }}>ผู้สมัครจะเห็นได้ทันที</p>
                    </div>
                    <button
                      type="button"
                      onClick={() => setFormData(p => ({ ...p, is_active: !p.is_active }))}
                      style={{
                        position: 'relative', display: 'inline-flex', height: 24, width: 44,
                        alignItems: 'center', borderRadius: 9999, border: 'none', cursor: 'pointer',
                        backgroundColor: formData.is_active ? '#22c55e' : '#cbd5e1',
                        transition: 'background-color 0.2s', flexShrink: 0
                      }}
                    >
                      <span style={{
                        display: 'inline-block', width: 16, height: 16, borderRadius: '50%',
                        backgroundColor: 'white',
                        transform: formData.is_active ? 'translateX(24px)' : 'translateX(4px)',
                        transition: 'transform 0.2s'
                      }} />
                    </button>
                  </div>
                </div>

                {/* ── Custom Form Fields ── */}
                <div className="profile-details">
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16, paddingBottom: 12, borderBottom: '2px solid #e2e8f0' }}>
                    <h4 style={{ fontSize: 16, fontWeight: 700, color: '#0f172a', margin: 0 }}>ฟิลด์แบบกำหนดเอง (Dynamic Form)</h4>
                    <button
                      type="button" onClick={addField}
                      className="btn btn-secondary"
                      style={{ padding: '8px 16px', fontSize: 13 }}
                    >
                      <Plus size={15} /> เพิ่มฟิลด์
                    </button>
                  </div>

                  {formFields.length === 0 ? (
                    <div style={{
                      textAlign: 'center', padding: '40px 24px',
                      border: '2px dashed #e2e8f0', borderRadius: 12, color: '#94a3b8'
                    }}>
                      ยังไม่มีฟิลด์แบบกำหนดเอง<br />
                    </div>
                  ) : (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
                      {formFields.map((field, index) => (
                        <div key={field.id} style={{
                          padding: 20, border: '2px solid #e2e8f0', borderRadius: 16,
                          background: '#f8fafc', display: 'flex', flexDirection: 'column', gap: 16
                        }}>
                          {/* Field header */}
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingBottom: 12, borderBottom: '1px solid #e2e8f0' }}>
                            <span style={{ fontSize: 11, fontWeight: 700, color: '#64748b', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
                              Field {index + 1}
                            </span>
                            <button type="button" onClick={() => removeField(field.id)}
                              style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#94a3b8', display: 'flex', alignItems: 'center' }}
                              onMouseEnter={e => (e.currentTarget.style.color = '#ef4444')}
                              onMouseLeave={e => (e.currentTarget.style.color = '#94a3b8')}
                            >
                              <Trash2 size={16} />
                            </button>
                          </div>

           
                          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: 0 }}>
                            
                            {/* Field ID (ภาษาอังกฤษ) */}
                            <div className="form-group" style={{ marginBottom: 0 }}>
                              <label className="form-label">
                                Field ID (ภาษาอังกฤษตัวพิมพ์เล็ก, ตัวเลข, _) <span style={{ color: '#ef4444' }}></span>
                              </label>
                              <input
                                type="text" 
                                value={field.field_key} 
                                required
                                onChange={e => {
                                  // 🔥 อนุญาตแค่ a-z, 0-9 และ _
                                  const val = e.target.value.replace(/[^a-z0-9_]/g, '').toLowerCase()
                                  updateField(field.id, 'field_key', val)
                                }}
                                placeholder="เช่น student_gpa"
                                className="form-input"
                              />
                            </div>

                            {/* Label (ชื่อที่โชว์ให้ผู้ใช้เห็น) */}
                            <div className="form-group" style={{ marginBottom: 0 }}>
                              <label className="form-label">
                                ชื่อคำถาม  <span style={{ color: '#ef4444' }}>*</span>
                              </label>
                              <input
                                type="text" 
                                value={field.label} 
                                required
                                onChange={e => updateField(field.id, 'label', e.target.value)}
                                placeholder="เช่น เกรดเฉลี่ยสะสม"
                                className="form-input"
                              />
                            </div>
                            
                          </div>

                          {/* Type */}
                          <div className="form-group" style={{ marginBottom: 0, maxWidth: '50%' }}>
                            <label className="form-label">ประเภท </label>
                            <select
                              value={field.type}
                              onChange={e => {
                                updateField(field.id, 'type', e.target.value)
                                if (e.target.value === 'select' && (!field.options || field.options.length === 0)) {
                                  updateField(field.id, 'options', ['Option 1'])
                                }
                              }}
                              className="form-select"
                            >
                              {fieldTypes.map(t => (
                                <option key={t.value} value={t.value}>{t.label}</option>
                              ))}
                            </select>
                          </div>

                          {/* Options (select type) */}
                          {field.type === 'select' && (
                            <div style={{
                              padding: '16px', background: 'rgba(59,130,246,0.05)',
                              borderRadius: 12, border: '1px solid rgba(59,130,246,0.15)'
                            }}>
                              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
                                <span style={{ fontSize: 12, fontWeight: 700, color: '#3b82f6', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
                                  ตัวเลือก (Choices)
                                </span>
                                <button type="button" onClick={() => addOption(field.id)}
                                  className="btn btn-secondary"
                                  style={{ padding: '6px 12px', fontSize: 12 }}
                                >
                                  <Plus size={13} /> เพิ่มตัวเลือก
                                </button>
                              </div>
                              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                                {field.options?.map((opt, optIdx) => (
                                  <div key={optIdx} style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
                                    <div style={{
                                      width: 28, height: 28, borderRadius: '50%',
                                      background: '#dbeafe', color: '#3b82f6',
                                      display: 'flex', alignItems: 'center', justifyContent: 'center',
                                      fontSize: 12, fontWeight: 700, flexShrink: 0
                                    }}>
                                      {optIdx + 1}
                                    </div>
                                    <input
                                      type="text" value={opt} required
                                      onChange={e => updateOption(field.id, optIdx, e.target.value)}
                                      placeholder={`ตัวเลือกที่ ${optIdx + 1}`}
                                      className="form-input"
                                      style={{ flex: 1, marginBottom: 0 }}
                                    />
                                    <button type="button" onClick={() => removeOption(field.id, optIdx)}
                                      style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#94a3b8', display: 'flex', alignItems: 'center' }}
                                      onMouseEnter={e => (e.currentTarget.style.color = '#ef4444')}
                                      onMouseLeave={e => (e.currentTarget.style.color = '#94a3b8')}
                                    >
                                      <X size={16} />
                                    </button>
                                  </div>
                                ))}
                                {(!field.options || field.options.length === 0) && (
                                  <p style={{ fontSize: 12, color: '#ef4444', margin: 0 }}>กรุณาเพิ่มตัวเลือกอย่างน้อย 1 ข้อ</p>
                                )}
                              </div>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                {/* ── Form Actions ── */}
                <div className="form-actions">
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>
                    ยกเลิก
                  </button>
                  <button type="submit" className="btn btn-primary" disabled={isSubmitting}
                    style={{ opacity: isSubmitting ? 0.6 : 1 }}>
                    {isSubmitting ? <Loader2 size={16} className="animate-spin" /> : <Check size={16} />}
                    {isEditing ? 'อัปเดตประเภทรางวัล' : 'บันทึกประเภทรางวัล'}
                  </button>
                </div>

              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}