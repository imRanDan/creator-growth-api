'use client'

import { useState, useEffect } from 'react'

const API_URL = process.env.NEXT_PUBLIC_API_URL || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:3000')

export default function Admin() {
  const [password, setPassword] = useState('')
  const [authenticated, setAuthenticated] = useState(false)
  const [entries, setEntries] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(0)
  const [limit] = useState(50)

  useEffect(() => {
    // Check if already authenticated (stored in sessionStorage)
    const storedAuth = typeof window !== 'undefined' ? sessionStorage.getItem('admin_authenticated') : null
    if (storedAuth === 'true') {
      setAuthenticated(true)
      fetchEntries(1)
    }
  }, [])

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setAuthenticated(true)
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('admin_authenticated', 'true')
    }
    fetchEntries(1)
  }

  const fetchEntries = async (pageNum: number) => {
    setLoading(true)
    setError('')
    
    try {
      const storedAuth = typeof window !== 'undefined' ? sessionStorage.getItem('admin_authenticated') : null
      if (storedAuth !== 'true') {
        setAuthenticated(false)
        return
      }

      // Get password from localStorage or prompt
      const adminPassword = typeof window !== 'undefined' ? (localStorage.getItem('admin_password') || password) : password
      
      const res = await fetch(
        `${API_URL}/api/admin/waitlist?page=${pageNum}&limit=${limit}`,
        {
          headers: {
            'X-Admin-Password': adminPassword,
          },
        }
      )

      if (res.status === 401) {
        setAuthenticated(false)
        if (typeof window !== 'undefined') {
          sessionStorage.removeItem('admin_authenticated')
        }
        setError('Invalid password')
        return
      }

      if (!res.ok) {
        throw new Error('Failed to fetch entries')
      }

      const data = await res.json()
      setEntries(data.entries || [])
      setTotal(data.total || 0)
      setPages(data.pages || 0)
      setPage(pageNum)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    setAuthenticated(false)
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('admin_authenticated')
      localStorage.removeItem('admin_password')
    }
    setPassword('')
    setEntries([])
  }

  const exportCSV = () => {
    const headers = ['Email', 'Signed Up']
    const rows = entries.map(entry => [
      entry.email,
      new Date(entry.created_at).toLocaleString()
    ])
    
    const csv = [
      headers.join(','),
      ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
    ].join('\n')
    
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `waitlist-${new Date().toISOString().split('T')[0]}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
  }

  if (!authenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-8 border border-white/10 max-w-md w-full">
          <h1 className="text-3xl font-bold text-white mb-6 text-center">🔐 Admin Login</h1>
          <form onSubmit={handleLogin}>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Admin Password"
              className="w-full bg-white/10 border border-white/20 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 mb-4"
              required
            />
            {error && (
              <p className="text-red-400 text-sm mb-4">{error}</p>
            )}
            <button
              type="submit"
              className="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold px-8 py-3 rounded-lg hover:opacity-90 transition"
            >
              Login
            </button>
          </form>
          <p className="text-gray-400 text-sm mt-4 text-center">
            Set ADMIN_PASSWORD in your .env file
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen p-6 bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">📊 Waitlist Admin</h1>
            <p className="text-gray-400">Total signups: {total}</p>
          </div>
          <div className="flex gap-3">
            <button
              onClick={exportCSV}
              className="bg-green-600 text-white font-semibold px-6 py-2 rounded-lg hover:opacity-90 transition"
            >
              📥 Export CSV
            </button>
            <button
              onClick={handleLogout}
              className="bg-red-600 text-white font-semibold px-6 py-2 rounded-lg hover:opacity-90 transition"
            >
              Logout
            </button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-3 gap-4 mb-6">
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10">
            <p className="text-gray-400 text-sm mb-1">Total Signups</p>
            <p className="text-3xl font-bold text-white">{total}</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10">
            <p className="text-gray-400 text-sm mb-1">This Page</p>
            <p className="text-3xl font-bold text-white">{entries.length}</p>
          </div>
          <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 border border-white/10">
            <p className="text-gray-400 text-sm mb-1">Page {page} of {pages}</p>
            <p className="text-3xl font-bold text-white">{pages}</p>
          </div>
        </div>

        {/* Entries Table */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl border border-white/10 overflow-hidden">
          {loading ? (
            <div className="p-12 text-center">
              <div className="text-4xl mb-4">⏳</div>
              <p className="text-gray-400">Loading entries...</p>
            </div>
          ) : entries.length === 0 ? (
            <div className="p-12 text-center">
              <p className="text-gray-400">No entries found</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-white/5">
                    <tr>
                      <th className="text-left p-4 text-white font-semibold">Email</th>
                      <th className="text-left p-4 text-white font-semibold">Signed Up</th>
                    </tr>
                  </thead>
                  <tbody>
                    {entries.map((entry, idx) => (
                      <tr
                        key={entry.id}
                        className={idx % 2 === 0 ? 'bg-white/5' : 'bg-white/0'}
                      >
                        <td className="p-4 text-gray-300">{entry.email}</td>
                        <td className="p-4 text-gray-400">
                          {new Date(entry.created_at).toLocaleString()}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {pages > 1 && (
                <div className="p-4 border-t border-white/10 flex justify-between items-center">
                  <button
                    onClick={() => fetchEntries(page - 1)}
                    disabled={page === 1}
                    className="bg-white/10 border border-white/20 text-white font-semibold px-4 py-2 rounded-lg hover:opacity-90 transition disabled:opacity-50"
                  >
                    ← Previous
                  </button>
                  <span className="text-gray-400">
                    Page {page} of {pages}
                  </span>
                  <button
                    onClick={() => fetchEntries(page + 1)}
                    disabled={page >= pages}
                    className="bg-white/10 border border-white/20 text-white font-semibold px-4 py-2 rounded-lg hover:opacity-90 transition disabled:opacity-50"
                  >
                    Next →
                  </button>
                </div>
              )}
            </>
          )}
        </div>

        {error && (
          <div className="mt-4 p-4 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-center">
            {error}
          </div>
        )}
      </div>
    </div>
  )
}

