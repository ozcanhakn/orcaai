/** @type {import('next').NextConfig} */
const nextConfig = {
  // Specify the output tracing root to avoid workspace detection issues
  outputFileTracingRoot: __dirname,
  
  // Ensure we're using the pages directory
  pageExtensions: ['js', 'jsx', 'ts', 'tsx'],
  
  // Specify the source directory
  dir: 'src',
  
  // Experimental features
  experimental: {
    // Enable if you're using app directory
    appDir: false,
  },
}

module.exports = nextConfig