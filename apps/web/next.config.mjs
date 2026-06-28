/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  // Permite que variáveis NEXT_PUBLIC_* sejam injetadas em build time
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000/api",
    NEXT_PUBLIC_SOCKET_URL: process.env.NEXT_PUBLIC_SOCKET_URL ?? "ws://localhost:8765",
  },
};

export default nextConfig;