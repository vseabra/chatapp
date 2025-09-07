import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Trash2, Users, Crown, MessageCircle } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useChatRooms } from '@/contexts/ChatRoomContext'
import { useAuth } from '@/contexts/AuthContext'

export function ChatRoomList() {
  const { chatRooms, isLoading, error, deleteChatRoom } = useChatRooms()
  const { user } = useAuth()
  const navigate = useNavigate()

  const handleDelete = async (roomId: string, roomTitle: string) => {
    if (window.confirm(`Are you sure you want to delete "${roomTitle}"?`)) {
      try {
        await deleteChatRoom(roomId)
      } catch (error) {
        console.error('Failed to delete chatroom:', error)
      }
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading chatrooms...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-600 mb-4">{error}</p>
        <Button onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {chatRooms.length === 0 ? (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <Users className="mx-auto h-12 w-12 text-gray-400 mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No chatrooms yet</h3>
              <p className="text-gray-600">Create your first chatroom to get started!</p>
            </div>
          </CardContent>
        </Card>
      ) : (
        chatRooms.map((room) => (
          <Card key={room.id} className="hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Users className="h-5 w-5 text-gray-500" />
                  <CardTitle className="text-lg">{room.title}</CardTitle>
                </div>
                <div className="flex items-center space-x-2">
                  {room.ownerId === user?.id && (
                    <div className="flex items-center text-sm text-gray-600">
                      <Crown className="h-4 w-4 mr-1" />
                      Owner
                    </div>
                  )}
                  {room.ownerId === user?.id && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDelete(room.id, room.title)}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>
              <CardDescription>
                Created by {room.ownerId === user?.id ? 'you' : 'another user'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-600">
                  Room ID: {room.id}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => navigate(`/chat/${room.id}`)}
                  className="hover:bg-blue-50 hover:border-blue-300"
                >
                  <MessageCircle className="h-4 w-4 mr-2" />
                  Join Chat
                </Button>
              </div>
            </CardContent>
          </Card>
        ))
      )}
    </div>
  )
}
