import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const API_URL = 'https://creator-growth-api-production.up.railway.app'

function Login() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [isSignup, setIsSignup] = useState(false)

  const signup = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/api/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      })
      const data = await res.json()
      if (data.token) {
        localStorage.setItem('token', data.token)
        navigate('/dashboard')
      } else {
        setError(data.error || 'Signup failed')
      }
    } catch (err) {
      setError('Connection error')
    }
    setLoading(false)
  }

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
        localStorage.setItem('token', data.token)
        navigate('/dashboard')
      } else {
        setError(data.error || 'Login failed')
      }
    } catch (err) {
      setError('Connection error')
    }
    setLoading(false)
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">📈 Creator Growth</h1>
          <p className="text-gray-400">Track your Instagram like a casual boss</p>
        </div>
        
        <form onSubmit={isSignup ? signup : login} className="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10">
          <div className="mb-4">
            <label className="block text-gray-300 text-sm mb-2">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
              placeholder="you@email.com"
              required
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
              required
            />
          </div>
          {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold py-3 rounded-lg hover:opacity-90 transition disabled:opacity-50"
          >
            {loading ? 'Loading...' : isSignup ? 'Sign Up' : 'Sign In'}
          </button>
          <p className="text-center text-gray-400 text-sm mt-4">
            {isSignup ? 'Already have an account?' : "Don't have an account?"}{' '}
            <button
              type="button"
              onClick={() => {
                setIsSignup(!isSignup)
                setError('')
              }}
              className="text-purple-400 hover:text-purple-300 font-semibold"
            >
              {isSignup ? 'Sign In' : 'Sign Up'}
            </button>
          </p>
          <p className="text-center text-gray-500 text-xs mt-4">
            <button
              type="button"
              onClick={() => navigate('/waitlist')}
              className="text-purple-400 hover:text-purple-300"
            >
              ← Back to waitlist
            </button>
          </p>
        </form>
      </div>
    </div>
  )
}

export default Login


