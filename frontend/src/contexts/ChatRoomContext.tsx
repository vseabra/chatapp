import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useAuth } from './AuthContext'

export interface ChatRoom {
  id: string
  title: string
  ownerId: string
}

interface ChatRoomContextType {
  chatRooms: ChatRoom[]
  isLoading: boolean
  error: string | null
  createChatRoom: (title: string) => Promise<void>
  deleteChatRoom: (id: string) => Promise<void>
  refreshChatRooms: () => Promise<void>
}

const ChatRoomContext = createContext<ChatRoomContextType | undefined>(undefined)

export function useChatRooms() {
  const context = useContext(ChatRoomContext)
  if (context === undefined) {
    throw new Error('useChatRooms must be used within a ChatRoomProvider')
  }
  return context
}

interface ChatRoomProviderProps {
  children: ReactNode
}

const getApiBaseUrl = () => {
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    return 'http://localhost:8080/api/v1/chatroom'
  } else {
    return 'http://chat-api:8080/api/v1/chatroom'
  }
}

const API_BASE_URL = getApiBaseUrl()

export function ChatRoomProvider({ children }: ChatRoomProviderProps) {
  const [chatRooms, setChatRooms] = useState<ChatRoom[]>([])
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

  const refreshChatRooms = async () => {
    if (!user) return

    setIsLoading(true)
    setError(null)

    try {
      console.log('Fetching chatrooms from:', `${API_BASE_URL}/all`)

      const response = await fetch(`${API_BASE_URL}/all`, {
        method: 'GET',
        headers: getAuthHeaders(),
      })

      console.log('Chatrooms response status:', response.status)

      if (!response.ok) {
        throw new Error('Failed to fetch chatrooms')
      }

      const data = await response.json()
      console.log('Chatrooms data:', data)

      setChatRooms(data.items || [])
    } catch (error) {
      console.error('Failed to fetch chatrooms:', error)
      setError('Failed to load chatrooms')
    } finally {
      setIsLoading(false)
    }
  }

  const createChatRoom = async (title: string) => {
    if (!user) return

    setIsLoading(true)
    setError(null)

    try {
      console.log('Creating chatroom:', { title })

      const response = await fetch(API_BASE_URL, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ title }),
      })

      console.log('Create chatroom response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Failed to create chatroom'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      const newChatRoom = await response.json()
      console.log('Created chatroom:', newChatRoom)

      // Refresh the chatrooms list
      await refreshChatRooms()
    } catch (error) {
      console.error('Failed to create chatroom:', error)
      setError('Failed to create chatroom')
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  const deleteChatRoom = async (id: string) => {
    if (!user) return

    setIsLoading(true)
    setError(null)

    try {
      console.log('Deleting chatroom:', id)

      const response = await fetch(`${API_BASE_URL}/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
      })

      console.log('Delete chatroom response status:', response.status)

      if (!response.ok) {
        let errorMessage = 'Failed to delete chatroom'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (e) {
          console.error('Failed to parse error response:', e)
        }
        throw new Error(errorMessage)
      }

      console.log('Chatroom deleted successfully')

      // Refresh the chatrooms list
      await refreshChatRooms()
    } catch (error) {
      console.error('Failed to delete chatroom:', error)
      setError('Failed to delete chatroom')
      throw error
    } finally {
      setIsLoading(false)
    }
  }

  // Load chatrooms when user logs in
  useEffect(() => {
    if (user) {
      refreshChatRooms()
    } else {
      setChatRooms([])
    }
  }, [user])

  const value: ChatRoomContextType = {
    chatRooms,
    isLoading,
    error,
    createChatRoom,
    deleteChatRoom,
    refreshChatRooms
  }

  return (
    <ChatRoomContext.Provider value={value}>
      {children}
    </ChatRoomContext.Provider>
  )
}
