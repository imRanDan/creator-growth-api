import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Analytics } from '@vercel/analytics/react'
import Waitlist from './Waitlist'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import Admin from './pages/Admin'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/waitlist" element={<Waitlist />} />
        <Route path="/" element={<Waitlist />} />
        <Route path="/login" element={<Login />} />
        <Route 
          path="/dashboard" 
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          } 
        />
        <Route path="/admin" element={<Admin />} />
        <Route path="*" element={<Navigate to="/waitlist" replace />} />
      </Routes>
      <Analytics />
    </BrowserRouter>
  )
}

// Protected route component
function ProtectedRoute({ children }) {
  const token = localStorage.getItem('token')
  return token ? children : <Navigate to="/login" replace />
}

export default App
