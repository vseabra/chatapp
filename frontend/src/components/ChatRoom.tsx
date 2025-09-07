import { useEffect, useRef, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ArrowLeft, Users, RefreshCw, Wifi, WifiOff } from 'lucide-react'
import { MessageList } from './MessageList'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { MessageInput } from './MessageInput'
import { useMessages } from '@/contexts/MessageContext'
import { useChatRooms } from '@/contexts/ChatRoomContext'
import type { Message } from '@/contexts/MessageContext'
import { useAuth } from '@/contexts/AuthContext'

export function ChatRoom() {
  const { roomId } = useParams<{ roomId: string }>()
  const navigate = useNavigate()
  const { loadMessages, clearMessages, isLoading, addRealTimeMessage, registerWebSocket, unregisterWebSocket } = useMessages()
  const { chatRooms } = useChatRooms()
  const { user } = useAuth()
  const loadedRoomRef = useRef<string | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<number | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [connectionAttempts, setConnectionAttempts] = useState(0)

  // Find the current room details
  const currentRoom = chatRooms.find(room => room.id === roomId)

  // WebSocket connection functions
  const getWebSocketUrl = useCallback(() => {
    const token = localStorage.getItem('accessToken')
    const baseUrl = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
      ? 'ws://localhost:8080'
      : 'ws://chat-api:8080'

    // Send JWT token as query parameter for authentication
    const params = new URLSearchParams({
      roomId: roomId || '',
    })

    // Add token if available
    if (token) {
      params.append('token', token)
    }

    return `${baseUrl}/api/v1/ws?${params.toString()}`
  }, [roomId])

  const connectWebSocket = useCallback(() => {
    if (!roomId || !user) return

    // Prevent multiple connections for the same room
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected for room:', roomId)
      return
    }

    console.log('Connecting to WebSocket for room:', roomId)

    // Close existing connection if any
    if (wsRef.current) {
      console.log('Closing existing WebSocket connection')
      wsRef.current.close()
    }

    const wsUrl = getWebSocketUrl()
    console.log('WebSocket URL:', wsUrl)

    const token = localStorage.getItem('accessToken')
    console.log('JWT Token present:', !!token)
    if (token) {
      console.log('Token length:', token.length)
    }

    const ws = new WebSocket(wsUrl)

    // Add connection timeout
    const connectionTimeout = setTimeout(() => {
      if (ws.readyState === WebSocket.CONNECTING) {
        console.log('WebSocket connection timeout, closing...')
        ws.close()
      }
    }, 10000) // 10 second timeout

    ws.onopen = () => {
      console.log('WebSocket connected for room:', roomId)
      clearTimeout(connectionTimeout)

      // Send JWT token as first message for authentication (fallback)
      const token = localStorage.getItem('accessToken')
      if (token) {
        console.log('Sending JWT token for authentication')
        ws.send(JSON.stringify({
          type: 'auth',
          token: token
        }))
      }

      // Register this WebSocket connection for sending messages
      if (roomId) {
        registerWebSocket(roomId, ws)
      }

      setIsConnected(true)
      setConnectionAttempts(0)
    }

    ws.onmessage = (event) => {
      console.log('Raw WebSocket message received:', event.data)

      try {
        const messageData = JSON.parse(event.data)
        console.log('Parsed WebSocket message:', messageData)

        // Handle authentication response
        if (messageData.type === 'auth' && messageData.status === 'success') {
          console.log('Authentication successful')
          return
        }

        // Handle authentication error
        if (messageData.type === 'auth' && messageData.status === 'error') {
          console.error('Authentication failed:', messageData.error)
          ws.close()
          return
        }

        // Check if it's a MessageCreated event
        // Bot messages have empty userId, so we need to check for roomId and text
        if (messageData.roomId && messageData.text && messageData.id) {
          console.log('Processing message:', messageData.id, 'type:', messageData.type)
          const newMessage: Message = {
            id: messageData.id,
            roomId: messageData.roomId,
            userId: messageData.userId || '', // Allow empty userId for bot messages
            userName: messageData.userName || 'Bot', // Default to 'Bot' if empty
            text: messageData.text,
            type: messageData.type || 'message',
            createdAt: messageData.createdAt
          }

          // Add the real-time message to the context
          addRealTimeMessage(newMessage)
          console.log('Message added to chat:', messageData.type, messageData.text)
        } else {
          console.log('Received unknown message type:', messageData)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
        console.error('Raw message data:', event.data)
      }
    }

    ws.onclose = (event) => {
      clearTimeout(connectionTimeout)
      console.log('WebSocket disconnected for room:', roomId, 'code:', event.code, 'reason:', event.reason)
      console.log('WebSocket readyState:', ws.readyState)
      setIsConnected(false)

      // Attempt to reconnect if it wasn't a clean disconnect
      if (event.code !== 1000 && connectionAttempts < 5) {
        const delay = Math.min(1000 * Math.pow(2, connectionAttempts), 30000)
        console.log(`Attempting to reconnect in ${delay}ms... (attempt ${connectionAttempts + 1}/5)`)

        reconnectTimeoutRef.current = setTimeout(() => {
          setConnectionAttempts(prev => prev + 1)
          connectWebSocket()
        }, delay)
      }
    }

    ws.onerror = (error) => {
      clearTimeout(connectionTimeout)
      console.error('WebSocket error for room:', roomId, error)
      console.error('WebSocket readyState:', ws.readyState)
      console.error('Error details:', error)
    }

    wsRef.current = ws
  }, [roomId, user, getWebSocketUrl, addRealTimeMessage, connectionAttempts])

  const disconnectWebSocket = useCallback(() => {
    console.log('Disconnecting WebSocket for room:', roomId)

    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    setIsConnected(false)
  }, [roomId])

  useEffect(() => {
    // Only load messages and connect WebSocket if we haven't already for this room
    if (roomId && loadedRoomRef.current !== roomId && !isLoading) {
      console.log('Loading messages for room:', roomId)
      loadedRoomRef.current = roomId
      loadMessages(roomId, 50)

      // Connect to WebSocket for real-time messages
      setTimeout(() => connectWebSocket(), 100) // Small delay to ensure messages are loaded first
    }

    return () => {
      // Clear messages and disconnect WebSocket when leaving the room
      if (roomId) {
        console.log('Clearing messages and disconnecting WebSocket for room:', roomId)
        loadedRoomRef.current = null
        unregisterWebSocket(roomId)
        disconnectWebSocket()
        clearMessages()
      }
    }
  }, [roomId]) // Only depend on roomId to prevent unnecessary re-runs

  const handleRefresh = () => {
    if (roomId) {
      // Reset the loaded room ref to allow re-loading
      loadedRoomRef.current = null
      loadMessages(roomId, 50)
    }
  }

  if (!currentRoom) {
    return (
      <div className="min-h-screen bg-gray-50 p-8">
        <div className="max-w-4xl mx-auto">
          <div className="text-center py-12">
            <Users className="mx-auto h-16 w-16 text-gray-400 mb-4" />
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Chat Room Not Found</h2>
            <p className="text-gray-600 mb-6">The chat room you're looking for doesn't exist or has been deleted.</p>
            <Button onClick={() => navigate('/chatrooms')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Chat Rooms
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-8 py-4">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => navigate('/chatrooms')}
              >
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Rooms
              </Button>
              <div>
                <h1 className="text-xl font-semibold text-gray-900">{currentRoom.title}</h1>
                <div className="flex items-center space-x-2">
                  <p className="text-sm text-gray-600">Chat Room</p>
                  <div className={`flex items-center space-x-1 text-xs px-2 py-1 rounded-full ${
                    isConnected
                      ? 'bg-green-100 text-green-700'
                      : 'bg-red-100 text-red-700'
                  }`}>
                    {isConnected ? (
                      <>
                        <Wifi className="h-3 w-3" />
                        <span>Live</span>
                      </>
                    ) : (
                      <>
                        <WifiOff className="h-3 w-3" />
                        <span>Disconnected</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleRefresh}
                disabled={isLoading}
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
              {!isConnected && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={connectWebSocket}
                  className="text-blue-600 hover:text-blue-700"
                >
                  <Wifi className="h-4 w-4 mr-2" />
                  Reconnect
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Chat Messages */}
      <div className="max-w-4xl mx-auto h-[calc(100vh-200px)] flex flex-col">
        <div className="h-full bg-white mx-8 my-4 rounded-lg shadow-sm border border-gray-200">
          <div className="h-full overflow-hidden">
            <MessageList />
          </div>
        </div>

        {/* Message Input */}
        <div className="mx-8 mb-4">
          {roomId && (
            <MessageInput
              roomId={roomId}
              disabled={!isConnected && !user}
            />
          )}
        </div>
      </div>
    </div>
  )
}
