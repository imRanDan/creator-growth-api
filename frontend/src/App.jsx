import { useState, useEffect } from 'react'

const API_URL = 'http://localhost:8080'

function App() {
  const [token, setToken] = useState(localStorage.getItem('token') || '')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [user, setUser] = useState(null)

  useEffect(() => {
    if (token) {
      fetchStats()
    }
  }, [token])

  const login = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      })
      const data = await res.json()
      if (data.token) {
        setToken(data.token)
        setUser(data.user)
        localStorage.setItem('token', data.token)
      } else {
        setError(data.error || 'Login failed')
      }
    } catch (err) {
      setError('Connection error')
    }
    setLoading(false)
  }

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

  const logout = () => {
    setToken('')
    setStats(null)
    setUser(null)
    localStorage.removeItem('token')
  }

  // Login Screen
  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-white mb-2">📈 Creator Growth</h1>
            <p className="text-gray-400">Track your Instagram like a casual boss</p>
          </div>
          
          <form onSubmit={login} className="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10">
            <div className="mb-4">
              <label className="block text-gray-300 text-sm mb-2">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
                placeholder="you@email.com"
              />
            </div>
            <div className="mb-6">
              <label className="block text-gray-300 text-sm mb-2">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
                placeholder="••••••••"
              />
            </div>
            {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold py-3 rounded-lg hover:opacity-90 transition disabled:opacity-50"
            >
              {loading ? 'Loading...' : 'Sign In'}
            </button>
          </form>
        </div>
      </div>
    )
  }

  // Dashboard
  return (
    <div className="min-h-screen p-6">
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
            <button className="bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold px-8 py-3 rounded-lg hover:opacity-90 transition">
              Connect Instagram
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

            {/* Refresh Button */}
            <button
              onClick={refreshPosts}
              disabled={loading}
              className="w-full bg-white/10 border border-white/20 text-white font-semibold py-3 rounded-lg hover:bg-white/20 transition disabled:opacity-50"
            >
              {loading ? 'Refreshing...' : '🔄 Refresh Stats'}
            </button>
          </>
        )}
      </div>
    </div>
  )
}

function StatCard({ label, value, icon }) {
  return (
    <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10 text-center">
      <div className="text-2xl mb-2">{icon}</div>
      <p className="text-3xl font-bold text-white">{value}</p>
      <p className="text-gray-400 text-sm">{label}</p>
    </div>
  )
}

function TrendCard({ label, value }) {
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

export default App
