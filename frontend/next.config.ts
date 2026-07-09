import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  devIndicators: false,
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        // Server-side proxy target: prefer API_INTERNAL_BASE_URL (e.g. the backend's
        // service name on the docker-compose network), since NEXT_PUBLIC_API_BASE_URL
        // is baked into the browser bundle and points at a host reachable from the
        // user's browser, not necessarily from inside this container.
        destination: `${process.env.API_INTERNAL_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://31.76.100.179:5000/v1"}/:path*`,
      },
    ];
  },
};

export default nextConfig;
