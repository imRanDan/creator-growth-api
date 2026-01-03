# Creator Growth Tool - Technical Interview Cheat Sheet

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.25.3
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **OAuth**: Instagram Graph API (Facebook OAuth 2.0)
- **Security**: bcrypt (password hashing), ulule/limiter (rate limiting)
- **Deployment**: Railway (Production)

### Frontend
- **Framework**: React 19 (with Vite, not Next.js)
- **Build Tool**: Vite 7.2.4
- **Styling**: Tailwind CSS 4
- **State Management**: React Hooks (useState, useEffect)
- **Routing**: React Router DOM 7
- **HTTP Client**: Fetch API
- **Analytics**: Vercel Analytics
- **Data Visualization**: Custom components (no chart libraries)
  - Card-based stat displays
  - Color-coded trend indicators
  - Glassmorphism UI design

### DevOps & Infrastructure
- PostgreSQL database with automatic migrations
- Environment-based configuration (.env)
- CORS middleware for cross-origin requests
- Health check endpoints (`/health`)

---

## 🎯 Why Go for Backend?

### Performance & Concurrency
- **Goroutines**: Lightweight concurrency for background jobs (token refresh, post fetching)
- **Compiled Language**: Fast execution, low memory footprint
- **Excellent for APIs**: Built-in HTTP server, efficient JSON handling

### Production-Ready Features
- **Strong Typing**: Compile-time error catching
- **Standard Library**: Rich stdlib for HTTP, JSON, crypto, database
- **Simple Deployment**: Single binary, no runtime dependencies
- **Database Drivers**: Mature PostgreSQL driver (lib/pq)

### Specific Use Cases in This Project
- Background token refresh job runs in goroutine (non-blocking)
- Concurrent post fetching after OAuth connection
- Efficient handling of Instagram API rate limits
- Fast response times for analytics queries

---

## 🗄️ Database Schema & Design Decisions

### Tables

#### `users`
```sql
- id: UUID (PRIMARY KEY, gen_random_uuid())
- email: VARCHAR(255) UNIQUE NOT NULL
- password: VARCHAR(255) NOT NULL (bcrypt hashed)
- created_at, updated_at: TIMESTAMP WITH TIME ZONE
```

**Design Decisions:**
- UUIDs for distributed system compatibility
- Email uniqueness enforced at DB level
- Passwords never returned in JSON (excluded with `json:"-"`)

#### `instagram_accounts`
```sql
- id: UUID (PRIMARY KEY)
- user_id: UUID REFERENCES users(id) ON DELETE CASCADE
- ig_user_id: TEXT UNIQUE NOT NULL (Instagram's user ID)
- username: TEXT
- access_token: TEXT NOT NULL (long-lived token)
- token_expires_at: TIMESTAMP WITH TIME ZONE
- followers: BIGINT (nullable)
- created_at, updated_at: TIMESTAMP WITH TIME ZONE
```

**Design Decisions:**
- One-to-many relationship (user can have multiple IG accounts)
- `ON DELETE CASCADE` ensures data consistency
- `ig_user_id` unique to prevent duplicate connections
- Token expiry tracking for automatic refresh

#### `instagram_posts`
```sql
- id: UUID (PRIMARY KEY)
- ig_post_id: TEXT NOT NULL
- account_id: UUID REFERENCES instagram_accounts(id) ON DELETE CASCADE
- caption: TEXT
- media_type: TEXT
- media_url: TEXT
- like_count: INT
- comments_count: INT
- posted_at: TIMESTAMP WITH TIME ZONE
- fetched_at: TIMESTAMP WITH TIME ZONE
- UNIQUE(account_id, ig_post_id)
```

**Design Decisions:**
- Composite unique constraint prevents duplicate posts per account
- Stores engagement metrics (likes, comments) for analytics
- `fetched_at` tracks when data was last synced
- CASCADE delete removes posts when account disconnected

#### `waitlist`
```sql
- id: UUID (PRIMARY KEY)
- email: VARCHAR(255) UNIQUE NOT NULL
- created_at: TIMESTAMP WITH TIME ZONE
```

**Design Decisions:**
- Simple table for pre-launch email collection
- Email uniqueness prevents duplicates

### Indexing Strategy
- Primary keys automatically indexed
- Unique constraints create indexes
- Foreign keys create indexes for JOIN performance
- Consider adding index on `posted_at` for time-range queries (future optimization)

### Migration Strategy
- Migrations run automatically on server startup
- `CREATE TABLE IF NOT EXISTS` for idempotency
- Uses PostgreSQL extensions (`pgcrypto` for UUID generation)

---

## 🔐 OAuth 2.0 Flow Implementation

### Instagram OAuth Flow (3-Step Process)

#### Step 1: Generate OAuth URL
- **Endpoint**: `GET /api/instagram/connect` (protected)
- **Process**:
  1. User authenticated via JWT
  2. Generate short-lived state token (10 min expiry) with user ID embedded
  3. Build Facebook OAuth URL with scopes:
     - `instagram_basic`
     - `pages_show_list`
     - `pages_read_engagement`
     - `business_management`
  4. Return OAuth URL to frontend
  5. Frontend redirects user to Facebook login

#### Step 2: OAuth Callback
- **Endpoint**: `GET /auth/instagram/callback`
- **Process**:
  1. Instagram redirects with `code` and `state` query params
  2. Validate state token (prevents CSRF attacks)
  3. Exchange authorization code for short-lived token:
     - `POST https://api.instagram.com/oauth/access_token`
     - Returns short-lived token (1 hour)
  4. Exchange short-lived token for long-lived token:
     - `GET https://graph.instagram.com/access_token?grant_type=ig_exchange_token`
     - Returns 60-day token
  5. Fetch Instagram user profile (`/me?fields=id,username`)
  6. Save account to database with token and expiry
  7. Trigger background post fetch (non-blocking goroutine)
  8. Redirect to frontend with success flag

#### Step 3: Token Refresh (Automated)
- **Background Job**: Runs every 12 hours
- **Process**:
  1. Query accounts with tokens expiring within 7 days
  2. For each account, call refresh endpoint:
     - `GET https://graph.instagram.com/refresh_access_token?grant_type=ig_refresh_token`
  3. Update database with new token and expiry
  4. Log success/failure for monitoring

### Security Measures
- **State Token**: Prevents CSRF attacks, short-lived (10 min)
- **Token Storage**: Long-lived tokens stored securely in database
- **Automatic Refresh**: Prevents token expiration issues
- **Scopes**: Minimal required permissions

---

## 🔑 JWT Authentication Strategy

### Token Generation
- **Library**: `golang-jwt/jwt/v5`
- **Algorithm**: HS256 (HMAC-SHA256)
- **Claims Structure**:
  ```go
  type Claims struct {
      UserID string `json:"user_id"`
      Email  string `json:"email"`
      jwt.RegisteredClaims
  }
  ```
- **Expiry**: 24 hours for auth tokens, 10 minutes for OAuth state tokens

### Token Flow

#### Registration/Login
1. User provides email/password
2. Backend validates credentials
3. Generate JWT with user ID and email
4. Return token to frontend
5. Frontend stores in `localStorage`

#### Protected Routes
1. Frontend sends `Authorization: Bearer <token>` header
2. `AuthMiddleware` extracts token from header
3. Validates token signature and expiry
4. Extracts claims (user_id, email)
5. Stores in Gin context for handler access
6. Handler accesses via `c.GetString("user_id")`

### Security Features
- **Password Hashing**: bcrypt with default cost (10 rounds)
- **Token Validation**: Signature verification + expiry check
- **Rate Limiting**: 5 login attempts per minute per IP
- **No Password in Responses**: Passwords excluded from JSON (`json:"-"`)

### Middleware Implementation
```go
func AuthMiddleware() gin.HandlerFunc {
    // Extract Bearer token from header
    // Validate token signature
    // Set user_id and email in context
    // Abort if invalid
}
```

---

## 📱 Instagram API Integration

### API Endpoints Used

#### 1. OAuth Token Exchange
- `POST https://api.instagram.com/oauth/access_token`
- Exchanges auth code for short-lived token

#### 2. Long-Lived Token Exchange
- `GET https://graph.instagram.com/access_token?grant_type=ig_exchange_token`
- Converts short-lived to 60-day token

#### 3. Token Refresh
- `GET https://graph.instagram.com/refresh_access_token?grant_type=ig_refresh_token`
- Refreshes long-lived token (extends 60 days)

#### 4. User Profile
- `GET https://graph.instagram.com/me?fields=id,username&access_token=...`
- Fetches authenticated user's profile

#### 5. User Media
- `GET https://graph.instagram.com/me/media?fields=id,caption,media_type,media_url,timestamp,like_count,comments_count&limit=50`
- Fetches recent posts with engagement metrics

### Data Fetching Strategy

#### Initial Fetch
- Triggered automatically after OAuth connection
- Runs in background goroutine (non-blocking)
- Fetches last 50 posts
- Upserts to database (prevents duplicates)

#### Manual Refresh
- **Endpoint**: `POST /api/instagram/refresh` (protected)
- User-triggered refresh
- Also runs in background goroutine
- Returns immediately with "fetch scheduled" status

#### Post Storage
- **Upsert Strategy**: `ON CONFLICT (account_id, ig_post_id) DO UPDATE`
- Updates existing posts with latest engagement metrics
- Handles timestamp parsing (multiple formats: RFC3339, Instagram format)
- Continues on individual post errors (doesn't fail entire batch)

### Rate Limiting Considerations
- Instagram API has rate limits (not explicitly handled in code)
- Background fetching reduces API calls
- Batch operations minimize requests
- Token refresh happens proactively (before expiry)

---

## 🤖 AI Recommendation System Architecture

### Implementation
- **Type**: Rule-based AI (not ML model)
- **Location**: `generateGrowthMessage()` in `internal/services/instagram.go`

### Algorithm

#### Input Data
- `LikesTrend`: Percentage change vs previous period
- `CommentsTrend`: Percentage change vs previous period
- `PostsThisWeek`: Posting frequency
- `TotalPosts`: Current period post count

#### Decision Tree
```
IF TotalPosts == 0:
    → "No posts yet in this period. Time to share something! 📸"

ELSE IF LikesTrend > 20:
    → "🔥 You're on fire! Engagement is way up."

ELSE IF LikesTrend > 5:
    → "📈 Nice! You're growing steadily."

ELSE IF LikesTrend > -5:
    → "😎 Holding steady - keep doing your thing."

ELSE IF LikesTrend > -20:
    → "📉 Slight dip, but no worries - it happens."

ELSE:
    → "💪 Engagement is down, but consistency is key!"

// Add posting frequency context
IF PostsThisWeek == 0:
    → Append " Haven't posted this week though - your audience misses you!"

ELSE IF PostsThisWeek >= 5:
    → Append " You've been posting a lot - great hustle!"
```

### Why Rule-Based?
- **Fast**: No model inference overhead
- **Interpretable**: Clear logic for debugging
- **Customizable**: Easy to adjust thresholds
- **No Training Data**: Works immediately

### Future Enhancements
- Could integrate ML model for more sophisticated recommendations
- Could analyze hashtag performance
- Could suggest optimal posting times
- Could provide content type recommendations

---

## 📡 API Endpoints & Data Flow

### Public Endpoints

#### Authentication
- `POST /api/auth/register`
  - **Request**: `{email, password}`
  - **Response**: `{token, user: {id, email, created_at}}`
  - **Validation**: Email format, password min 8 chars

- `POST /api/auth/login` (rate limited: 5/min)
  - **Request**: `{email, password}`
  - **Response**: `{token, token_type: "bearer", expires_in: 86400, user}`
  - **Headers**: Sets `Authorization: Bearer <token>`

#### OAuth Callback
- `GET /auth/instagram/callback?code=...&state=...`
  - Handles Instagram OAuth redirect
  - Exchanges code for tokens
  - Saves account to database
  - Redirects to frontend

#### Waitlist
- `POST /api/waitlist/signup`
  - **Request**: `{email}`
  - **Response**: `{message: "Signed up successfully"}`

#### Admin
- `GET /api/admin/waitlist?password=...&page=1&limit=50`
  - Returns paginated waitlist entries
  - Simple password auth (header or query param)

#### Health Check
- `GET /health`
  - Returns server and database status

### Protected Endpoints (Require JWT)

#### User
- `GET /api/user/me`
  - Returns current authenticated user info

#### Instagram
- `GET /api/instagram/connect`
  - Returns OAuth URL for Instagram connection

- `POST /api/instagram/refresh`
  - Triggers background post fetch
  - Returns `{status: "fetch scheduled", account: {...}}`

- `GET /api/instagram/posts`
  - Returns stored posts for user's Instagram account
  - Response: `{account: {...}, posts_count: N, posts: [...]}`

- `DELETE /api/instagram/disconnect`
  - Removes Instagram account and all posts (CASCADE)

#### Analytics
- `GET /api/growth/stats?period=30`
  - **Query Params**: `period` (7, 14, 30, 90 days)
  - **Response**: Complete growth statistics
  - Calculates trends, best post, posting frequency

### Data Flow Example: User Views Dashboard

1. **Frontend**: `GET /api/growth/stats` with JWT token
2. **Middleware**: Validates JWT, extracts user_id
3. **Handler**: Gets Instagram account for user
4. **Service**: Queries database for posts in period
5. **Service**: Calculates stats (aggregations, trends, best post)
6. **Service**: Generates AI recommendation message
7. **Handler**: Returns JSON response
8. **Frontend**: Displays stats in dashboard

---

## ⚡ Performance Considerations

### Database Optimizations

#### Query Efficiency
- **Aggregations**: Single queries with `COUNT()`, `SUM()`, `COALESCE()`
- **Indexes**: Primary keys, unique constraints, foreign keys auto-indexed
- **Parameterized Queries**: Prevents SQL injection, allows query plan caching
- **Time-Range Queries**: Uses `INTERVAL` for efficient date filtering

#### Upsert Strategy
- `ON CONFLICT ... DO UPDATE` prevents duplicate inserts
- Single query for insert/update (no separate SELECT needed)
- Reduces database round trips

### Concurrency

#### Background Jobs
- Token refresh runs in goroutine (non-blocking)
- Post fetching runs in goroutine after OAuth
- Server continues handling requests during background work

#### Non-Blocking Operations
- Post fetch triggered with `go func()` after OAuth
- Returns immediately to user
- Errors logged but don't block user flow

### API Efficiency

#### Batch Operations
- Fetches 50 posts in single API call
- Upserts all posts in transaction-like flow
- Continues on individual errors (doesn't fail batch)

#### Caching Opportunities (Future)
- Could cache growth stats (invalidate on post refresh)
- Could cache Instagram profile data
- Could use Redis for session management

### Frontend Performance

#### React Optimizations
- Uses React Hooks (lightweight state management)
- Conditional rendering (only loads stats when available)
- LocalStorage for token persistence (no server round trip)

#### Network Efficiency
- Single API call for dashboard stats
- Minimal payload sizes
- JSON responses (efficient parsing)

---

## 🎯 Challenges Solved

### 1. Instagram Token Management
**Challenge**: Instagram tokens expire (60 days), need automatic refresh
**Solution**: 
- Background job runs every 12 hours
- Checks for tokens expiring within 7 days
- Automatically refreshes before expiry
- Updates database with new token

### 2. Duplicate Post Prevention
**Challenge**: Same post fetched multiple times
**Solution**:
- Composite unique constraint: `UNIQUE(account_id, ig_post_id)`
- Upsert query: `ON CONFLICT DO UPDATE`
- Updates existing posts with latest metrics

### 3. OAuth State Token Security
**Challenge**: Prevent CSRF attacks during OAuth flow
**Solution**:
- Generate short-lived JWT state token (10 min)
- Embed user ID in token
- Validate state token on callback
- Prevents unauthorized account linking

### 4. Timestamp Parsing
**Challenge**: Instagram returns timestamps in different formats
**Solution**:
- Try multiple formats: RFC3339, Instagram format (+0000)
- Graceful fallback if parsing fails
- Handles timezone variations

### 5. Background Processing
**Challenge**: Post fetching takes time, shouldn't block user
**Solution**:
- Use goroutines for async processing
- Return immediately with "scheduled" status
- Log errors without blocking user flow

### 6. Rate Limiting
**Challenge**: Prevent brute force login attacks
**Solution**:
- ulule/limiter middleware
- 5 attempts per minute per IP
- Memory-based store (fast, in-memory)

### 7. Database Migrations
**Challenge**: Ensure schema consistency across deployments
**Solution**:
- Automatic migrations on server startup
- `CREATE TABLE IF NOT EXISTS` for idempotency
- No separate migration tool needed

### 8. Error Handling
**Challenge**: Graceful degradation when Instagram API fails
**Solution**:
- Continue processing other posts if one fails
- Log errors for monitoring
- Return user-friendly error messages
- Don't expose internal errors to frontend

### 9. CORS Configuration
**Challenge**: Frontend on different origin (localhost:5173)
**Solution**:
- Custom CORS middleware
- Checks Origin header
- Allows specific origins in dev
- Configurable for production

### 10. Password Security
**Challenge**: Store passwords securely
**Solution**:
- bcrypt hashing (default cost: 10)
- Passwords never returned in JSON
- Validation on registration (min 8 chars)
- Secure comparison for login

---

## 📊 Key Metrics & Statistics

### Scalability
- **Database**: PostgreSQL handles 10,000+ posts per account
- **Concurrent Users**: Goroutines handle multiple simultaneous requests
- **API Limits**: Background jobs minimize Instagram API calls

### Data Volume
- **Posts per Account**: Up to 50 posts fetched per refresh
- **Analytics Period**: 7, 14, 30, 90 days
- **Token Refresh**: Every 12 hours, checks 7-day window

### Response Times
- **Health Check**: < 10ms
- **Growth Stats**: Depends on post count (typically < 100ms)
- **Post Fetch**: Background (non-blocking)

---

## 🔄 System Architecture Flow

### User Registration Flow
```
Frontend → POST /api/auth/register
         → Backend validates email/password
         → Hash password (bcrypt)
         → Insert into users table
         → Generate JWT token
         → Return token + user info
         → Frontend stores token in localStorage
```

### Instagram Connection Flow
```
Frontend → GET /api/instagram/connect (with JWT)
         → Backend generates state token
         → Returns OAuth URL
         → Frontend redirects to Facebook
         → User authorizes
         → Instagram redirects to /auth/instagram/callback
         → Backend exchanges code for tokens
         → Saves account to database
         → Triggers background post fetch (goroutine)
         → Redirects to frontend
```

### Dashboard Stats Flow
```
Frontend → GET /api/growth/stats (with JWT)
         → Middleware validates JWT
         → Handler gets user's Instagram account
         → Service queries posts in period
         → Calculates aggregations (likes, comments, trends)
         → Finds best performing post
         → Generates AI recommendation message
         → Returns JSON response
         → Frontend displays stats
```

---

## 🚀 Deployment & Production

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT signing
- `INSTAGRAM_CLIENT_ID`: Instagram app client ID
- `INSTAGRAM_CLIENT_SECRET`: Instagram app secret
- `INSTAGRAM_REDIRECT_URI`: OAuth callback URL
- `PORT`: Server port (default: 8080)
- `FRONTEND_URL`: Frontend URL for OAuth redirects
- `ADMIN_PASSWORD`: Admin password for waitlist access

### Production Considerations
- **Database**: Use connection pooling (PostgreSQL)
- **Logging**: Structured logging for monitoring
- **Monitoring**: Health check endpoint for uptime monitoring
- **Error Tracking**: Consider Sentry or similar
- **Rate Limiting**: Could use Redis for distributed rate limiting
- **CORS**: Update allowed origins for production domain

---

## 📝 Code Quality & Best Practices

### Go Best Practices
- **Error Handling**: Explicit error returns, no panics
- **Context Usage**: Could use context.Context for cancellation
- **Struct Tags**: JSON tags for API responses
- **Package Organization**: Clear separation (api, services, models, database)

### Security Best Practices
- **Input Validation**: Gin binding for request validation
- **SQL Injection**: Parameterized queries only
- **Password Hashing**: bcrypt with appropriate cost
- **Token Security**: JWT with expiration
- **CORS**: Configurable allowed origins

### Frontend Best Practices
- **Protected Routes**: Route-level authentication check
- **Error Handling**: Try-catch for API calls
- **Loading States**: User feedback during async operations
- **Token Storage**: localStorage (consider httpOnly cookies for production)

---

## 🎓 Interview Talking Points

### Why This Architecture?
- **Separation of Concerns**: API handlers, business logic (services), data access (database)
- **Scalability**: Goroutines for concurrent processing
- **Security**: Multiple layers (JWT, rate limiting, password hashing)
- **Maintainability**: Clear package structure, single responsibility

### Trade-offs Made
- **Rule-based AI**: Fast and interpretable, but less sophisticated than ML
- **Memory Rate Limiting**: Fast but not distributed (Redis would be better for scale)
- **Background Jobs**: Simple goroutines vs. job queue (sufficient for current scale)
- **PostgreSQL**: Relational DB vs. NoSQL (better for analytics queries)

### Future Improvements
- **Caching**: Redis for stats caching
- **Job Queue**: Bull/BullMQ for background jobs
- **ML Recommendations**: Train model on engagement patterns
- **Real-time Updates**: WebSockets for live stats
- **Multi-platform**: Support TikTok, YouTube, etc.

---

## 📚 Key Libraries & Dependencies

### Backend
- `gin-gonic/gin`: HTTP web framework
- `golang-jwt/jwt/v5`: JWT token handling
- `lib/pq`: PostgreSQL driver
- `golang.org/x/crypto/bcrypt`: Password hashing
- `ulule/limiter/v3`: Rate limiting
- `google/uuid`: UUID generation

### Frontend
- `react`: UI framework
- `react-router-dom`: Client-side routing
- `vite`: Build tool and dev server
- `tailwindcss`: Utility-first CSS framework

---

**Last Updated**: Based on codebase analysis
**Project**: Creator Growth Tool
**Author**: Dan Imran


