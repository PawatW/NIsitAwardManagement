"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { jwtDecode } from "jwt-decode";

interface UserToken {
  user_id: string;
  email: string;
  role: "ADMIN" | "STUDENT" | "COMMITTEE_CHAIR" | "DEAN" | "HEAD_OF_DEPARTMENT" | "VICE_DEAN";
  exp: number;
  iat?: number;
}

interface AuthContextType {
  user: UserToken | null;
  loading: boolean;
  login: (token: string) => UserToken | null; 
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserToken | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();

  const login = (token: string) => {
    try {
      localStorage.setItem("nisit_token", token);
      const decoded = jwtDecode<UserToken>(token);
      setUser(decoded); 
      return decoded;
    } catch (error) {
      console.error("Login failed:", error);
      return null;
    }
  };

  const logout = () => {
    localStorage.removeItem("nisit_token");
    setUser(null);
    router.replace("/");
  };

  useEffect(() => {
    const checkAuth = () => {
      setLoading(true); 
      const token = localStorage.getItem("nisit_token");
  
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        const decoded = jwtDecode<UserToken>(token);
        if (decoded.exp * 1000 < Date.now()) {
          logout();
          return;
        }
        setUser(decoded);
      } catch (error) {
        logout();
      } finally {
        setLoading(false);
      }
    };

    checkAuth();

  }, [pathname]); 

  return (

    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);

// "use client";
// import { createContext, useContext, useEffect, useState } from "react";
// import { useRouter, usePathname } from "next/navigation";
// import { jwtDecode } from "jwt-decode";

// // 🛑 ตั้งเป็น true เพื่อเทส / ตั้งเป็น false เพื่อใช้งานจริง
// const BYPASS_LOGIN = true;  

// interface UserToken {
//   user_id: string;
//   email: string;
//   role: "ADMIN" | "STUDENT" | "COMMITTEE_CHAIR" | "DEAN" | "HEAD_OF_DEPARTMENT" | "VICE_DEAN";
//   exp: number;
//   iat?: number;
// }

// interface AuthContextType {
//   user: UserToken | null;
//   loading: boolean;
//   logout: () => void;
// }

// const AuthContext = createContext<AuthContextType>(null!);

// export function AuthProvider({ children }: { children: React.ReactNode }) {
//   const [user, setUser] = useState<UserToken | null>(null);
//   const [loading, setLoading] = useState(true);
//   const router = useRouter();
//   const pathname = usePathname();

//   const logout = () => {
//     localStorage.removeItem("nisit_token");
//     setUser(null);
//     router.replace("/");
//   };

//   useEffect(() => {
//     const checkAuth = () => {

//       if (BYPASS_LOGIN) {

//         const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjkyMDdlMGMtMmZkMC00MWE1LWFiMDUtZWUwYjhkMjg3ZTY3IiwiZW1haWwiOiJjb21taXR0ZWVAdGVzdC5jb20iLCJyb2xlIjoiQ09NTUlUVEVFX0NIQUlSIiwiZXhwIjoxNzczMjkyNTc3LCJpYXQiOjE3NzMyMDYxNzd9.zpXYAt-KbC214xxJAkZC_RQTZZv2ZmoFSjcVc-FfMSo";
//         localStorage.setItem("nisit_token", testToken);
        

//       }

//       const token = localStorage.getItem("nisit_token");
      
//       if (!token) {
//         setLoading(false);
//         return;
//       }

//       try {
//         const decoded = jwtDecode<UserToken>(token);
        
//         // ตรวจสอบวันหมดอายุ
//         if (decoded.exp * 1000 < Date.now()) {
//           logout();
//           return;
//         }

//         setUser(decoded);
//       } catch (error) {
//         console.error("Invalid token:", error);
//         logout();
//       } finally {
//         setLoading(false);
//       }
//     };

//     checkAuth();
//   }, [pathname]); 

//   return (
//     <AuthContext.Provider value={{ user, loading, logout }}>
//       {children}
//     </AuthContext.Provider>
//   );
// }

// export const useAuth = () => useContext(AuthContext);