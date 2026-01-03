'use client'

import { useState } from 'react'
import Link from 'next/link'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:3000')

export default function Home() {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
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
        <div className="text-center mb-8 md:mb-12">
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-white mb-3 md:mb-4">
            📈 Creator Growth
          </h1>
          <p className="text-lg sm:text-xl text-gray-300 mb-2">
            Track your Instagram like a casual boss
          </p>
          <p className="text-sm sm:text-base text-gray-400 px-4">
            Analytics for creators who post somewhat often, not every day
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8 md:mb-12">
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-4 sm:p-6 border border-white/10 text-center">
            <div className="text-2xl sm:text-3xl mb-2">📊</div>
            <h3 className="text-white font-semibold mb-2 text-sm sm:text-base">Real-time Analytics</h3>
            <p className="text-gray-400 text-xs sm:text-sm">Track engagement, growth trends, and top posts</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-4 sm:p-6 border border-white/10 text-center">
            <div className="text-2xl sm:text-3xl mb-2">🚀</div>
            <h3 className="text-white font-semibold mb-2 text-sm sm:text-base">Smart Insights</h3>
            <p className="text-gray-400 text-xs sm:text-sm">Get actionable recommendations to grow</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-4 sm:p-6 border border-white/10 text-center">
            <div className="text-2xl sm:text-3xl mb-2">⚡</div>
            <h3 className="text-white font-semibold mb-2 text-sm sm:text-base">Easy Setup</h3>
            <p className="text-gray-400 text-xs sm:text-sm">Connect Instagram in seconds</p>
          </div>
        </div>

        {/* Waitlist Form */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-4 sm:p-6 md:p-8 border border-white/10">
          {success ? (
            <div className="text-center">
              <div className="text-4xl sm:text-5xl mb-4">🎉</div>
              <h2 className="text-xl sm:text-2xl font-bold text-white mb-2">You're in!</h2>
              <p className="text-sm sm:text-base text-gray-300 mb-6">
                We'll notify you when we launch. Thanks for joining the waitlist!
              </p>
              <button
                onClick={() => {
                  setSuccess(false)
                  setEmail('')
                }}
                className="text-purple-400 hover:text-purple-300 font-semibold text-sm sm:text-base"
              >
                Add another email
              </button>
            </div>
          ) : (
            <>
              <h2 className="text-2xl sm:text-3xl font-bold text-white mb-3 sm:mb-4 text-center">
                Sign up on the waitlist today!
              </h2>
              <p className="text-sm sm:text-base text-gray-400 text-center mb-4 sm:mb-6">
                Be the first to know when we launch
              </p>
              <form onSubmit={handleSubmit}>
                <div className="flex flex-col sm:flex-row gap-3">
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="your@email.com"
                    className="flex-1 bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 text-base"
                    required
                  />
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full sm:w-auto bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold px-6 sm:px-8 py-3 rounded-lg hover:opacity-90 transition disabled:opacity-50 whitespace-nowrap"
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
        <div className="text-center mt-6 sm:mt-8">
          <p className="text-gray-400 mb-3 sm:mb-4 text-sm sm:text-base">Already have access?</p>
          <Link
            href="/login"
            className="text-purple-400 hover:text-purple-300 font-semibold underline text-sm sm:text-base"
          >
            Sign in to your account →
          </Link>
        </div>
      </div>
    </div>
  )
}

