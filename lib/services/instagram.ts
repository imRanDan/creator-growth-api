import { sql } from '@/lib/db'
import { v4 as uuidv4 } from 'uuid'

export interface InstagramAccount {
  id: string
  userID: string
  igUserID: string
  username: string
  accessToken: string
  tokenExpiresAt: Date
}

export interface InstagramPost {
  id: string
  igPostID: string
  accountID: string
  caption: string
  mediaType: string
  mediaURL: string
  likeCount: number
  commentCount: number
  postedAt: Date
  fetchedAt: Date
}

export async function saveInstagramAccount(data: {
  userID: string
  igUserID: string
  username: string
  accessToken: string
  tokenExpiresAt: Date
}): Promise<InstagramAccount> {
  const id = uuidv4()
  const now = new Date()
  
  const result = await sql`
    INSERT INTO instagram_accounts
      (id, user_id, ig_user_id, username, access_token, token_expires_at, created_at, updated_at)
    VALUES (${id}, ${data.userID}, ${data.igUserID}, ${data.username}, ${data.accessToken}, ${data.tokenExpiresAt}, ${now}, ${now})
    ON CONFLICT (ig_user_id) DO UPDATE
      SET username = EXCLUDED.username,
          access_token = EXCLUDED.access_token,
          token_expires_at = EXCLUDED.token_expires_at,
          updated_at = EXCLUDED.updated_at
    RETURNING id, user_id, ig_user_id, username, access_token, token_expires_at
  `
  
  return {
    id: result.rows[0].id,
    userID: result.rows[0].user_id,
    igUserID: result.rows[0].ig_user_id,
    username: result.rows[0].username,
    accessToken: result.rows[0].access_token,
    tokenExpiresAt: result.rows[0].token_expires_at,
  }
}

export async function getInstagramAccountByUserID(userID: string): Promise<InstagramAccount | null> {
  const result = await sql`
    SELECT id, user_id, ig_user_id, username, access_token, token_expires_at
    FROM instagram_accounts
    WHERE user_id = ${userID}
    LIMIT 1
  `
  
  if (result.rows.length === 0) {
    return null
  }
  
  return {
    id: result.rows[0].id,
    userID: result.rows[0].user_id,
    igUserID: result.rows[0].ig_user_id,
    username: result.rows[0].username,
    accessToken: result.rows[0].access_token,
    tokenExpiresAt: result.rows[0].token_expires_at,
  }
}

export async function getInstagramAccountByID(id: string): Promise<InstagramAccount | null> {
  const result = await sql`
    SELECT id, user_id, ig_user_id, username, access_token, token_expires_at
    FROM instagram_accounts
    WHERE id = ${id}
    LIMIT 1
  `
  
  if (result.rows.length === 0) {
    return null
  }
  
  return {
    id: result.rows[0].id,
    userID: result.rows[0].user_id,
    igUserID: result.rows[0].ig_user_id,
    username: result.rows[0].username,
    accessToken: result.rows[0].access_token,
    tokenExpiresAt: result.rows[0].token_expires_at,
  }
}

export async function deleteInstagramAccountByUserID(userID: string): Promise<void> {
  await sql`DELETE FROM instagram_accounts WHERE user_id = ${userID}`
}

export async function fetchAndStorePosts(accountID: string): Promise<void> {
  const account = await getInstagramAccountByID(accountID)
  
  if (!account) {
    throw new Error('Account not found')
  }
  
  // Fetch media from Instagram
  const mediaResp = await fetch(
    `https://graph.instagram.com/me/media?fields=id,caption,media_type,media_url,timestamp,like_count,comments_count&access_token=${account.accessToken}&limit=50`
  )
  
  if (!mediaResp.ok) {
    throw new Error('Failed to fetch media')
  }
  
  const mediaData = await mediaResp.json()
  
  // Store posts
  for (const post of mediaData.data || []) {
    let postedAt = new Date()
    if (post.timestamp) {
      // Try to parse timestamp
      const parsed = new Date(post.timestamp)
      if (!isNaN(parsed.getTime())) {
        postedAt = parsed
      }
    }
    
    await sql`
      INSERT INTO instagram_posts
        (ig_post_id, account_id, caption, media_type, media_url, like_count, comments_count, posted_at, fetched_at)
      VALUES (${post.id}, ${accountID}, ${post.caption || ''}, ${post.media_type}, ${post.media_url}, ${post.like_count || 0}, ${post.comments_count || 0}, ${postedAt}, NOW())
      ON CONFLICT (account_id, ig_post_id) DO UPDATE
        SET caption = EXCLUDED.caption,
            media_type = EXCLUDED.media_type,
            media_url = EXCLUDED.media_url,
            like_count = EXCLUDED.like_count,
            comments_count = EXCLUDED.comments_count,
            posted_at = EXCLUDED.posted_at,
            fetched_at = NOW()
    `
  }
}

export async function getPostsByAccountID(accountID: string, limit: number = 50): Promise<InstagramPost[]> {
  const result = await sql`
    SELECT id, ig_post_id, account_id, caption, media_type, media_url, like_count, comments_count, posted_at, fetched_at
    FROM instagram_posts
    WHERE account_id = ${accountID}
    ORDER BY posted_at DESC
    LIMIT ${limit}
  `
  
  return result.rows.map((row: any) => ({
    id: row.id,
    igPostID: row.ig_post_id,
    accountID: row.account_id,
    caption: row.caption || '',
    mediaType: row.media_type || '',
    mediaURL: row.media_url || '',
    likeCount: row.like_count || 0,
    commentCount: row.comments_count || 0,
    postedAt: row.posted_at,
    fetchedAt: row.fetched_at,
  }))
}

