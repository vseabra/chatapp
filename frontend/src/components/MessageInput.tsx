import { useState, useRef, useEffect } from 'react'
import type { FormEvent, KeyboardEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Send, Loader2, AlertCircle } from 'lucide-react'
import { useMessages } from '@/contexts/MessageContext'
import { useAuth } from '@/contexts/AuthContext'

interface MessageInputProps {
  roomId: string
  disabled?: boolean
}

export function MessageInput({ roomId, disabled = false }: MessageInputProps) {
  const [message, setMessage] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const { sendMessage, sendError, clearSendError } = useMessages()
  const { user } = useAuth()

  const maxLength = 500
  const isEmpty = message.trim().length === 0
  const isTooLong = message.length > maxLength

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()

    if (isEmpty || isTooLong || !user || isSubmitting) {
      return
    }

    const messageToSend = message.trim()
    setIsSubmitting(true)

    try {
      await sendMessage(roomId, messageToSend, user.id, user.name)
      setMessage('')

      // Reset textarea height
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto'
      }
    } catch (error) {
      console.error('Failed to send message:', error)
      // Error is handled by the context, so we don't need to do anything here
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e as any)
    }
  }

  const handleTextareaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value)

    // Clear send error when user starts typing
    if (sendError) {
      clearSendError()
    }

    // Auto-resize textarea
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`
    }
  }

  // Clear error when component unmounts or room changes
  useEffect(() => {
    return () => {
      if (sendError) {
        clearSendError()
      }
    }
  }, [roomId])

  const isDisabled = disabled || isSubmitting

  return (
    <div className="border-t border-gray-200 bg-white p-4">
      <form onSubmit={handleSubmit} className="flex items-end gap-3">
        <div className="flex-1 relative">
          <Textarea
            ref={textareaRef}
            value={message}
            onChange={handleTextareaChange}
            onKeyDown={handleKeyDown}
            placeholder="Type your message..."
            className="min-h-[40px] max-h-[120px] resize-none pr-12"
            disabled={isDisabled}
            rows={1}
          />

          {/* Character count */}
          <div className="absolute bottom-2 right-2 text-xs text-gray-400">
            {message.length}/{maxLength}
          </div>
        </div>

        <Button
          type="submit"
          size="sm"
          disabled={isDisabled || isEmpty || isTooLong}
          className="h-10 px-4"
        >
          {isSubmitting ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Send className="h-4 w-4" />
          )}
        </Button>
      </form>

      {/* Error states */}
      <div className="flex flex-col space-y-1 mt-2">
        {sendError && (
          <div className="flex items-center space-x-1 text-xs text-red-600">
            <AlertCircle className="h-3 w-3" />
            <span>{sendError}</span>
          </div>
        )}
        {isTooLong && (
          <p className="text-xs text-red-500">
            Message is too long (max {maxLength} characters)
          </p>
        )}
      </div>
    </div>
  )
}
