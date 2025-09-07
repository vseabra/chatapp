import { useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ArrowLeft, Users, RefreshCw } from 'lucide-react'
import { MessageList } from './MessageList'
import { useMessages } from '@/contexts/MessageContext'
import { useChatRooms } from '@/contexts/ChatRoomContext'

export function ChatRoom() {
  const { roomId } = useParams<{ roomId: string }>()
  const navigate = useNavigate()
  const { loadMessages, clearMessages, isLoading } = useMessages()
  const { chatRooms } = useChatRooms()
  const loadedRoomRef = useRef<string | null>(null)

  // Find the current room details
  const currentRoom = chatRooms.find(room => room.id === roomId)

  useEffect(() => {
    // Only load messages if we haven't already loaded for this room
    if (roomId && loadedRoomRef.current !== roomId) {
      console.log('Loading messages for room:', roomId)
      loadedRoomRef.current = roomId
      loadMessages(roomId, 50)
    }

    return () => {
      // Clear messages when leaving the room
      if (roomId) {
        console.log('Clearing messages for room:', roomId)
        loadedRoomRef.current = null
        clearMessages()
      }
    }
  }, [roomId]) // Only depend on roomId

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
                <p className="text-sm text-gray-600">Chat Room</p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleRefresh}
              disabled={isLoading}
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>
        </div>
      </div>

      {/* Chat Messages */}
      <div className="max-w-4xl mx-auto h-[calc(100vh-120px)] flex flex-col">
        <div className="flex-1 bg-white mx-8 my-4 rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <MessageList />
        </div>

        {/* Message Input (placeholder for future) */}
        <div className="mx-8 mb-4">
          <Card>
            <CardContent className="pt-6">
              <div className="text-center text-gray-500">
                <p className="text-sm">💬 Message input will be implemented in the next step</p>
                <p className="text-xs mt-1">Currently showing historical messages only</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
