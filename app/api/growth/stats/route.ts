import { NextRequest, NextResponse } from 'next/server'
import { getAuthUser } from '@/lib/utils/auth'
import { getInstagramAccountByUserID } from '@/lib/services/instagram'
import { getGrowthStats } from '@/lib/services/growth'

export async function GET(req: NextRequest) {
  const claims = getAuthUser(req)
  
  if (!claims) {
    return NextResponse.json(
      { error: 'user not authenticated' },
      { status: 401 }
    )
  }
  
  const account = await getInstagramAccountByUserID(claims.user_id)
  
  if (!account) {
    return NextResponse.json(
      { error: 'no instagram account connected' },
      { status: 404 }
    )
  }
  
  // Get period from query param, default 30 days
  let periodDays = 30
  const period = req.nextUrl.searchParams.get('period')
  if (period) {
    switch (period) {
      case '7':
      case 'week':
        periodDays = 7
        break
      case '14':
        periodDays = 14
        break
      case '30':
      case 'month':
        periodDays = 30
        break
      case '90':
        periodDays = 90
        break
    }
  }
  
  try {
    const stats = await getGrowthStats(account.id, periodDays)
    
    return NextResponse.json({
      account: {
        id: account.id,
        username: account.username,
      },
      stats,
    })
  } catch (error: any) {
    console.error('GetGrowthStats error:', error)
    return NextResponse.json(
      { error: 'failed to calculate growth stats' },
      { status: 500 }
    )
  }
}

