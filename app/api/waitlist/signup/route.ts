import { NextRequest, NextResponse } from 'next/server'
import { getPool } from '@/lib/db'
import { sendWelcomeEmail } from '@/lib/services/email'

const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/

export async function POST(req: NextRequest) {
  try {
    const { email } = await req.json()
    
    if (!email) {
      return NextResponse.json(
        { error: 'Email is required' },
        { status: 400 }
      )
    }
    
    if (!emailRegex.test(email)) {
      return NextResponse.json(
        { error: 'Invalid email format' },
        { status: 400 }
      )
    }
    
    const pool = getPool()
    
    // Insert with ON CONFLICT
    const result = await pool.query(
      `INSERT INTO waitlist (email)
       VALUES ($1)
       ON CONFLICT (email) DO NOTHING
       RETURNING id`,
      [email]
    )
    
    const isNewSignup = result.rows.length > 0
    
    // Send email in background (non-blocking)
    if (isNewSignup) {
      sendWelcomeEmail(email).catch(err => {
        console.error(`Failed to send welcome email to ${email}:`, err)
      })
    }
    
    return NextResponse.json({
      message: 'Successfully added to waitlist!',
      email,
    })
  } catch (error: any) {
    console.error('Waitlist signup error:', error)
    return NextResponse.json(
      { error: 'Failed to add email to waitlist', details: error.message },
      { status: 500 }
    )
  }
}

