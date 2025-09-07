import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { MessageSquare } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/contexts/AuthContext'

export function Dashboard() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <Button onClick={logout} variant="outline">
            Sign out
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Welcome back{user?.name ? `, ${user.name}` : ''}!</CardTitle>
            <CardDescription>
              You are successfully logged in to your account.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="space-y-2">
                {user?.name && <p><strong>Name:</strong> {user.name}</p>}
                <p><strong>Email:</strong> {user?.email}</p>
                <p><strong>User ID:</strong> {user?.id}</p>
              </div>
              <div className="pt-4 border-t">
                <Button
                  onClick={() => navigate('/chatrooms')}
                  className="w-full"
                  size="lg"
                >
                  <MessageSquare className="h-5 w-5 mr-2" />
                  Manage Chat Rooms
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
