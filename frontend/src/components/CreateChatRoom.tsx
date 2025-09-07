import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Plus } from 'lucide-react'
import { useChatRooms } from '@/contexts/ChatRoomContext'

export function CreateChatRoom() {
  const [title, setTitle] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [isExpanded, setIsExpanded] = useState(false)
  const { createChatRoom } = useChatRooms()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return

    setIsSubmitting(true)
    setError('')

    try {
      await createChatRoom(title.trim())
      setTitle('')
      setIsExpanded(false)
    } catch (error) {
      console.error('Failed to create chatroom:', error)
      setError('Failed to create chatroom')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    setTitle('')
    setError('')
    setIsExpanded(false)
  }

  if (!isExpanded) {
    return (
      <Card>
        <CardContent className="pt-6">
          <Button
            onClick={() => setIsExpanded(true)}
            className="w-full"
            size="lg"
          >
            <Plus className="h-5 w-5 mr-2" />
            Create New Chatroom
          </Button>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <Plus className="h-5 w-5 mr-2" />
          Create New Chatroom
        </CardTitle>
        <CardDescription>
          Create a new chatroom that others can join
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Chatroom Title</Label>
            <Input
              id="title"
              type="text"
              placeholder="Enter chatroom title..."
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              maxLength={120}
              autoFocus
            />
            <p className="text-xs text-gray-500">
              {title.length}/120 characters
            </p>
          </div>

          {error && (
            <div className="text-sm text-red-600">
              {error}
            </div>
          )}

          <div className="flex space-x-2">
            <Button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="flex-1"
            >
              {isSubmitting ? 'Creating...' : 'Create Chatroom'}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={handleCancel}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
