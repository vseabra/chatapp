import { createContext, useContext, useState, ReactNode, useCallback } from 'react'
import { useAuth } from './AuthContext'

export interface Message {
  id: string
  roomId: string
  userId: string
  userName: string
  text: string
  type: string
  createdAt: string
}

interface MessageContextType {
  messages: Message[]
  isLoading: boolean
  error: string | null
  loadMessages: (roomId: string, limit?: number) => Promise<void>
  clearMessages: () => void
}

const MessageContext = createContext<MessageContextType | undefined>(undefined)

export function useMessages() {
  const context = useContext(MessageContext)
  if (context === undefined) {
    throw new Error('useMessages must be used within a MessageProvider')
  }
  return context
}

interface MessageProviderProps {
  children: ReactNode
}

const getApiBaseUrl = () => {
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return 'http://localhost:8080/api/v1/rooms'
  } else {
    return 'http://chat-api:8080/api/v1/rooms'
  }
}

const API_BASE_URL = getApiBaseUrl()

export function MessageProvider({ children }: MessageProviderProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { user } = useAuth()

  const getAuthHeaders = () => {
    const token = localStorage.getItem('accessToken')
    return {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {})
    }
  }

  const loadMessages = useCallback(async (roomId: string, limit: number = 50) => {
    if (!user) return

    setIsLoading(true)
    setError(null)

    try {
      console.log('Loading messages for room:', roomId, 'limit:', limit)

      const url = `${API_BASE_URL}/${roomId}/messages?limit=${limit}`
      console.log('Messages API URL:', url)

      const response = await fetch(url, {
        method: 'GET',
        headers: getAuthHeaders(),
      })

      console.log('Messages response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Failed to load messages'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      console.log('Messages data:', data)

      // Sort messages by createdAt timestamp (oldest first for proper ordering)
      const sortedMessages = (data.items || []).sort((a: Message, b: Message) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      )

      setMessages(sortedMessages)
    } catch (error) {
      console.error('Failed to load messages:', error)
      setError('Failed to load messages')
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [user])

  const clearMessages = useCallback(() => {
    setMessages([])
    setError(null)
  }, [])

  const value: MessageContextType = {
    messages,
    isLoading,
    error,
    loadMessages,
    clearMessages
  }

  return (
    <MessageContext.Provider value={value}>
      {children}
    </MessageContext.Provider>
  )
}
