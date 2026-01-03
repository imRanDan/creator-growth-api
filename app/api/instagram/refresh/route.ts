import { NextRequest, NextResponse } from 'next/server'
import { getAuthUser } from '@/lib/utils/auth'
import { getInstagramAccountByUserID, fetchAndStorePosts } from '@/lib/services/instagram'

export async function POST(req: NextRequest) {
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
  
  // Trigger fetch in background (non-blocking)
  fetchAndStorePosts(account.id).catch(err => {
    console.error('FetchAndStorePosts error:', err)
  })
  
  return NextResponse.json({
    status: 'fetch scheduled',
    account: {
      id: account.id,
      username: account.username,
    },
  }, { status: 202 })
}

