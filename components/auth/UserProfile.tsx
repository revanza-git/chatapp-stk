'use client'

import { useState } from 'react'
import { User, LogOut, Settings, Key } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { useAuth, type User as UserType } from '@/lib/auth'

interface UserProfileProps {
  variant?: 'compact' | 'full'
}

export default function UserProfile({ variant = 'compact' }: UserProfileProps) {
  const { user, logout } = useAuth()
  const [isProfileOpen, setIsProfileOpen] = useState(false)
  
  if (!user) return null

  if (variant === 'compact') {
    return (
      <div className="flex items-center gap-2">
        <Dialog open={isProfileOpen} onOpenChange={setIsProfileOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" size="sm" className="gap-2">
              <User className="w-4 h-4" />
              <span className="hidden sm:inline">{user.first_name}</span>
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>User Profile</DialogTitle>
            </DialogHeader>
            <UserProfileContent user={user} onClose={() => setIsProfileOpen(false)} />
          </DialogContent>
        </Dialog>
        
        <Button 
          variant="outline" 
          size="sm" 
          onClick={logout}
          className="text-red-600 hover:text-red-700 gap-2"
        >
          <LogOut className="w-4 h-4" />
          <span className="hidden sm:inline">Logout</span>
        </Button>
      </div>
    )
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <User className="w-5 h-5" />
          User Profile
        </CardTitle>
      </CardHeader>
      <CardContent>
        <UserProfileContent user={user} />
      </CardContent>
    </Card>
  )
}

function UserProfileContent({ user, onClose }: { user: UserType, onClose?: () => void }) {
  const { updateProfile, changePassword, logout } = useAuth()
  const [isEditing, setIsEditing] = useState(false)
  const [isChangingPassword, setIsChangingPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const [profileData, setProfileData] = useState({
    first_name: user.first_name,
    last_name: user.last_name,
    email: user.email,
  })

  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  })

  const handleProfileUpdate = async () => {
    setError('')
    setSuccess('')
    setIsLoading(true)

    try {
      const result = await updateProfile(profileData)
      
      if (result.success) {
        setSuccess('Profile updated successfully!')
        setIsEditing(false)
        setTimeout(() => setSuccess(''), 3000)
      } else {
        setError(result.error || 'Failed to update profile')
      }
    } catch (err) {
      setError('An unexpected error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  const handlePasswordChange = async () => {
    setError('')
    setSuccess('')

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setError('New passwords do not match')
      return
    }

    if (passwordData.newPassword.length < 8) {
      setError('New password must be at least 8 characters long')
      return
    }

    setIsLoading(true)

    try {
      const result = await changePassword(passwordData.currentPassword, passwordData.newPassword)
      
      if (result.success) {
        setSuccess('Password changed successfully!')
        setIsChangingPassword(false)
        setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' })
        setTimeout(() => setSuccess(''), 3000)
      } else {
        setError(result.error || 'Failed to change password')
      }
    } catch (err) {
      setError('An unexpected error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      {/* User Info Display */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold">Account Information</h3>
          {!isEditing && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsEditing(true)}
              className="gap-2"
            >
              <Settings className="w-4 h-4" />
              Edit
            </Button>
          )}
        </div>

        {!isEditing ? (
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-sm font-medium">Name:</span>
              <span className="text-sm">{user.first_name} {user.last_name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Email:</span>
              <span className="text-sm">{user.email}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Username:</span>
              <span className="text-sm">{user.username}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Role:</span>
              <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                {user.role}
              </span>
            </div>
          </div>
        ) : (
          <div className="space-y-3">
            <div className="grid grid-cols-2 gap-2">
              <div>
                <label className="text-sm font-medium">First Name</label>
                <Input
                  value={profileData.first_name}
                  onChange={(e) => setProfileData(prev => ({ ...prev, first_name: e.target.value }))}
                  disabled={isLoading}
                />
              </div>
              <div>
                <label className="text-sm font-medium">Last Name</label>
                <Input
                  value={profileData.last_name}
                  onChange={(e) => setProfileData(prev => ({ ...prev, last_name: e.target.value }))}
                  disabled={isLoading}
                />
              </div>
            </div>
            <div>
              <label className="text-sm font-medium">Email</label>
              <Input
                type="email"
                value={profileData.email}
                onChange={(e) => setProfileData(prev => ({ ...prev, email: e.target.value }))}
                disabled={isLoading}
              />
            </div>
            <div className="flex gap-2">
              <Button
                onClick={handleProfileUpdate}
                disabled={isLoading}
                size="sm"
              >
                {isLoading ? 'Saving...' : 'Save Changes'}
              </Button>
              <Button
                variant="outline"
                onClick={() => {
                  setIsEditing(false)
                  setProfileData({
                    first_name: user.first_name,
                    last_name: user.last_name,
                    email: user.email,
                  })
                }}
                disabled={isLoading}
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Password Change */}
      <div className="pt-4 border-t">
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-semibold">Password</h3>
          {!isChangingPassword && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsChangingPassword(true)}
              className="gap-2"
            >
              <Key className="w-4 h-4" />
              Change
            </Button>
          )}
        </div>

        {isChangingPassword && (
          <div className="space-y-3">
            <div>
              <label className="text-sm font-medium">Current Password</label>
              <Input
                type="password"
                value={passwordData.currentPassword}
                onChange={(e) => setPasswordData(prev => ({ ...prev, currentPassword: e.target.value }))}
                disabled={isLoading}
              />
            </div>
            <div>
              <label className="text-sm font-medium">New Password</label>
              <Input
                type="password"
                value={passwordData.newPassword}
                onChange={(e) => setPasswordData(prev => ({ ...prev, newPassword: e.target.value }))}
                disabled={isLoading}
                minLength={8}
              />
            </div>
            <div>
              <label className="text-sm font-medium">Confirm New Password</label>
              <Input
                type="password"
                value={passwordData.confirmPassword}
                onChange={(e) => setPasswordData(prev => ({ ...prev, confirmPassword: e.target.value }))}
                disabled={isLoading}
              />
            </div>
            <div className="flex gap-2">
              <Button
                onClick={handlePasswordChange}
                disabled={isLoading}
                size="sm"
              >
                {isLoading ? 'Changing...' : 'Change Password'}
              </Button>
              <Button
                variant="outline"
                onClick={() => {
                  setIsChangingPassword(false)
                  setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' })
                }}
                disabled={isLoading}
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Messages */}
      {error && (
        <div className="text-sm text-red-600 bg-red-50 p-3 rounded-md">
          {error}
        </div>
      )}

      {success && (
        <div className="text-sm text-green-600 bg-green-50 p-3 rounded-md">
          {success}
        </div>
      )}

      {/* Logout Button for modal view */}
      {onClose && (
        <div className="pt-4 border-t">
          <Button
            variant="outline"
            onClick={() => {
              logout()
              onClose()
            }}
            className="w-full text-red-600 hover:text-red-700 gap-2"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </Button>
        </div>
      )}
    </div>
  )
} 