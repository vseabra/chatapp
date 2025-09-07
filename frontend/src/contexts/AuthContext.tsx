import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'

interface User {
  id: string
  name: string
  email: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
}

// API configuration - dynamically determine URL based on environment
const getApiBaseUrl = () => {
  // Check if we're running in Docker (frontend container)
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    // Running locally - use localhost
    return 'http://localhost:8080/api/v1/auth'
  } else {
    // Running in Docker - use service name
    return 'http://chat-api:8080/api/v1/auth'
  }
}

const API_BASE_URL = getApiBaseUrl()

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Check for existing session on mount
  useEffect(() => {
    const checkAuth = async () => {
      try {
        // TODO: Check if user is already logged in (e.g., check localStorage or make API call)
        const storedUser = localStorage.getItem('user')
        if (storedUser) {
          setUser(JSON.parse(storedUser))
        }
      } catch (error) {
        console.error('Auth check failed:', error)
      } finally {
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [])

  const login = async (email: string, password: string) => {
    setIsLoading(true)
    try {
      console.log('Making login request to:', `${API_BASE_URL}/login`)

      const response = await fetch(`${API_BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      console.log('Login response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Login failed'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      console.log('Login response data:', data)

      // Store the JWT token
      localStorage.setItem('accessToken', data.accessToken)

      // Use the actual user data from the login response
      const user: User = {
        id: data.userId,
        name: data.userName,
        email: email // Email comes from the login form, not the response
      }

      setUser(user)
      localStorage.setItem('user', JSON.stringify(user))
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const register = async (name: string, email: string, password: string) => {
    setIsLoading(true)
    try {
      console.log('Making register request to:', `${API_BASE_URL}/register`)

      const response = await fetch(`${API_BASE_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      })

      console.log('Register response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Registration failed'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      console.log('Register response data:', data)

      // After successful registration, automatically log the user in
      console.log('Making auto-login request to:', `${API_BASE_URL}/login`)

      const loginResponse = await fetch(`${API_BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      console.log('Auto-login response status:', loginResponse.status)

      if (!loginResponse.ok) {
        const loginError = await loginResponse.json()
        throw new Error(loginError.error || 'Auto-login after registration failed')
      }

      const loginData = await loginResponse.json()
      console.log('Auto-login response data:', loginData)

      // Store the JWT token
      localStorage.setItem('accessToken', loginData.accessToken)

      // Use the actual user data from the login response (more reliable and consistent)
      const user: User = {
        id: loginData.userId,
        name: loginData.userName,
        email: email // Email comes from the registration form
      }

      setUser(user)
      localStorage.setItem('user', JSON.stringify(user))
    } catch (error) {
      console.error('Registration failed:', error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const logout = () => {
    setUser(null)
    localStorage.removeItem('user')
    localStorage.removeItem('accessToken')
  }

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}
