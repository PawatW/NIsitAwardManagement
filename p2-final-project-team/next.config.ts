import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "lh3.googleusercontent.com", // อนุญาตโดเมนรูปของ Google
      },
      {
        protocol: "https",
        hostname: "*.googleusercontent.com", // เผื่อไว้สำหรับ subdomain อื่นๆ
      },
    ],
  },
};

export default nextConfig;