import { LucideIcon } from 'lucide-react'

// Common Types
export interface IconData {
  IconComponent: LucideIcon
  label?: string
}

// Dashboard Types
export interface AwardPeriod {
  id: number
  campus: string
  semester: string
  openingDate: string
  closingDate: string
  daysRemaining: number
  isActive: boolean
  awardCategories: AwardCategory[]
}

export interface AwardCategory {
  id: number
  name: string
  isActive: boolean
}

export interface StatData {
  label: string
  value: string
  IconComponent: LucideIcon
  iconBg: string
}

// User Types
export interface User {
  id: number
  name: string
  email: string
  role: 'Student' | 'Faculty' | 'Staff' | 'Committee'
  faculty: string
  department: string
  campus: string
  status: 'Active' | 'Inactive'
  lastLogin?: string
  createdDate?: string
}

export interface UserFormData {
  name: string
  email: string
  role: string
  faculty: string
  department: string
  campus: string
  status: string
}

export interface ActivityHistory {
  id: number
  type: 'submission' | 'login' | 'status_change' | 'account_created'
  title: string
  description: string
  timestamp: string
  timeAgo: string
}

// Awards Types
export interface Award {
  id: number
  name: string
  description: string
  IconComponent: LucideIcon
  iconBg: string
  benefitType: string
  benefitColor: string
  requiredDocs: string[]
  updated: string
}

export interface FormSection {
  id: number
  title: string
  isLocked?: boolean
  fields: FormField[]
}

export interface FormField {
  id: string
  label: string
  fieldType: 'text' | 'textarea' | 'number' | 'dropdown' | 'date' | 'email' | 'phone'
  required: boolean
  placeholder?: string
  disabled?: boolean
  options?: string[]
}

export interface FieldType {
  value: string
  label: string
  IconComponent: LucideIcon
}

export interface DocumentUploadField {
  id: string
  label: string
  required: boolean
  description?: string
}

export interface DocumentUpload {
  enabled: boolean
  fields: DocumentUploadField[]
}

// Settings Types
export interface ProfileData {
  name: string
  email: string
  role: string
  avatar: File | null
  avatarPreview: string | null
}

export interface PasswordData {
  currentPassword: string
  newPassword: string
  confirmPassword: string
}

// Register Types
export interface RegisterFormData {
  firstName: string
  lastName: string
  studentId: string
  email: string
  phone: string
  facultyId: string
  majorId: string
  password: string
  confirmPassword: string
}

export interface Faculty {
  id: number
  name: string
}

export interface Major {
  id: number
  name: string
}

export interface ValidationErrors {
  [key: string]: string
}
