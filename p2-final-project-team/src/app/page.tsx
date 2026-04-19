"use client";

import Image from "next/image";
import { FcGoogle } from "react-icons/fc";
import { HiSparkles } from "react-icons/hi2";
import { useState } from "react";

export default function LoginPage() {

  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async () => {
    setIsLoading(true);

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/google`, {
        method: "GET",
      });

      if (!res.ok) throw new Error("Failed to get auth url");

      const data = await res.json(); 

      if (data.auth_url) {
        window.location.href = data.auth_url; 
      }

    } catch (error) {
      console.error("Login Error:", error);
      setIsLoading(false); 
      alert("เกิดข้อผิดพลาดในการเชื่อมต่อกับ Google");
    }
  };

  return (
    <div className="flex min-h-screen w-full bg-white">
      {/* Background Image */}
      <div className="hidden lg:block lg:w-1/2 relative">
        <Image
          src="/login-bg.jpg" 
          alt="Kasetsart University Building"
          fill
          className="object-cover"
          priority
        />
      </div>

      {/* Login */}
      <div className="w-full lg:w-1/2 flex flex-col justify-center items-center p-8 lg:p-24">
        <div className="w-full max-w-md space-y-8">
          
          {/* Header */}
          <div className="text-left space-y-2">
            <div className="flex items-center gap-2 mb-6">
              <HiSparkles className="text-green-700 w-8 h-8" /> 
              <span className="text-2xl font-bold text-black">Nisit Deeden</span>
            </div>
            
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">
              Welcome 
            </h1>
            <p className="text-sm text-gray-500">
              Sign in via Google Account
            </p>
          </div>

          <div className="mt-8">
              <button
              onClick={handleLogin} 
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-3 bg-white border border-gray-300 rounded-lg px-4 py-3 text-sm font-medium text-gray-700 hover:bg-gray-50 hover:shadow-sm transition-all"
            >
              {isLoading ? (
                <span className="text-gray-500">Redirecting to Google...</span>
              ) : (
                <>
                  <FcGoogle className="w-6 h-6" />
                  <span>Sign in with Google</span>
                </>
              )}
            </button>
          </div>

        </div>
      </div>
    </div>
  );
}