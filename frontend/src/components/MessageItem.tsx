import { Message } from '@/contexts/MessageContext'
import { useAuth } from '@/contexts/AuthContext'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

interface MessageItemProps {
  message: Message
  showAvatar?: boolean
}

export function MessageItem({ message, showAvatar = true }: MessageItemProps) {
  const { user } = useAuth()
  const isOwnMessage = message.userId === user?.id

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
      isOwnMessage ? "bg-blue-50" : ""
    )}>
      {showAvatar && (
        <Avatar className="w-8 h-8 flex-shrink-0">
          <AvatarFallback className={cn(
            "text-xs font-medium",
            isOwnMessage ? "bg-blue-500 text-white" : "bg-gray-500 text-white"
          )}>
            {getInitials(message.userName)}
          </AvatarFallback>
        </Avatar>
      )}

      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2 mb-1">
          <span className={cn(
            "font-medium text-sm",
            isOwnMessage ? "text-blue-700" : "text-gray-900"
          )}>
            {message.userName}
            {isOwnMessage && (
              <span className="text-xs text-blue-500 ml-1">(You)</span>
            )}
          </span>
          <span className="text-xs text-gray-500">
            {formatTime(message.createdAt)}
          </span>
        </div>

        <div className={cn(
          "text-sm break-words",
          isOwnMessage ? "text-blue-800" : "text-gray-700"
        )}>
          {message.text}
        </div>
      </div>
    </div>
  )
}
