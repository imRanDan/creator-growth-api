import { NextRequest, NextResponse } from 'next/server'
import { createUser } from '@/lib/services/auth'
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
    
    if (password.length < 8) {
      return NextResponse.json(
        { error: 'Password must be at least 8 characters' },
        { status: 400 }
      )
    }
    
    const user = await createUser(email, password)
    const token = generateToken(user.id, user.email)
    
    return NextResponse.json({
      token,
      user: {
        id: user.id,
        email: user.email,
        created_at: user.created_at,
      },
    }, { status: 201 })
  } catch (error: any) {
    if (error.message === 'Email already in use') {
      return NextResponse.json(
        { error: 'Email already registered' },
        { status: 409 }
      )
    }
    
    return NextResponse.json(
      { error: 'Failed to create account' },
      { status: 500 }
    )
  }
}

