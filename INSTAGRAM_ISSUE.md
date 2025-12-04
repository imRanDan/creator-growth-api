# 🐛 Instagram Business Account Data Not Fetching

## Status: OAuth Working ✅ | Data Fetching ❌

## Problem

Instagram OAuth flow completes successfully, but no posts are being fetched. The app shows 0 posts, likes, comments, and engagement.

## Root Cause

**The OAuth flow connects the Facebook account instead of the Instagram Business Account.**

Current flow:
1. ✅ User authorizes via Facebook OAuth
2. ✅ Backend receives Facebook access token
3. ✅ Token saved to database
4. ❌ Backend queries `GET https://graph.instagram.com/me/media`
5. ❌ "me" refers to Facebook user, not Instagram Business Account
6. ❌ Returns 0 posts because Facebook user != Instagram account

## What's Working

- ✅ Frontend: React app deployed on Vercel
- ✅ Backend: Go API deployed on Railway
- ✅ OAuth: Facebook login completes successfully
- ✅ Database: Instagram account record created
- ✅ Token storage: Long-lived token saved
- ✅ API calls: `/api/instagram/refresh` returns 202
- ✅ Meta App: Domains and redirect URIs configured

## What Needs Fixing

### Backend Changes Required

**File: `internal/services/instagram.go`**

The Instagram Business API requires:
1. Get Facebook Pages the user manages
2. Get Instagram Business Account ID linked to each page
3. Store the Instagram Business Account ID (not just Facebook user ID)
4. Use Instagram Business Account ID for all media queries

**Current (incorrect):**
```go
// Calls /me/media which queries Facebook user
meURL := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username&access_token=%s", accessToken)
```

**Should be:**
```go
// 1. Get Facebook Pages
GET https://graph.facebook.com/me/accounts?access_token={token}

// 2. For each page, get linked Instagram Business Account
GET https://graph.facebook.com/{page-id}?fields=instagram_business_account&access_token={token}

// 3. Store instagram_business_account.id in database

// 4. Query posts using IG Business Account ID
GET https://graph.facebook.com/{ig-business-account-id}/media?fields=...&access_token={token}
```

### Database Schema Update

**Table: `instagram_accounts`**

Add column:
```sql
ALTER TABLE instagram_accounts ADD COLUMN fb_page_id VARCHAR(255);
```

Current columns store:
- `ig_user_id` (Instagram user ID) ✅
- `username` (Instagram username) ✅
- `access_token` (Facebook token) ✅

Need to also store:
- `fb_page_id` (Facebook Page ID)
- Use `ig_user_id` as the Instagram Business Account ID

## Testing Requirements

**Prerequisite:**
- Instagram account must be **Business or Creator** type
- Instagram must be **linked to a Facebook Page**
- Account must have **posts**

**Test accounts:**
- @5amdany (Instagram) - linked to Facebook Page
- Check Meta Accounts Center to verify connection

## References

- [Instagram Graph API - Get User Media](https://developers.facebook.com/docs/instagram-api/reference/ig-user/media)
- [Facebook Pages API](https://developers.facebook.com/docs/graph-api/reference/page)
- [Get Instagram Business Account from Page](https://developers.facebook.com/docs/instagram-api/getting-started#get-pages)

## Priority

**Medium** - App works end-to-end except data fetching. Critical for production but can demo with mock data for interview.

## Next Steps

1. [ ] Update `InstagramCallback` to fetch Facebook Pages
2. [ ] Get Instagram Business Account ID from page
3. [ ] Store IG Business Account ID in database
4. [ ] Update `FetchUserMedia` to use IG Business Account ID instead of "me"
5. [ ] Test with @5amdany account
6. [ ] Verify posts are fetched and displayed

---

**Created:** Dec 3, 2025  
**Interview:** Next week (Salt XC - $90K)  
**Status:** OAuth infrastructure complete, data fetching needs fix


