import { useMessages } from '@/contexts/MessageContext'
import { MessageItem } from './MessageItem'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Loader2 } from 'lucide-react'

export function MessageList() {
  const { messages, isLoading, error } = useMessages()

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-gray-400" />
          <p className="text-gray-500">Loading messages...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">Failed to load messages</p>
          <p className="text-gray-500 text-sm">{error}</p>
        </div>
      </div>
    )
  }

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <div className="text-gray-400 mb-4">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No messages yet</h3>
          <p className="text-gray-600">Be the first to send a message in this chat room!</p>
        </div>
      </div>
    )
  }

  return (
    <ScrollArea className="flex-1">
      <div className="space-y-0">
        {messages.map((message, index) => {
          // Show avatar only if the previous message is from a different user
          const showAvatar = index === 0 || messages[index - 1].userId !== message.userId

          return (
            <MessageItem
              key={message.id}
              message={message}
              showAvatar={showAvatar}
            />
          )
        })}
      </div>

      {/* Auto-scroll anchor */}
      <div id="messages-end" className="h-4" />
    </ScrollArea>
  )
}
