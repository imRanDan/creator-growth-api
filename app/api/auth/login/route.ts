import { NextRequest, NextResponse } from 'next/server'
import { authenticateUser } from '@/lib/services/auth'
import { generateToken } from '@/lib/services/jwt'

export async function POST(req: NextRequest) {
  try {
    const { email, password } = await req.json()
    
    if (!email || !password) {
      return NextResponse.json(
        { error: 'Email and password are required' },
        { status: 400 }
      )
    }
    
    const user = await authenticateUser(email, password)
    const token = generateToken(user.id, user.email)
    
    return NextResponse.json({
      token,
      token_type: 'bearer',
      expires_in: 24 * 60 * 60,
      user: {
        id: user.id,
        email: user.email,
        created_at: user.created_at,
      },
    })
  } catch (error: any) {
    return NextResponse.json(
      { error: 'Invalid email or password' },
      { status: 401 }
    )
  }
}

