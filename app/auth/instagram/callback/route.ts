import { NextRequest, NextResponse } from 'next/server'
import { validateToken } from '@/lib/services/jwt'
import { saveInstagramAccount, fetchAndStorePosts } from '@/lib/services/instagram'

export async function GET(req: NextRequest) {
  const code = req.nextUrl.searchParams.get('code')
  const state = req.nextUrl.searchParams.get('state')
  
  if (!code || !state) {
    return NextResponse.json(
      { error: 'missing code or state' },
      { status: 400 }
    )
  }
  
  try {
    const claims = validateToken(state)
    const userID = claims.user_id
    
    const clientID = process.env.INSTAGRAM_CLIENT_ID!
    const clientSecret = process.env.INSTAGRAM_CLIENT_SECRET!
    const redirectURI = process.env.INSTAGRAM_REDIRECT_URI!
    
    // Exchange code for short-lived token
    const shortResp = await fetch('https://api.instagram.com/oauth/access_token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({
        client_id: clientID,
        client_secret: clientSecret,
        grant_type: 'authorization_code',
        redirect_uri: redirectURI,
        code,
      }),
    })
    
    if (!shortResp.ok) {
      throw new Error('Token exchange failed')
    }
    
    const shortData = await shortResp.json()
    
    // Exchange for long-lived token
    const longURL = `https://graph.instagram.com/access_token?grant_type=ig_exchange_token&client_secret=${clientSecret}&access_token=${shortData.access_token}`
    const longResp = await fetch(longURL)
    
    if (!longResp.ok) {
      throw new Error('Failed to get long-lived token')
    }
    
    const longData = await longResp.json()
    
    // Get IG user info
    const meResp = await fetch(`https://graph.instagram.com/me?fields=id,username&access_token=${longData.access_token}`)
    if (!meResp.ok) {
      throw new Error('Failed to fetch IG user')
    }
    
    const me = await meResp.json()
    
    // Save account
    const expiresAt = new Date(Date.now() + longData.expires_in * 1000)
    const account = await saveInstagramAccount({
      userID,
      igUserID: me.id,
      username: me.username,
      accessToken: longData.access_token,
      tokenExpiresAt: expiresAt,
    })
    
    // Auto-fetch posts in background
    fetchAndStorePosts(account.id).catch(err => {
      console.error('Auto-fetch posts error:', err)
    })
    
    // Redirect to frontend
    const frontendURL = process.env.FRONTEND_URL || 'http://localhost:3000'
    return NextResponse.redirect(`${frontendURL}/dashboard?connected=true`)
  } catch (error: any) {
    console.error('Instagram callback error:', error)
    return NextResponse.json(
      { error: error.message || 'Failed to connect Instagram' },
      { status: 500 }
    )
  }
}

