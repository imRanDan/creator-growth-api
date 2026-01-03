import { NextRequest } from 'next/server'
import { validateToken, Claims } from '@/lib/services/jwt'

export function getAuthUser(req: NextRequest): Claims | null {
  const authHeader = req.headers.get('authorization')
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return null
  }
  
  const token = authHeader.replace('Bearer ', '')
  
  try {
    return validateToken(token)
  } catch {
    return null
  }
}

