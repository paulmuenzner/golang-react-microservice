
// app/frontend/public/next.config.js
const nextConfig = {
  reactStrictMode: true,

  // Set root for Turbopack
  turbopack: {
    root: __dirname, // Current directory
  },
  
  // Production optimization
  output: 'standalone',
}

export default nextConfig;
