import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    const billingApiUrl = process.env.BILLING_API_URL?.replace(/\/+$/, "");
    if (!billingApiUrl) return [];
    return [
      {
        source: "/backend/:path*",
        destination: `${billingApiUrl}/:path*`,
      },
    ];
  },
};

export default nextConfig;
