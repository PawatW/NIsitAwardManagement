"use client";

import Link from 'next/link'
import { useAuth } from '../../../context/AuthContext'; 
import { useEffect, useState } from 'react';
import api from '../../../lib/axios'; 

interface SidebarProps {
  activePage: 'review' | 'announcements';
}

export default function Sidebar({ activePage }: SidebarProps) {
  const { user, logout } = useAuth();
  

  const [profileUrl, setProfileUrl] = useState<string | null>(null);


  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!user) return;
      try {
        const res = await api.get('/api/v1/user/me');
        const userData = res.data.data || res.data;
        
        if (userData.profile_url) {
          setProfileUrl(userData.profile_url);
        }
      } catch (error) {
        console.error("Failed to fetch user profile:", error);
      }
    };

    fetchUserProfile();
  }, [user]);

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <div className="logo">
          <div className="logo-icon">ST</div>
          <div className="logo-text">
            <div className="logo-title">Nisit D-Den</div>
            <div className="logo-subtitle">ระบบพิจารณารางวัลนิสิตดีเด่น</div>
          </div>
        </div>
      </div>

      <nav className="sidebar-nav">
        <Link href="/staff-page/app/review" className={`nav-item ${activePage === 'review' ? 'active' : ''}`}>
          <span className="nav-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path d="M4 10L8 14L16 6" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </span>
          <span className="nav-label">ตรวจสอบใบสมัคร</span>
        </Link>

        <Link href="/announcements" className={`nav-item ${activePage === 'announcements' ? 'active' : ''}`}>
          <span className="nav-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
              <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6" />
              <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18" />
              <path d="M4 22h16" />
              <path d="M10 14.66V17c0 .55-.47 .98-.97 1.21C7.85 18.75 7 20.24 7 22" />
              <path d="M14 14.66V17c0 .55.47 .98.97 1.21C16.15 18.75 17 20.24 17 22" />
              <path d="M18 2H6v7a6 6 0 0 0 12 0V2Z" />
            </svg>
          </span>
          <span className="nav-label">ทำเนียบนิสิตดีเด่น</span>
        </Link>
      </nav> 
        
      <div className="sidebar-footer">
        <div className="user-profile">
          <div className="user-avatar overflow-hidden flex items-center justify-center bg-green-500 text-white font-bold">
            {profileUrl ? (
              <img 
                src={profileUrl} 
                alt="Profile" 
                className="w-full h-full object-cover"
                onError={(e) => {
                  e.currentTarget.style.display = 'none';
                }}
              />
            ) : (

              <span>{user?.email ? user.email.charAt(0).toUpperCase() : 'S'}</span>
            )}
          </div>
          
          <div className="user-info">
            <div className="user-name" title={user?.email}>{user?.email || 'กำลังโหลด...'}</div>
            <div className="user-id">Role: Staff </div>
          </div>
        </div>
        
        <button className="logout-button" onClick={() => {
          if (window.confirm('ต้องการออกจากระบบหรือไม่?')) {
            logout();
          }
        }}>
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M13 14L17 10M17 10L13 6M17 10H7M7 17H4C3.46957 17 2.96086 16.7893 2.58579 16.4142C2.21071 16.0391 2 15.5304 2 15V5C2 4.46957 2.21071 3.96086 2.58579 3.58579C2.96086 3.21071 3.46957 3 4 3H7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
          ออกจากระบบ
        </button>
      </div>
    </aside>
  );
}