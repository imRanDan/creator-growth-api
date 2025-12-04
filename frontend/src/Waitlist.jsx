import { useState } from 'react'

const API_URL = 'https://creator-growth-api-production.up.railway.app'

function Waitlist() {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await fetch(`${API_URL}/api/waitlist/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      })

      const data = await res.json()
      
      if (res.ok) {
        setSuccess(true)
        setEmail('')
      } else {
        setError(data.error || 'Something went wrong')
      }
    } catch (err) {
      setError('Connection error. Please try again.')
    }
    
    setLoading(false)
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
      <div className="w-full max-w-2xl">
        {/* Hero Section */}
        <div className="text-center mb-12">
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
            📈 Creator Growth
          </h1>
          <p className="text-xl text-gray-300 mb-2">
            Track your Instagram like a casual boss
          </p>
          <p className="text-gray-400">
            Analytics for creators who post somewhat often, not every day
          </p>
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-3 gap-4 mb-12">
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10 text-center">
            <div className="text-3xl mb-2">📊</div>
            <h3 className="text-white font-semibold mb-2">Real-time Analytics</h3>
            <p className="text-gray-400 text-sm">Track engagement, growth trends, and top posts</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10 text-center">
            <div className="text-3xl mb-2">🚀</div>
            <h3 className="text-white font-semibold mb-2">Smart Insights</h3>
            <p className="text-gray-400 text-sm">Get actionable recommendations to grow</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10 text-center">
            <div className="text-3xl mb-2">⚡</div>
            <h3 className="text-white font-semibold mb-2">Easy Setup</h3>
            <p className="text-gray-400 text-sm">Connect Instagram in seconds</p>
          </div>
        </div>

        {/* Waitlist Form */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10">
          {success ? (
            <div className="text-center">
              <div className="text-5xl mb-4">🎉</div>
              <h2 className="text-2xl font-bold text-white mb-2">You're in!</h2>
              <p className="text-gray-300 mb-6">
                We'll notify you when we launch. Thanks for joining the waitlist!
              </p>
              <button
                onClick={() => {
                  setSuccess(false)
                  setEmail('')
                }}
                className="text-purple-400 hover:text-purple-300 font-semibold"
              >
                Add another email
              </button>
            </div>
          ) : (
            <>
              <h2 className="text-2xl font-bold text-white mb-2 text-center">
                Join the Waitlist
              </h2>
              <p className="text-gray-400 text-center mb-6">
                Be the first to know when we launch
              </p>
              <form onSubmit={handleSubmit}>
                <div className="flex gap-3">
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="your@email.com"
                    className="flex-1 bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
                    required
                  />
                  <button
                    type="submit"
                    disabled={loading}
                    className="bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold px-8 py-3 rounded-lg hover:opacity-90 transition disabled:opacity-50 whitespace-nowrap"
                  >
                    {loading ? 'Adding...' : 'Join Waitlist'}
                  </button>
                </div>
                {error && (
                  <p className="text-red-400 text-sm mt-3 text-center">{error}</p>
                )}
              </form>
            </>
          )}
        </div>

        {/* CTA */}
        <div className="text-center mt-8">
          <p className="text-gray-400 mb-4">Already have access?</p>
          <a
            href="/login"
            className="text-purple-400 hover:text-purple-300 font-semibold underline"
          >
            Sign in to your account →
          </a>
        </div>
      </div>
    </div>
  )
}

export default Waitlist

