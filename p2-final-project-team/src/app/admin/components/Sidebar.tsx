'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation' 
import { 
  LayoutDashboard, 
  Users, 
  Trophy, 
  GraduationCap,
  ChevronLeft,
  ChevronRight,
  LogOut,
  Building2
} from 'lucide-react'

import { useAuth } from '../../../context/AuthContext' 

interface MenuItem {
  path: string
  icon: any
  label: string
}

export default function Sidebar() {
  const [collapsed, setCollapsed] = useState<boolean>(false)
  const pathname = usePathname()
  const router = useRouter() 
  
  
  const { user, logout } = useAuth() 

  const menuItems: MenuItem[] = [
    { path: '/admin/app', icon: LayoutDashboard, label: 'Dashboard' },
    { path: '/admin/app/users', icon: Users, label: 'ผู้ใช้งาน' },
    { path: '/admin/app/awards', icon: Trophy, label: 'รางวัล' },
    { path: '/admin/app/organization', icon: Building2, label: 'องค์กร' },
    {path: '/announcements', icon: GraduationCap, label: 'ทำเนียบนิสิตดีเด่น' },
  ]


  const handleLogout = async () => {
    if (window.confirm('คุณต้องการออกจากระบบใช่หรือไม่?')) {
      try {

        if (logout) {
          await logout()
        } else {

          localStorage.removeItem('token') 
          sessionStorage.removeItem('token')
        }
        
        router.push('/') 
      } catch (error) {
        console.error('Logout failed', error)
      }
    }
  }



  return (
    <aside className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
      <div className="sidebar-header">
        <div className="logo">
          <div className="logo-icon">
            <GraduationCap size={28} strokeWidth={2.5} />
          </div>
          <div className="logo-text">
            <div className="logo-title">แอดมิน</div>
            <div className="logo-subtitle">ระบบจัดการนิสิตดีเด่น</div>
          </div>
        </div>
        <button 
          className="toggle-btn" 
          onClick={() => setCollapsed(!collapsed)}
          title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {collapsed ? <ChevronRight size={16} /> : <ChevronLeft size={16} />}
        </button>
      </div>

      <nav className="sidebar-nav">
        {menuItems.map((item) => {
          const Icon = item.icon
          if (item.path === '/announcements') {
             return (
              <a
                key={item.path}
                href={item.path}
                className={`nav-item ${pathname === item.path ? 'active' : ''}`}
                title={collapsed ? item.label : ''}
              >
                <span className="nav-icon">
                  <Icon size={20} strokeWidth={2} />
                </span>
                <span className="nav-label">{item.label}</span>
              </a>
            )
          }
          return (
            <Link
              key={item.path}
              href={item.path}
              className={`nav-item ${pathname === item.path ? 'active' : ''}`}
              title={collapsed ? item.label : ''}
            >
              <span className="nav-icon">
                <Icon size={20} strokeWidth={2} />
              </span>
              <span className="nav-label">{item.label}</span>
            </Link>
          )
        })}
      </nav>

      <div className="sidebar-footer">
        <div className="user-profile">

          <div className="user-info">
         
            <div className="user-name">{'Admin'}</div>
            <div className="user-role uppercase text-[10px]">{user?.role || 'System Admin'}</div>
          </div>
        </div>
        
   
        <button 
          className="logout-btn" 
          title="Logout"
          onClick={handleLogout}
        >
          <span className="logout-icon">
            <LogOut size={16} strokeWidth={2} />
          </span>
          <span className="logout-text">ออกจากระบบ</span>
        </button>
      </div>
    </aside>
  )
}