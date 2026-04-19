"use client";

import { useState, useEffect } from "react";
import { GraduationCap, LogOut, Trophy, User } from "lucide-react"; 
import Link from "next/link"; 
import { useAuth } from "../../context/AuthContext";
import api from "../../lib/axios"; 

const roleTranslationMap: Record<string, string> = {
  "ADMIN": "ผู้ดูแลระบบ",
  "STUDENT": "นิสิต",
  "COMMITTEE_CHAIR": "คณะกรรมการ",
  "DEAN": "คณบดี",
  "VICE_DEAN": "รองคณบดี",
  "HEAD_OF_DEPARTMENT": "หัวหน้าภาควิชา"
};

export default function Navbar() {
  const { user, logout } = useAuth();
  
  const [profileUrl, setProfileUrl] = useState<string | null>(null);
  const [profileName, setProfileName] = useState<string | null>(null);

  const emailName = user?.email ? user.email.split('@')[0] : ""; 
  const role = user?.role || "STUDENT"; 

  const getThaiRole = (englishRole: string) => {
    return roleTranslationMap[englishRole.toUpperCase()] || englishRole.replace(/_/g, ' ');
  };

  useEffect(() => {
    const fetchUserProfile = async () => {

      if (!user) return;
      try {
        const res = await api.get('/api/v1/user/me');
        const userData = res.data.data || res.data;
        
        if (userData.profile_url) {
          setProfileUrl(userData.profile_url);
        }

        if (userData.first_name && userData.last_name) {
          setProfileName(`${userData.first_name} ${userData.last_name}`);
        } else if (userData.name) {
          setProfileName(userData.name);
        }
      } catch (error) {
        console.error("Failed to fetch user profile:", error);
      }
    };

    fetchUserProfile();
  }, [user]);

  return (
    <nav className="w-full bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center sticky top-0 z-50 shadow-sm">
      {/* Logo & Title */}
      <Link href="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
        <div className="text-green-600">
          <GraduationCap size={40} strokeWidth={1.5} />
        </div>
        <div className="flex flex-col">
          <h1 className="text-xl font-bold text-gray-900 leading-tight">
            ระบบสมัครนิสิตดีเด่น 
          </h1>
          <span className="text-sm text-green-600 font-medium">
            มหาวิทยาลัยเกษตรศาสตร์
          </span>
        </div>
      </Link>

      {/* Right Section */}
      <div className="flex items-center gap-4 sm:gap-6">
        <Link 
          href="/announcements" 
          className="flex items-center gap-2 px-3 py-2 text-gray-600 hover:text-yellow-600 hover:bg-yellow-50 rounded-xl transition-all duration-200 font-medium"
          title="ดูทำเนียบนิสิตดีเด่น"
        >
          <Trophy size={20} />
          <span className="hidden sm:inline">ทำเนียบรางวัล</span>
        </Link>


        {user ? (
          <div className="flex items-center gap-3 border-l border-gray-200 pl-4 sm:pl-6">
            <div className="hidden sm:flex flex-col items-end">
              <span className="text-sm font-bold text-gray-900 line-clamp-1 max-w-[150px]">
                {profileName ? profileName : emailName} 
              </span>
              <span className="text-[12px] bg-green-50 text-green-700 px-2 py-0.5 rounded-full font-bold tracking-wider">
   
                {getThaiRole(role)}
              </span>
            </div>
            
            {/* Avatar */}
            <div className="w-10 h-10 rounded-full bg-green-100 border border-green-200 flex items-center justify-center overflow-hidden shrink-0 text-green-700 font-bold shadow-inner">
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
                 <span className="text-lg uppercase">
                   {(profileName || emailName || "?").charAt(0)}
                 </span>
               )}
            </div>

            {/* Logout Button */}
            <button 
              onClick={logout}
              className="ml-1 p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-full transition-all duration-200"
              title="ออกจากระบบ"
            >
              <LogOut size={20} />
            </button>
          </div>
        ) : (

          <Link 
            href="/"
            className="flex items-center gap-2 px-5 py-2 bg-green-600 text-white rounded-xl hover:bg-green-700 transition-all font-bold text-sm shadow-md shadow-green-100"
          >
            <User size={18} />
            เข้าสู่ระบบ
          </Link>
        )}
      </div>
    </nav>
  );
}