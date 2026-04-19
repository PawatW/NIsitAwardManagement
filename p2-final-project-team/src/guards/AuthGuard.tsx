"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation"; 
import { useAuth } from "../context/AuthContext"; 
import { Loader2 } from "lucide-react";

// กำหนดสิทธิ์ว่า Role ไหน เข้า Path ไหนได้บ้าง
const roleAccessMap: Record<string, string[]> = {
  ADMIN: ["/admin"],
  COMMITTEE_CHAIR: ["/staff-page"],
  DEAN: ["/approval"],
  HEAD_OF_DEPARTMENT: ["/approval"], 
  VICE_DEAN: ["/approval"],
  STUDENT: ["/dashboard", "/application", "/my-applications"],
};

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth(); 
  const router = useRouter();
  const pathname = usePathname(); 

  const publicPaths = ["/", "/register", "/oauth/callback", "/announcements"]; 

  const isPublicPath = publicPaths.includes(pathname);
  const allowedPaths = user ? roleAccessMap[user.role] || [] : [];
  const hasAccess = user ? allowedPaths.some(path => pathname.startsWith(path)) : false;

  useEffect(() => {
    if (loading) return;

    if (!user && !isPublicPath) {
      router.replace("/"); 
      return;
    }

    if (user && !isPublicPath && !hasAccess) {
      console.warn(`Access Denied: ${user.role} cannot access ${pathname}`);
      switch (user.role) {
        case "ADMIN": router.replace("/admin/app"); break;
        case "COMMITTEE_CHAIR": router.replace("/staff-page/app/review"); break;
        case "DEAN": 
        case "HEAD_OF_DEPARTMENT": 
        case "VICE_DEAN": router.replace("/approval"); break;
        case "STUDENT": 
        default: router.replace("/dashboard"); break;
      }
    }
  }, [user, loading, router, pathname, isPublicPath, hasAccess]); 


  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
         <Loader2 className="w-10 h-10 animate-spin text-green-600" />
         <p className="ml-2 text-gray-500 font-medium">Checking permissions...</p>
      </div>
    );
  }

  if (!user && !isPublicPath) return null;

 
  if (user && !isPublicPath && !hasAccess) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
         <Loader2 className="w-10 h-10 animate-spin text-red-500" />
         <p className="ml-2 text-gray-500 font-medium">Redirecting to your workspace...</p>
      </div>
    );
  }

  
  return <>{children}</>;
}