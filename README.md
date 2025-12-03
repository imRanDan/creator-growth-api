# 📈 Creator Growth Tool

A full-stack Instagram analytics platform that helps content creators track engagement metrics, analyze growth trends, and get AI-powered recommendations to optimize their content strategy.

![Go](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## 🎯 Overview

Creator Growth Tool connects to Instagram via OAuth 2.0 and analyzes up to 10,000+ posts to provide actionable insights. It calculates engagement rates, tracks posting trends, identifies top-performing content, and delivers personalized growth recommendations.

Built as a scalable SaaS solution with production-grade security (JWT authentication, rate limiting, secure token storage) and automated token refresh for long-lived Instagram access.

## ✨ Features

### 📊 Analytics Dashboard
- **Real-time Metrics**: Track total posts, likes, comments, and engagement
- **Trend Analysis**: Compare current performance vs. previous periods with percentage changes
- **Best Post Identification**: Automatically identifies your top-performing content
- **Engagement Rates**: Calculate average likes and comments per post

### 🔐 Authentication & Security
- JWT-based authentication with secure token storage
- Instagram OAuth 2.0 integration
- Rate limiting (5 attempts/minute) to prevent brute force attacks
- Long-lived token management with automatic refresh (60-day tokens)
- Password hashing with bcrypt

### 📈 Growth Intelligence
- AI-powered growth messages based on engagement trends
- Posting frequency tracking (weekly/monthly)
- Period-based comparison analytics (30/60/90 day views)
- Hashtag extraction and analysis

### 🔄 Data Management
- Automated post synchronization from Instagram
- Background jobs for token refresh
- Database-backed analytics with PostgreSQL
- Efficient data upserts to prevent duplicates

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.25.3
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **OAuth**: Instagram Graph API
- **Security**: bcrypt password hashing, rate limiting (ulule/limiter)
- **Deployment**: Railway (Production)

### Frontend
- **Framework**: React 19
- **Build Tool**: Vite
- **Styling**: Tailwind CSS 4
- **State Management**: React Hooks
- **HTTP Client**: Fetch API

### DevOps & Infrastructure
- PostgreSQL database with migrations
- Environment-based configuration
- CORS middleware for cross-origin requests
- Health check endpoints for monitoring

## 📁 Project Structure

```
creator-growth-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── auth.go             # Authentication handlers
│   │   ├── instagram.go        # Instagram OAuth & webhooks
│   │   └── middleware.go       # JWT middleware
│   ├── database/
│   │   ├── db.go               # Database connection
│   │   └── migrations.go       # Schema migrations
│   ├── models/
│   │   └── user.go             # User data models
│   └── services/
│       ├── auth.go             # Auth business logic
│       ├── instagram.go        # Instagram API integration
│       ├── jobs.go             # Background jobs
│       └── jwt.go              # Token generation/validation
├── frontend/
│   ├── src/
│   │   ├── App.jsx             # Main React component
│   │   ├── main.jsx            # React entry point
│   │   └── index.css           # Global styles
│   ├── package.json
│   └── vite.config.js
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Getting Started

### Prerequisites
- Go 1.25+ installed
- PostgreSQL 14+ running
- Node.js 18+ and npm
- Instagram Developer Account with approved app

### Backend Setup

1. **Clone the repository**
```bash
git clone https://github.com/imRanDan/creator-growth-api.git
cd creator-growth-api
```

2. **Set up environment variables**
```bash
# Create .env file in root directory
cat > .env << EOF
# Database
DATABASE_URL=postgres://user:password@localhost:5432/creator_growth?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Instagram OAuth
INSTAGRAM_CLIENT_ID=your-instagram-client-id
INSTAGRAM_CLIENT_SECRET=your-instagram-client-secret
INSTAGRAM_REDIRECT_URI=http://localhost:8080/auth/instagram/callback

# Server
PORT=8080
EOF
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run database migrations**
```bash
go run cmd/server/main.go
# Migrations run automatically on startup
```

5. **Start the backend server**
```bash
go run cmd/server/main.go
# Server starts on http://localhost:8080
```

### Frontend Setup

1. **Navigate to frontend directory**
```bash
cd frontend
```

2. **Install dependencies**
```bash
npm install
```

3. **Update API URL** (optional for local dev)
```javascript
// src/App.jsx - change for local development
const API_URL = 'http://localhost:8080'
```

4. **Start development server**
```bash
npm run dev
# Frontend runs on http://localhost:5173
```

## 📡 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Create new user account | ❌ |
| POST | `/api/auth/login` | Login with credentials | ❌ |
| GET | `/api/user/me` | Get current user info | ✅ |

### Instagram Integration
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/instagram/connect` | Get OAuth URL for Instagram | ✅ |
| GET | `/auth/instagram/callback` | OAuth callback handler | ❌ |
| POST | `/api/instagram/refresh` | Manually refresh posts from IG | ✅ |
| GET | `/api/instagram/posts` | Get user's Instagram posts | ✅ |

### Analytics
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/growth/stats` | Get engagement analytics | ✅ |

### System
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check endpoint | ❌ |

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Instagram Accounts Table
```sql
CREATE TABLE instagram_accounts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    ig_user_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255),
    access_token TEXT NOT NULL,
    token_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Instagram Posts Table
```sql
CREATE TABLE instagram_posts (
    id UUID PRIMARY KEY,
    ig_post_id VARCHAR(255) NOT NULL,
    account_id UUID REFERENCES instagram_accounts(id),
    caption TEXT,
    media_type VARCHAR(50),
    media_url TEXT,
    like_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    posted_at TIMESTAMP,
    fetched_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(account_id, ig_post_id)
);
```

## 🔒 Security Features

- **Password Security**: Bcrypt hashing with salt rounds
- **JWT Tokens**: Secure token-based authentication
- **Rate Limiting**: 5 login attempts per minute per IP
- **OAuth 2.0**: Industry-standard Instagram authentication
- **Long-lived Tokens**: 60-day Instagram access tokens with auto-refresh
- **CORS Protection**: Configurable allowed origins
- **SQL Injection Prevention**: Parameterized queries
- **Environment Variables**: Sensitive data stored in .env files

## 🎨 Frontend Features

### Responsive Design
- Mobile-first approach with Tailwind CSS
- Glassmorphism UI elements
- Gradient accents and modern color palette
- Emoji-driven visual hierarchy

### User Experience
- Loading states for async operations
- Error handling with user-friendly messages
- Persistent authentication (localStorage)
- Smooth transitions and hover effects
- Real-time stat updates

## 📊 Sample Analytics Output

```json
{
  "stats": {
    "total_posts": 45,
    "total_likes": 12500,
    "total_comments": 890,
    "total_engagement": 13390,
    "avg_likes_per_post": 277.8,
    "avg_comments_per_post": 19.8,
    "engagement_rate": 297.6,
    "likes_trend": 15.3,
    "comments_trend": 8.7,
    "posting_trend": -5.2,
    "posts_this_week": 3,
    "posts_this_month": 12,
    "period_days": 30,
    "message": "📈 Nice! You're growing steadily.",
    "best_post": {
      "caption": "Just launched my new project! Link in bio 🚀",
      "like_count": 856,
      "comment_count": 47,
      "engagement": 903,
      "posted_at": "2024-11-15T10:30:00Z"
    }
  }
}
```

## 🚧 Future Enhancements

- [ ] Multi-platform support (TikTok, YouTube, Twitter)
- [ ] Advanced analytics (demographics, best posting times)
- [ ] Content scheduling and calendar
- [ ] Competitor analysis
- [ ] Export reports to PDF
- [ ] Email notifications for milestones
- [ ] Team collaboration features
- [ ] Custom dashboard widgets
- [ ] AI-powered caption suggestions
- [ ] Hashtag performance tracking

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Dan Imran**
- GitHub: [@imRanDan](https://github.com/imRanDan)
- LinkedIn: [Danyal Imran](https://linkedin.com/in/danyalimran)

## 🙏 Acknowledgments

- Instagram Graph API for providing the data access
- Gin framework for the excellent Go web toolkit
- React and Vite teams for modern frontend tooling
- Railway for reliable hosting infrastructure

## 📞 Support

For support, email dan.imran97@gmail.com or open an issue in the GitHub repository.

---

**Built with ❤️ by Dan Imran | 2025**

