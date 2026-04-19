"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { 
  LayoutDashboard, 
  FileText, 
  User, 
  LogOut, 
  Settings, 
  CheckSquare, 
  Users 
} from "lucide-react";
import { useAuth } from "../../context/AuthContext"; 

import api from "../../lib/axios"; 

export default function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  const [profileUrl, setProfileUrl] = useState<string | null>(null);

  const [profileName, setProfileName] = useState<string | null>(null);

  const commonMenu = [
    { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  ];

  const studentMenu = [
    ...commonMenu,
    { name: "คำร้องของฉัน", href: "/my-applications", icon: FileText },
    { name: "ผลการพิจารณานิสิดีเด่น", href: "/announcements", icon: CheckSquare },
  ];
  
  const menuItems = user?.role === "STUDENT" ? studentMenu : commonMenu;

  const emailName = user?.email ? user.email.split('@')[0] : "Guest";
  const displayName = profileName || emailName;

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

  const getDisplayRole = (role: string | undefined) => {
    if (!role) return "Guest";
    if (role === "STUDENT") return "นิสิต";
    return role.toLowerCase();
  };

  return (
    <aside className="h-screen w-64 bg-white border-r border-gray-200 flex flex-col p-6 fixed left-0 top-0 font-sans z-40">
      
      {/* User Profile Section */}
      <div className="flex items-center gap-3 mb-10">
        <div className="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center flex-shrink-0 text-green-700 font-bold text-lg border border-green-200 overflow-hidden">
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
            displayName.charAt(0).toUpperCase()
          )}
        </div>
        
        <div className="flex flex-col overflow-hidden">
          <span className="font-bold text-gray-800 text-sm truncate" title={displayName}>
            {displayName}
          </span>
          <span className="text-xs text-gray-500 capitalize">
            {getDisplayRole(user?.role)}
          </span>
        </div>
      </div>

      {/* Navigation Menu */}
      <nav className="flex-1 space-y-2">
        {menuItems.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
          
          if (item.href === '/announcements') {
             return (
              <a
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium ${
                  isActive
                    ? "bg-green-100 text-green-800" 
                    : "text-gray-600 hover:bg-gray-50 hover:text-gray-900" 
                }`}
              >
                <item.icon size={20} />
                {item.name}
              </a>
            );
          }

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium ${
                isActive
                  ? "bg-green-100 text-green-800" 
                  : "text-gray-600 hover:bg-gray-50 hover:text-gray-900" 
              }`}
            >
              <item.icon size={20} />
              {item.name}
            </Link>
          );
        })}
      </nav>

      {/* Logout Button */}
      <div className="mt-auto pt-6 border-t border-gray-100">
        <button
          onClick={logout} 
          className="flex items-center gap-3 px-4 py-3 w-full text-gray-600 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors text-sm font-medium"
        >
          <LogOut size={20} />
          ออกจากระบบ
        </button>
      </div>
    </aside>
  );
}