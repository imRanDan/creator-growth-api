import { NextRequest, NextResponse } from 'next/server'
import { getAuthUser } from '@/lib/utils/auth'
import { getInstagramAccountByUserID, deleteInstagramAccountByUserID } from '@/lib/services/instagram'

export async function DELETE(req: NextRequest) {
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
    await deleteInstagramAccountByUserID(claims.user_id)
    
    return NextResponse.json({
      message: 'Instagram account disconnected successfully',
      account: {
        id: account.id,
        username: account.username,
      },
    })
  } catch (error: any) {
    console.error('DisconnectInstagram error:', error)
    return NextResponse.json(
      {
        error: 'failed to disconnect instagram account',
        details: error.message,
      },
      { status: 500 }
    )
  }
}

