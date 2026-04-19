"use client";

import { usePathname, useRouter } from "next/navigation"; 
import { useEffect } from "react"; 
import { useAuth } from "../../context/AuthContext"; 
import Sidebar from "../../component/layouts/StudentSidebar";
import Navbar from "../../component/layouts/Navbar";
import { Loader2 } from "lucide-react"; 

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, loading } = useAuth(); 
  const pathname = usePathname();
  const router = useRouter(); 

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/"); 
    }
  }, [user, loading, router]);


  const showNavbarPaths = [
    "/application",
    "/approval",
    

  ];

  const showSidebarPaths = [
    "/dashboard",
    "/my-applications",
    
  ];

  const shouldShow = (pathList: string[]) => {
 
    return pathList.some((path) => pathname === path || pathname.startsWith(`${path}/`));
  };

  const showNavbar = shouldShow(showNavbarPaths);
  const showSidebar = shouldShow(showSidebarPaths);

  if (loading) {
    return (
      <div className="h-screen w-full flex flex-col items-center justify-center bg-gray-50">
        <Loader2 className="animate-spin text-green-600 mb-4" size={40} />
        <p className="text-gray-500 animate-pulse">Checking authentication...</p>
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50">
      {showNavbar && (
        <div className="fixed top-0 left-0 right-0 z-50">
          <Navbar />
        </div>
      )}

      <div className="flex">
        {showSidebar && (
          <div className={`fixed left-0 h-full w-64 ${showNavbar ? "top-16" : "top-0"}`}>
             <Sidebar />
          </div>
        )}

        <main 
          className={`
            flex-1 p-8 transition-all duration-300
            ${showNavbar ? "pt-20" : "pt-8"}
            ${showSidebar ? "ml-64" : "ml-0"}
          `}
        >
          {children}
        </main>
      </div>
    </div>
  );
}