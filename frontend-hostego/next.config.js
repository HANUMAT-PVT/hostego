/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    NEXT_PUBLIC_AWS_REGION: process.env.NEXT_PUBLIC_AWS_REGION,
    NEXT_PUBLIC_AWS_ACCESS_KEY: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
    NEXT_PUBLIC_AWS_SECRET_KEY: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
    NEXT_PUBLIC_AWS_BUCKET_NAME: process.env.NEXT_PUBLIC_AWS_BUCKET_NAME,
  },
  // Add any production-specific configurations
  productionBrowserSourceMaps: true, // Enable source maps in production
}

module.exports = nextConfig 