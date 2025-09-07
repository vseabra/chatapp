import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/contexts/AuthContext'
import { ChatRoomProvider } from '@/contexts/ChatRoomContext'
import { MessageProvider } from '@/contexts/MessageContext'
import { LoginForm } from '@/components/LoginForm'
import { RegisterForm } from '@/components/RegisterForm'
import { Dashboard } from '@/components/Dashboard'
import { ChatRooms } from '@/components/ChatRooms'
import { ChatRoom } from '@/components/ChatRoom'
import './App.css'

// Protected Route component
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

// Public Route component (redirects to dashboard if already authenticated)
function PublicRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

  return isAuthenticated ? <Navigate to="/" replace /> : <>{children}</>
}

function AppRoutes() {
  return (
    <Routes>
      <Route
        path="/login"
        element={
          <PublicRoute>
            <LoginForm />
          </PublicRoute>
        }
      />
      <Route
        path="/register"
        element={
          <PublicRoute>
            <RegisterForm />
          </PublicRoute>
        }
      />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/chatrooms"
        element={
          <ProtectedRoute>
            <ChatRooms />
          </ProtectedRoute>
        }
      />
      <Route
        path="/chat/:roomId"
        element={
          <ProtectedRoute>
            <ChatRoom />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <AuthProvider>
      <ChatRoomProvider>
        <MessageProvider>
          <Router>
            <AppRoutes />
          </Router>
        </MessageProvider>
      </ChatRoomProvider>
    </AuthProvider>
  )
}

export default App
