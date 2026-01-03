import { NextRequest, NextResponse } from 'next/server'
import { sql } from '@/lib/db'

export async function GET(req: NextRequest) {
  const adminPassword = req.headers.get('x-admin-password') || 
                        req.nextUrl.searchParams.get('password')
  
  const expectedPassword = process.env.ADMIN_PASSWORD || 'changeme'
  
  if (adminPassword !== expectedPassword) {
    return NextResponse.json(
      { error: 'Unauthorized' },
      { status: 401 }
    )
  }
  
  const page = Math.max(1, parseInt(req.nextUrl.searchParams.get('page') || '1'))
  const limit = Math.min(100, Math.max(1, parseInt(req.nextUrl.searchParams.get('limit') || '50')))
  const offset = (page - 1) * limit
  
  try {
    // Get total count
    const countResult = await sql`SELECT COUNT(*) FROM waitlist`
    const total = parseInt(countResult.rows[0].count)
    
    // Get entries
    const entries = await sql`
      SELECT id, email, created_at
      FROM waitlist
      ORDER BY created_at DESC
      LIMIT ${limit}
      OFFSET ${offset}
    `
    
    return NextResponse.json({
      entries: entries.rows,
      total,
      page,
      limit,
      pages: Math.ceil(total / limit),
    })
  } catch (error) {
    return NextResponse.json(
      { error: 'Failed to fetch entries' },
      { status: 500 }
    )
  }
}

