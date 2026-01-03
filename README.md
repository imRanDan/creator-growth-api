# 📈 Creator Growth Tool

A full-stack Instagram analytics platform built with Next.js that helps content creators track engagement metrics, analyze growth trends, and get AI-powered recommendations to optimize their content strategy.

![Next.js](https://img.shields.io/badge/Next.js-15-black?style=flat&logo=next.js)
![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0-3178C6?style=flat&logo=typescript)
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
- Rate limiting to prevent brute force attacks
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

### Full Stack
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript
- **UI**: React 19 + Tailwind CSS
- **Database**: PostgreSQL (via `pg` / Vercel Postgres)
- **Authentication**: JWT (jsonwebtoken)
- **OAuth**: Instagram Graph API
- **Email**: Resend API
- **Deployment**: Vercel

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ 
- PostgreSQL database (local Docker or Vercel Postgres)
- Instagram App credentials (Facebook Developer Console)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd creator-growth-api
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Set up environment variables**
   
   Create `.env.local` in the root directory:
   ```bash
   # Database
   POSTGRES_URL=postgresql://username:password@localhost:5432/dbname
   
   # JWT Secret (min 32 characters)
   JWT_SECRET=your-super-secret-random-string-here
   
   # Instagram OAuth
   INSTAGRAM_CLIENT_ID=your-instagram-client-id
   INSTAGRAM_CLIENT_SECRET=your-instagram-client-secret
   INSTAGRAM_REDIRECT_URI=http://localhost:3000/auth/instagram/callback
   
   # Frontend URL
   FRONTEND_URL=http://localhost:3000
   NEXT_PUBLIC_API_URL=http://localhost:3000
   
   # Email (Optional - for waitlist notifications)
   RESEND_API_KEY=your-resend-api-key
   
   # Admin (Optional)
   ADMIN_PASSWORD=your-admin-password
   ```

4. **Set up database**
   
   Using Docker:
   ```bash
   docker run --name cg-postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=creator_growth \
     -p 5432:5432 \
     -d postgres:15
   ```

5. **Run migrations**
   ```bash
   npm run migrate
   ```

6. **Start development server**
   ```bash
   npm run dev
   ```

   Visit `http://localhost:3000`

## 📁 Project Structure

```
creator-growth-api/
├── app/                    # Next.js App Router
│   ├── api/               # API routes
│   │   ├── auth/          # Authentication endpoints
│   │   ├── instagram/     # Instagram OAuth & data
│   │   ├── waitlist/      # Waitlist signup
│   │   ├── admin/         # Admin dashboard API
│   │   └── growth/        # Analytics endpoints
│   ├── auth/              # OAuth callback routes
│   ├── dashboard/         # Dashboard page
│   ├── login/             # Login page
│   ├── admin/             # Admin page
│   ├── page.tsx           # Home/waitlist page
│   ├── layout.tsx         # Root layout
│   └── globals.css        # Global styles
├── lib/                   # Shared utilities
│   ├── db.ts             # Database connection
│   ├── services/         # Business logic
│   │   ├── auth.ts       # Authentication
│   │   ├── jwt.ts        # JWT tokens
│   │   ├── instagram.ts  # Instagram API
│   │   ├── email.ts      # Email sending
│   │   └── growth.ts     # Analytics
│   └── utils/            # Helper functions
│       └── auth.ts       # Auth middleware helpers
├── scripts/              # Utility scripts
│   └── migrate.ts        # Database migrations
├── middleware.ts         # Next.js middleware (JWT auth)
├── package.json
├── tsconfig.json
└── .env.local           # Environment variables (not in git)
```

## 🔑 API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user (protected)

### Instagram
- `GET /api/instagram/connect` - Get OAuth URL (protected)
- `GET /auth/instagram/callback` - OAuth callback
- `POST /api/instagram/refresh` - Refresh posts (protected)
- `DELETE /api/instagram/disconnect` - Disconnect account (protected)
- `GET /api/instagram/posts` - Get user's posts (protected)

### Analytics
- `GET /api/growth/stats` - Get growth statistics (protected)

### Waitlist
- `POST /api/waitlist/signup` - Join waitlist

### Admin
- `GET /api/admin/waitlist` - View waitlist entries (admin only)

## 🚢 Deployment

### Vercel (Recommended)

1. **Push to GitHub**
2. **Connect to Vercel**
   - Import your repository
   - Vercel will auto-detect Next.js
3. **Add environment variables** in Vercel dashboard
4. **Set up Vercel Postgres** (optional, or use your own PostgreSQL)
5. **Deploy!**

### Environment Variables for Production

Make sure to set all environment variables in your deployment platform:
- `POSTGRES_URL` - Your production database URL
- `JWT_SECRET` - Strong random secret
- `INSTAGRAM_CLIENT_ID` - Your Instagram app ID
- `INSTAGRAM_CLIENT_SECRET` - Your Instagram app secret
- `INSTAGRAM_REDIRECT_URI` - Your production callback URL (e.g., `https://your-domain.com/auth/instagram/callback`)

## 📝 Database Schema

- **users** - User accounts
- **instagram_accounts** - Connected Instagram accounts
- **instagram_posts** - Fetched Instagram posts
- **waitlist** - Waitlist signups

## 🔒 Security

- JWT tokens for authentication
- Password hashing with bcrypt
- Rate limiting on auth endpoints
- Secure token storage
- CORS protection
- Environment variable security

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions welcome! Please open an issue or submit a pull request.

## 📧 Support

For issues and questions, please open a GitHub issue.
