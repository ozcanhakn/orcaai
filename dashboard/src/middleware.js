import { NextResponse } from 'next/server';

export function middleware(request) {
  // Check if the user is accessing protected routes
  if (request.nextUrl.pathname.startsWith('/dashboard') || 
      request.nextUrl.pathname.startsWith('/playground')) {
    
    // Check for token in localStorage
    const token = request.cookies.get('token');
    
    if (!token) {
      // Redirect to login if no token found
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/playground/:path*']
};
