'use client'

import { useState, useEffect, Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:3000')

function DashboardContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') || '' : ''
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!token) {
      router.push('/login')
      return
    }
    
    // Check if user just connected Instagram
    if (searchParams.get('connected') === 'true') {
      window.history.replaceState({}, '', window.location.pathname)
    }
    
    // Fetch stats once after handling query parameters
    fetchStats()
  }, [token, router, searchParams])

  const fetchStats = async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/api/growth/stats`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const data = await res.json()
      if (data.stats) {
        setStats(data)
      } else if (data.error?.includes('no instagram')) {
        setStats({ needsConnect: true })
      }
    } catch (err) {
      console.error(err)
    }
    setLoading(false)
  }

  const refreshPosts = async () => {
    setLoading(true)
    await fetch(`${API_URL}/api/instagram/refresh`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    setTimeout(fetchStats, 3000)
  }

  const connectInstagram = async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/api/instagram/connect`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      const data = await res.json()
      if (data.url) {
        window.location.href = data.url
      }
    } catch (err) {
      console.error(err)
    }
    setLoading(false)
  }

  const disconnectInstagram = async () => {
    if (!window.confirm('Are you sure you want to disconnect your Instagram account? All your stats and data will be removed.')) {
      return
    }

    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/api/instagram/disconnect`, {
        method: 'DELETE',
        headers: { 
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })
      
      // Handle non-JSON responses (like 404 HTML pages)
      let data
      const contentType = res.headers.get('content-type')
      if (contentType && contentType.includes('application/json')) {
        data = await res.json()
      } else {
        const text = await res.text()
        console.error('Non-JSON response:', text)
        throw new Error(`Server returned ${res.status}: ${res.statusText}`)
      }
      
      if (res.ok) {
        setStats({ needsConnect: true })
      } else {
        console.error('Disconnect error:', data)
        alert(data.details ? `${data.error}\n\nDetails: ${data.details}` : data.error || 'Failed to disconnect Instagram account')
      }
    } catch (err: any) {
      console.error('Disconnect error:', err)
      if (err.message.includes('404')) {
        alert('Disconnect endpoint not found. Please deploy the latest code to production.')
      } else {
        alert(`Failed to disconnect Instagram account: ${err.message}`)
      }
    }
    setLoading(false)
  }

  const logout = () => {
    localStorage.removeItem('token')
    router.push('/')
  }

  return (
    <div className="min-h-screen p-6 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white">📈 Creator Growth</h1>
          <button
            onClick={logout}
            className="text-gray-400 hover:text-white transition"
          >
            Logout
          </button>
        </div>

        {loading && !stats && (
          <div className="text-center py-20">
            <div className="text-4xl mb-4">⏳</div>
            <p className="text-gray-400">Loading your stats...</p>
          </div>
        )}

        {stats?.needsConnect && (
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10 text-center">
            <div className="text-6xl mb-4">📸</div>
            <h2 className="text-2xl font-bold text-white mb-2">Connect Instagram</h2>
            <p className="text-gray-400 mb-6">Link your IG to see your growth stats</p>
            <button 
              onClick={connectInstagram}
              disabled={loading}
              className="bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold px-8 py-3 rounded-lg hover:opacity-90 transition disabled:opacity-50"
            >
              {loading ? 'Connecting...' : 'Connect Instagram'}
            </button>
          </div>
        )}

        {stats?.stats && (
          <>
            {/* Message Banner */}
            <div className="bg-gradient-to-r from-purple-600/20 to-pink-600/20 border border-purple-500/30 rounded-2xl p-6 mb-6">
              <p className="text-xl text-white">{stats.stats.message}</p>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <StatCard label="Posts" value={stats.stats.total_posts} icon="📝" />
              <StatCard label="Likes" value={stats.stats.total_likes} icon="❤️" />
              <StatCard label="Comments" value={stats.stats.total_comments} icon="💬" />
              <StatCard label="Engagement" value={stats.stats.total_engagement} icon="🔥" />
            </div>

            {/* Averages */}
            <div className="grid grid-cols-2 gap-4 mb-6">
              <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
                <p className="text-gray-400 text-sm mb-1">Avg Likes/Post</p>
                <p className="text-3xl font-bold text-white">{stats.stats.avg_likes_per_post?.toFixed(1) || 0}</p>
              </div>
              <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
                <p className="text-gray-400 text-sm mb-1">Avg Comments/Post</p>
                <p className="text-3xl font-bold text-white">{stats.stats.avg_comments_per_post?.toFixed(1) || 0}</p>
              </div>
            </div>

            {/* Best Post */}
            {stats.stats.best_post && (
              <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10 mb-6">
                <div className="flex items-center gap-2 mb-4">
                  <span className="text-2xl">🏆</span>
                  <h3 className="text-lg font-semibold text-white">Best Performing Post</h3>
                </div>
                <p className="text-gray-300 mb-4">"{stats.stats.best_post.caption}"</p>
                <div className="flex gap-4">
                  <span className="text-pink-400">❤️ {stats.stats.best_post.like_count}</span>
                  <span className="text-blue-400">💬 {stats.stats.best_post.comment_count}</span>
                </div>
              </div>
            )}

            {/* Trends */}
            <div className="grid grid-cols-3 gap-4 mb-6">
              <TrendCard label="Likes Trend" value={stats.stats.likes_trend} />
              <TrendCard label="Comments Trend" value={stats.stats.comments_trend} />
              <TrendCard label="Posting Trend" value={stats.stats.posting_trend} />
            </div>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-3">
              <button
                onClick={refreshPosts}
                disabled={loading}
                className="flex-1 bg-white/10 border border-white/20 text-white font-semibold py-3 rounded-lg hover:bg-white/20 transition disabled:opacity-50"
              >
                {loading ? 'Refreshing...' : '🔄 Refresh Stats'}
              </button>
              <button
                onClick={disconnectInstagram}
                disabled={loading}
                className="flex-1 bg-red-600/20 border border-red-500/30 text-red-400 font-semibold py-3 rounded-lg hover:bg-red-600/30 transition disabled:opacity-50"
              >
                {loading ? 'Disconnecting...' : '🔌 Disconnect Instagram'}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default function Dashboard() {
  return (
    <Suspense fallback={
      <div className="min-h-screen p-6 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
        <div className="max-w-4xl mx-auto">
          <div className="text-center py-20">
            <div className="text-4xl mb-4">⏳</div>
            <p className="text-gray-400">Loading...</p>
          </div>
        </div>
      </div>
    }>
      <DashboardContent />
    </Suspense>
  )
}

function StatCard({ label, value, icon }: { label: string, value: number, icon: string }) {
  return (
    <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10 text-center">
      <div className="text-2xl mb-2">{icon}</div>
      <p className="text-3xl font-bold text-white">{value}</p>
      <p className="text-gray-400 text-sm">{label}</p>
    </div>
  )
}

function TrendCard({ label, value }: { label: string, value: number }) {
  const isPositive = value > 0
  const isNeutral = value === 0
  
  return (
    <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-4 border border-white/10 text-center">
      <p className="text-gray-400 text-xs mb-1">{label}</p>
      <p className={`text-xl font-bold ${isNeutral ? 'text-gray-400' : isPositive ? 'text-green-400' : 'text-red-400'}`}>
        {isNeutral ? '—' : isPositive ? `+${value.toFixed(0)}%` : `${value.toFixed(0)}%`}
      </p>
    </div>
  )
}

