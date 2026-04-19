"use client";

import { useEffect, Suspense, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "../../../context/AuthContext"; 

function OAuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth(); 
  
  const hasProcessed = useRef(false);

  useEffect(() => {
    if (hasProcessed.current) return;
    
    const errorMsg = searchParams.get("msg");
    if (errorMsg) {
      hasProcessed.current = true;
      alert(errorMsg); 
      router.replace("/"); 
      return;
    }


    const accessToken = searchParams.get("accessToken");

    if (!accessToken) {
      router.replace("/");
      return;
    }

    hasProcessed.current = true;

    const decodedUser = login(accessToken);

    if (decodedUser) {
      switch (decodedUser.role) {
        case "ADMIN":
          router.replace("/admin/app");
          break;
        case "COMMITTEE_CHAIR":
          router.replace("/staff-page/app/review");
          break;
        case "DEAN":
        case "HEAD_OF_DEPARTMENT":
        case "VICE_DEAN":
          router.replace("/approval");
          break;
        case "STUDENT":
        default:
          router.replace("/dashboard");
          break;
      }
    } else {
      router.replace("/");
    }
  }, [router, searchParams, login]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-gray-500 animate-pulse">
        Checking your identity...
      </p>
    </div>
  );
}

export default function OAuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-gray-500 animate-pulse">Loading...</p>
      </div>
    }>
      <OAuthCallbackContent />
    </Suspense>
  );
}

// "use client";

// import { useEffect, Suspense } from "react";
// import { useRouter, useSearchParams } from "next/navigation";
// import { jwtDecode } from "jwt-decode"; 

// interface UserToken {
//   role: "ADMIN" | "STUDENT" | "COMMITTEE_CHAIR" | "DEAN" | "HEAD_OF_DEPARTMENT" | "VICE_DEAN";
  
// }

// function OAuthCallbackContent() {
//   const router = useRouter();
//   const searchParams = useSearchParams();

//   useEffect(() => {
//     const accessToken = searchParams.get("accessToken");

//     if (!accessToken) {
//       router.replace("/");
//       return;
//     }

//    localStorage.setItem("nisit_token", accessToken);

//     try {
//       const decoded = jwtDecode<UserToken>(accessToken);

//       switch (decoded.role) {
//         case "ADMIN":
//           router.replace("/admin/app");
//           break;
//         case "COMMITTEE_CHAIR":
//           router.replace("/staff-page/app/review");
//           break;
//         case "DEAN":
//         case "HEAD_OF_DEPARTMENT":
//         case "VICE_DEAN":
//           router.replace("/approval");
//           break;
//         case "STUDENT":
//         default:
//           router.replace("/dashboard");
//           break;
//       }
//     } catch (error) {
//       console.error("Token invalid:", error);
//       router.replace("/");
//     }
//   }, [router, searchParams]);

//   return (
//     <div className="flex min-h-screen items-center justify-center">
//       <p className="text-gray-500 animate-pulse">
//         Checking your identity...
//       </p>
//     </div>
//   );
// }


// export default function OAuthCallbackPage() {
//   return (
//     <Suspense fallback={
//       <div className="flex min-h-screen items-center justify-center">
//         <p className="text-gray-500 animate-pulse">Loading...</p>
//       </div>
//     }>
//       <OAuthCallbackContent />
//     </Suspense>
//   );
// }