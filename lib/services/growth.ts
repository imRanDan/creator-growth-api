import { sql } from '@/lib/db'

export interface GrowthStats {
  total_posts: number
  total_likes: number
  total_comments: number
  total_engagement: number
  avg_likes_per_post: number
  avg_comments_per_post: number
  engagement_rate: number
  best_post?: {
    id: string
    caption: string
    media_type: string
    media_url: string
    like_count: number
    comment_count: number
    engagement: number
    posted_at: Date
  }
  likes_trend: number
  comments_trend: number
  posting_trend: number
  posts_this_week: number
  posts_this_month: number
  period_days: number
  message: string
}

export async function getGrowthStats(accountID: string, periodDays: number = 30): Promise<GrowthStats> {
  if (periodDays <= 0) {
    periodDays = 30
  }

  // Get current period stats
  const currentResult = await sql`
    SELECT 
      COUNT(*)::int as post_count,
      COALESCE(SUM(like_count), 0)::int as total_likes,
      COALESCE(SUM(comments_count), 0)::int as total_comments
    FROM instagram_posts
    WHERE account_id = ${accountID}
      AND posted_at >= NOW() - INTERVAL '1 day' * ${periodDays}
  `

  const current = currentResult.rows[0]
  const totalPosts = current.post_count || 0
  const totalLikes = current.total_likes || 0
  const totalComments = current.total_comments || 0
  const totalEngagement = totalLikes + totalComments

  // Calculate averages
  const avgLikesPerPost = totalPosts > 0 ? totalLikes / totalPosts : 0
  const avgCommentsPerPost = totalPosts > 0 ? totalComments / totalPosts : 0
  const engagementRate = totalPosts > 0 ? totalEngagement / totalPosts : 0

  // Get previous period stats for trends
  const prevResult = await sql`
    SELECT 
      COUNT(*)::int as post_count,
      COALESCE(SUM(like_count), 0)::int as total_likes,
      COALESCE(SUM(comments_count), 0)::int as total_comments
    FROM instagram_posts
    WHERE account_id = ${accountID}
      AND posted_at >= NOW() - INTERVAL '1 day' * ${periodDays * 2}
      AND posted_at < NOW() - INTERVAL '1 day' * ${periodDays}
  `

  let likesTrend = 0
  let commentsTrend = 0
  let postingTrend = 0

  if (prevResult.rows.length > 0) {
    const prev = prevResult.rows[0]
    const prevPosts = prev.post_count || 0
    const prevLikes = prev.total_likes || 0
    const prevComments = prev.total_comments || 0

    if (prevLikes > 0) {
      likesTrend = ((totalLikes - prevLikes) / prevLikes) * 100
    }
    if (prevComments > 0) {
      commentsTrend = ((totalComments - prevComments) / prevComments) * 100
    }
    if (prevPosts > 0) {
      postingTrend = ((totalPosts - prevPosts) / prevPosts) * 100
    }
  }

  // Get best performing post
  const bestResult = await sql`
    SELECT id, caption, media_type, media_url, like_count, comments_count, posted_at
    FROM instagram_posts
    WHERE account_id = ${accountID}
      AND posted_at >= NOW() - INTERVAL '1 day' * ${periodDays}
    ORDER BY (like_count + comments_count) DESC
    LIMIT 1
  `

  let bestPost = undefined
  if (bestResult.rows.length > 0) {
    const best = bestResult.rows[0]
    let caption = best.caption || ''
    if (caption.length > 100) {
      caption = caption.substring(0, 100) + '...'
    }
    bestPost = {
      id: best.id,
      caption,
      media_type: best.media_type || '',
      media_url: best.media_url || '',
      like_count: best.like_count || 0,
      comment_count: best.comments_count || 0,
      engagement: (best.like_count || 0) + (best.comments_count || 0),
      posted_at: best.posted_at,
    }
  }

  // Get posting frequency
  const weekResult = await sql`
    SELECT COUNT(*)::int as count
    FROM instagram_posts
    WHERE account_id = ${accountID}
      AND posted_at >= NOW() - INTERVAL '7 days'
  `
  const postsThisWeek = weekResult.rows[0]?.count || 0

  const monthResult = await sql`
    SELECT COUNT(*)::int as count
    FROM instagram_posts
    WHERE account_id = ${accountID}
      AND posted_at >= NOW() - INTERVAL '30 days'
  `
  const postsThisMonth = monthResult.rows[0]?.count || 0

  // Generate friendly message
  let message = ''
  if (totalPosts === 0) {
    message = "No posts yet in this period. Time to share something! 📸"
  } else {
    if (likesTrend > 20) {
      message = "🔥 You're on fire! Engagement is way up."
    } else if (likesTrend > 5) {
      message = "📈 Nice! You're growing steadily."
    } else if (likesTrend > -5) {
      message = "😎 Holding steady - keep doing your thing."
    } else if (likesTrend > -20) {
      message = "📉 Slight dip, but no worries - it happens."
    } else {
      message = "💪 Engagement is down, but consistency is key!"
    }

    if (postsThisWeek === 0) {
      message += " Haven't posted this week though - your audience misses you!"
    } else if (postsThisWeek >= 5) {
      message += " You've been posting a lot - great hustle!"
    }
  }

  return {
    total_posts: totalPosts,
    total_likes: totalLikes,
    total_comments: totalComments,
    total_engagement: totalEngagement,
    avg_likes_per_post: avgLikesPerPost,
    avg_comments_per_post: avgCommentsPerPost,
    engagement_rate: engagementRate,
    best_post: bestPost,
    likes_trend: likesTrend,
    comments_trend: commentsTrend,
    posting_trend: postingTrend,
    posts_this_week: postsThisWeek,
    posts_this_month: postsThisMonth,
    period_days: periodDays,
    message,
  }
}

