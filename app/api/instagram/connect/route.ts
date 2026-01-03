import { NextRequest, NextResponse } from 'next/server'
import { getAuthUser } from '@/lib/utils/auth'
import { generateStateToken } from '@/lib/services/jwt'

export async function GET(req: NextRequest) {
  const claims = getAuthUser(req)
  
  if (!claims) {
    return NextResponse.json(
      { error: 'Unauthorized' },
      { status: 401 }
    )
  }
  
  const clientID = process.env.INSTAGRAM_CLIENT_ID
  const redirectURI = process.env.INSTAGRAM_REDIRECT_URI
  
  if (!clientID || !redirectURI) {
    return NextResponse.json(
      { error: 'Instagram app configuration missing' },
      { status: 500 }
    )
  }
  
  const state = generateStateToken(claims.user_id, claims.email)
  
  const authURL = `https://www.facebook.com/v18.0/dialog/oauth?client_id=${clientID}&redirect_uri=${encodeURIComponent(redirectURI)}&scope=instagram_basic,pages_show_list,pages_read_engagement,business_management&response_type=code&state=${state}`
  
  return NextResponse.json({ url: authURL })
}

