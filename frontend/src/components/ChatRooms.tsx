import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { ArrowLeft, RefreshCw } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { ChatRoomList } from './ChatRoomList'
import { CreateChatRoom } from './CreateChatRoom'
import { useChatRooms } from '@/contexts/ChatRoomContext'

export function ChatRooms() {
  const navigate = useNavigate()
  const { refreshChatRooms, isLoading } = useChatRooms()

  const handleRefresh = async () => {
    try {
      await refreshChatRooms()
    } catch (error) {
      console.error('Failed to refresh chatrooms:', error)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto p-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center space-x-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate('/')}
            >
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Dashboard
            </Button>
            <div>
              <h1 className="text-3xl font-bold">Chat Rooms</h1>
              <p className="text-gray-600">Manage and join chat rooms</p>
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

        {/* Create Chatroom */}
        <div className="mb-8">
          <CreateChatRoom />
        </div>

        {/* Chatroom List */}
        <div>
          <h2 className="text-xl font-semibold mb-4">Available Chat Rooms</h2>
          <ChatRoomList />
        </div>
      </div>
    </div>
  )
}
