/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  trailingSlash: true,
  async rewrites() {
      return [
          {
              source: "/api/:path*",
              destination: "http://localhost:8081/api/:path*",
          },
      ];
  },
  images: {
    remotePatterns: [
      {
        hostname: 'tailwindui.com',
      },
      {
        hostname: 'images.unsplash.com',
      },
    ],
  },
}

module.exports = nextConfig
