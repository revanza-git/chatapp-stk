import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone',
  eslint: {
    // Disable ESLint during builds for Docker
    ignoreDuringBuilds: true,
  },
  typescript: {
    // Ignore TypeScript errors during builds for Docker
    ignoreBuildErrors: true,
  },
};

export default nextConfig;
