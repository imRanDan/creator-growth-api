import { NextRequest, NextResponse } from 'next/server'
import { getAuthUser } from '@/lib/utils/auth'
import { getInstagramAccountByUserID, getPostsByAccountID } from '@/lib/services/instagram'

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
  
  try {
    const posts = await getPostsByAccountID(account.id, 50)
    
    return NextResponse.json({
      account: {
        id: account.id,
        username: account.username,
      },
      posts_count: posts.length,
      posts,
    })
  } catch (error) {
    return NextResponse.json(
      { error: 'failed to fetch posts' },
      { status: 500 }
    )
  }
}

