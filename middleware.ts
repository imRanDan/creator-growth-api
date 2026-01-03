import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { validateToken } from '@/lib/services/jwt'

export function middleware(request: NextRequest) {
  // Public routes
  const publicRoutes = ['/api/auth', '/api/waitlist', '/auth/instagram/callback']
  const isPublic = publicRoutes.some(route => request.nextUrl.pathname.startsWith(route))
  
  if (isPublic) {
    return NextResponse.next()
  }
  
  // Protected API routes
  if (request.nextUrl.pathname.startsWith('/api')) {
    const authHeader = request.headers.get('authorization')
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      )
    }
    
    const token = authHeader.replace('Bearer ', '')
    
    try {
      validateToken(token)
      return NextResponse.next()
    } catch {
      return NextResponse.json(
        { error: 'Invalid token' },
        { status: 401 }
      )
    }
  }
  
  return NextResponse.next()
}

export const config = {
  matcher: ['/api/:path*'],
}

