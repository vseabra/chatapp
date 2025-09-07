import type { Message } from '@/contexts/MessageContext'
import { useAuth } from '@/contexts/AuthContext'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Bot } from 'lucide-react'
import { cn } from '@/lib/utils'

interface MessageItemProps {
  message: Message
  showAvatar?: boolean
}

export function MessageItem({ message, showAvatar = true }: MessageItemProps) {
  const { user } = useAuth()
  const isOwnMessage = message.userId === user?.id
  const isBotMessage = message.type === 'bot' || (!message.userId && !message.userName)

  // Format timestamp
  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  // Get user initials for avatar
  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  return (
    <div className={cn(
      "flex gap-3 p-4 hover:bg-gray-50 transition-colors",
      isOwnMessage ? "bg-blue-50" : "",
      isBotMessage ? "bg-green-50 border-l-4 border-green-400" : ""
    )}>
      {showAvatar && (
        <Avatar className="w-8 h-8 flex-shrink-0">
          <AvatarFallback className={cn(
            "text-xs font-medium",
            isOwnMessage ? "bg-blue-500 text-white" :
            isBotMessage ? "bg-green-500 text-white" :
            "bg-gray-500 text-white"
          )}>
            {isBotMessage ? (
              <Bot className="w-4 h-4" />
            ) : (
              getInitials(message.userName || 'User')
            )}
          </AvatarFallback>
        </Avatar>
      )}

      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2 mb-1">
          <span className={cn(
            "font-medium text-sm flex items-center gap-1",
            isOwnMessage ? "text-blue-700" :
            isBotMessage ? "text-green-700" :
            "text-gray-900"
          )}>
            {isBotMessage ? (
              <>
                <Bot className="w-4 h-4" />
                {message.userName || 'Bot'}
              </>
            ) : (
              <>
                {message.userName}
                {isOwnMessage && (
                  <span className="text-xs text-blue-500">(You)</span>
                )}
              </>
            )}
          </span>
          <span className="text-xs text-gray-500">
            {formatTime(message.createdAt)}
          </span>
        </div>

        <div className={cn(
          "text-sm break-words",
          isOwnMessage ? "text-blue-800" :
          isBotMessage ? "text-green-800 font-medium" :
          "text-gray-700"
        )}>
          {message.text}
        </div>
      </div>
    </div>
  )
}
