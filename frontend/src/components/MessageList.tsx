import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useMessages } from '@/contexts/MessageContext'
import { MessageItem } from './MessageItem'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { ScrollArea } from '@/components/ui/scroll-area'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Button } from '@/components/ui/button'
import { Loader2, ArrowDown, ChevronUp } from 'lucide-react'

export function MessageList() {
  const { roomId } = useParams<{ roomId: string }>()
  const {
    messages,
    isLoading,
    isLoadingMore,
    error,
    hasMoreMessages,
    loadMoreMessages
  } = useMessages()
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const [showScrollButton, setShowScrollButton] = useState(false)

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [messages])

  // Show/hide scroll to bottom button based on scroll position
  const handleScroll = (event: React.UIEvent<HTMLDivElement>) => {
    const target = event.target as HTMLDivElement
    const isNearBottom = target.scrollHeight - target.scrollTop - target.clientHeight < 100
    setShowScrollButton(!isNearBottom)
  }

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' })
      setShowScrollButton(false)
    }
  }

  const handleLoadMore = async () => {
    if (roomId && !isLoadingMore) {
      try {
        await loadMoreMessages(roomId)
      } catch (error) {
        console.error('Failed to load more messages:', error)
      }
    }
  }

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
    <div className="relative h-full">
      <ScrollArea className="h-full w-full rounded-md">
      <div
        className="min-h-full p-4 space-y-0"
        onScroll={handleScroll}
      >
        {/* Load More Button */}
        {hasMoreMessages && (
          <div className="flex justify-center py-4">
            <Button
              onClick={handleLoadMore}
              disabled={isLoadingMore}
              variant="outline"
              size="sm"
              className="flex items-center gap-2"
            >
              {isLoadingMore ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Loading...
                </>
              ) : (
                <>
                  <ChevronUp className="h-4 w-4" />
                  Load More Messages
                </>
              )}
            </Button>
          </div>
        )}

        {/* Loading indicator for load more */}
        {isLoadingMore && (
          <div className="flex justify-center py-2">
            <Loader2 className="h-4 w-4 animate-spin text-gray-400" />
          </div>
        )}

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
        <div ref={messagesEndRef} className="h-4" />
      </ScrollArea>

      {/* Scroll to bottom button */}
      {showScrollButton && messages.length > 5 && (
        <Button
          onClick={scrollToBottom}
          className="absolute bottom-4 right-4 rounded-full shadow-lg"
          size="sm"
        >
          <ArrowDown className="h-4 w-4" />
        </Button>
      )}
    </div>
  )
}
