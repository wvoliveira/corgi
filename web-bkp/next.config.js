/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  trailingSlash: true,
  
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  },

  async rewrites() {
      return [
          {
              source: "/api/:path*",
              destination: "http://localhost:8081/api/:path*",
          },
      ];
  },
}

module.exports = nextConfig
