import { createContext, useContext, useState, type ReactNode, useCallback, useRef } from 'react'
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
  isLoadingMore: boolean
  error: string | null
  sendError: string | null
  hasMoreMessages: boolean
  loadMessages: (roomId: string, limit?: number) => Promise<void>
  loadMoreMessages: (roomId: string) => Promise<void>
  sendMessage: (roomId: string, text: string, userId: string, userName: string) => Promise<void>
  registerWebSocket: (roomId: string, ws: WebSocket) => void
  unregisterWebSocket: (roomId: string) => void
  addRealTimeMessage: (message: Message) => void
  clearMessages: () => void
  clearSendError: () => void
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
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [sendError, setSendError] = useState<string | null>(null)
  const [hasMoreMessages, setHasMoreMessages] = useState(false)
  const [nextCursor, setNextCursor] = useState<string | null>(null)
  const [currentRoomId, setCurrentRoomId] = useState<string | null>(null)
  const { user } = useAuth()

  // Store WebSocket connections for sending messages
  const wsConnectionsRef = useRef<Map<string, WebSocket>>(new Map())

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
    setCurrentRoomId(roomId)

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
      setHasMoreMessages(!!data.nextCursor)
      setNextCursor(data.nextCursor || null)
    } catch (error) {
      console.error('Failed to load messages:', error)
      setError('Failed to load messages')
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [user])

  const loadMoreMessages = useCallback(async (roomId: string) => {
    if (!user || !nextCursor || isLoadingMore || currentRoomId !== roomId) return

    setIsLoadingMore(true)

    try {
      console.log('Loading more messages for room:', roomId, 'cursor:', nextCursor)

      const url = `${API_BASE_URL}/${roomId}/messages?limit=50&cursor=${nextCursor}`
      console.log('Load more messages API URL:', url)

      const response = await fetch(url, {
        method: 'GET',
        headers: getAuthHeaders(),
      })

      console.log('Load more response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Failed to load more messages'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      console.log('Load more messages data:', data)

      // Sort messages by createdAt timestamp (oldest first for proper ordering)
      const newMessages = (data.items || []).sort((a: Message, b: Message) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      )

      // Prepend older messages to the existing messages
      setMessages(prevMessages => [...newMessages, ...prevMessages])
      setHasMoreMessages(!!data.nextCursor)
      setNextCursor(data.nextCursor || null)
    } catch (error) {
      console.error('Failed to load more messages:', error)
      setError('Failed to load more messages')
      throw error
    } finally {
      setIsLoadingMore(false)
    }
  }, [user, nextCursor, isLoadingMore, currentRoomId])

  const addRealTimeMessage = useCallback((message: Message) => {
    console.log('Adding real-time message:', message)
    setMessages(prevMessages => {
      // Check if message already exists to prevent duplicates
      const messageExists = prevMessages.some(m => m.id === message.id)
      if (messageExists) {
        console.log('Message already exists, skipping:', message.id)
        return prevMessages
      }

      // Add the new message and sort by timestamp
      const newMessages = [...prevMessages, message].sort((a, b) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      )

      return newMessages
    })
  }, [])

  const sendMessage = useCallback(async (roomId: string, text: string, userId: string, userName: string) => {
    if (!user) {
      throw new Error('User not authenticated')
    }

    setSendError(null)

    // Try to send via WebSocket first
    const ws = wsConnectionsRef.current.get(roomId)
    if (ws && ws.readyState === WebSocket.OPEN) {
      try {
        console.log('Sending message via WebSocket:', { roomId, text, userId, userName })

        const messageData = {
          type: 'submit',
          roomId,
          userId,
          userName,
          text
        }

        ws.send(JSON.stringify(messageData))
        console.log('Message sent via WebSocket')
        return
      } catch (error) {
        console.error('WebSocket send failed, falling back to HTTP:', error)
        // Fall through to HTTP fallback
      }
    }

    // Fallback to HTTP API
    try {
      console.log('Sending message via HTTP fallback:', { roomId, text })

      const response = await fetch(`${API_BASE_URL}/submit`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          roomId,
          userId,
          userName,
          text
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to send message')
      }

      console.log('Message sent via HTTP fallback')
    } catch (error) {
      console.error('HTTP send failed:', error)
      setSendError('Failed to send message. Please try again.')
      throw error
    }
  }, [user])

  const clearSendError = useCallback(() => {
    setSendError(null)
  }, [])

  const clearMessages = useCallback(() => {
    setMessages([])
    setError(null)
    setSendError(null)
    setHasMoreMessages(false)
    setNextCursor(null)
    setCurrentRoomId(null)
  }, [])

  // Function to register WebSocket connection for sending (used internally)
  const registerWebSocketConnection = useCallback((roomId: string, ws: WebSocket) => {
    wsConnectionsRef.current.set(roomId, ws)
    console.log('Registered WebSocket connection for room:', roomId)
  }, [])

  // Function to unregister WebSocket connection (used internally)
  const unregisterWebSocketConnection = useCallback((roomId: string) => {
    wsConnectionsRef.current.delete(roomId)
    console.log('Unregistered WebSocket connection for room:', roomId)
  }, [])

  const value: MessageContextType = {
    messages,
    isLoading,
    isLoadingMore,
    error,
    sendError,
    hasMoreMessages,
    loadMessages,
    loadMoreMessages,
    sendMessage,
    registerWebSocket: registerWebSocketConnection,
    unregisterWebSocket: unregisterWebSocketConnection,
    addRealTimeMessage,
    clearMessages,
    clearSendError
  }

  return (
    <MessageContext.Provider value={value}>
      {children}
    </MessageContext.Provider>
  )
}
