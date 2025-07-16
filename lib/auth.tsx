'use client'

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

// API Base URL - same logic as api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// Types
export interface User {
  id: number
  username: string
  email: string
  first_name: string
  last_name: string
  role: string
  is_active: boolean
  last_login?: string
  created_at: string
  updated_at: string
}

export interface AuthState {
  user: User | null
  token: string | null
  isLoading: boolean
  isAuthenticated: boolean
}

export interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<{ success: boolean; error?: string }>
  register: (data: RegisterData) => Promise<{ success: boolean; error?: string }>
  logout: () => void
  updateProfile: (data: Partial<User>) => Promise<{ success: boolean; error?: string }>
  changePassword: (currentPassword: string, newPassword: string) => Promise<{ success: boolean; error?: string }>
}

export interface RegisterData {
  username: string
  email: string
  password: string
  first_name: string
  last_name: string
}

// Create context
const AuthContext = createContext<AuthContextType | undefined>(undefined)

// Auth Provider
export function AuthProvider({ children }: { children: ReactNode }) {
  const [authState, setAuthState] = useState<AuthState>({
    user: null,
    token: null,
    isLoading: true,
    isAuthenticated: false
  })

  // Initialize auth state from localStorage on mount
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const token = localStorage.getItem('auth_token')
        const userData = localStorage.getItem('user_data')
        
        if (token && userData) {
          const user = JSON.parse(userData)
          
          // Verify token is still valid by calling profile endpoint
          const response = await fetch(`${API_BASE_URL}/api/profile`, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json'
            }
          })
          
          if (response.ok) {
            const profileData = await response.json()
            setAuthState({
              user: profileData.user,
              token,
              isLoading: false,
              isAuthenticated: true
            })
          } else {
            // Token invalid, clear storage
            localStorage.removeItem('auth_token')
            localStorage.removeItem('user_data')
            setAuthState({
              user: null,
              token: null,
              isLoading: false,
              isAuthenticated: false
            })
          }
        } else {
          setAuthState({
            user: null,
            token: null,
            isLoading: false,
            isAuthenticated: false
          })
        }
      } catch (error) {
        console.error('Auth initialization error:', error)
        setAuthState({
          user: null,
          token: null,
          isLoading: false,
          isAuthenticated: false
        })
      }
    }

    initializeAuth()
  }, [])

  const login = async (username: string, password: string): Promise<{ success: boolean; error?: string }> => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
      })

      if (response.ok) {
        const data = await response.json()
        const { token, user } = data
        
        // Store in localStorage
        localStorage.setItem('auth_token', token)
        localStorage.setItem('user_data', JSON.stringify(user))
        
        // Update state
        setAuthState({
          user,
          token,
          isLoading: false,
          isAuthenticated: true
        })
        
        return { success: true }
      } else {
        const errorData = await response.json()
        return { success: false, error: errorData.error || 'Login failed' }
      }
    } catch (error) {
      console.error('Login error:', error)
      return { success: false, error: 'Network error. Please try again.' }
    }
  }

  const register = async (data: RegisterData): Promise<{ success: boolean; error?: string }> => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })

      if (response.ok) {
        const responseData = await response.json()
        const { token, user } = responseData
        
        // Store in localStorage
        localStorage.setItem('auth_token', token)
        localStorage.setItem('user_data', JSON.stringify(user))
        
        // Update state
        setAuthState({
          user,
          token,
          isLoading: false,
          isAuthenticated: true
        })
        
        return { success: true }
      } else {
        const errorData = await response.json()
        return { success: false, error: errorData.error || 'Registration failed' }
      }
    } catch (error) {
      console.error('Registration error:', error)
      return { success: false, error: 'Network error. Please try again.' }
    }
  }

  const logout = () => {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user_data')
    setAuthState({
      user: null,
      token: null,
      isLoading: false,
      isAuthenticated: false
    })
  }

  const updateProfile = async (data: Partial<User>): Promise<{ success: boolean; error?: string }> => {
    if (!authState.token) {
      return { success: false, error: 'Not authenticated' }
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/profile`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authState.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })

      if (response.ok) {
        const responseData = await response.json()
        const updatedUser = responseData.user
        
        // Update localStorage
        localStorage.setItem('user_data', JSON.stringify(updatedUser))
        
        // Update state
        setAuthState(prev => ({
          ...prev,
          user: updatedUser
        }))
        
        return { success: true }
      } else {
        const errorData = await response.json()
        return { success: false, error: errorData.error || 'Update failed' }
      }
    } catch (error) {
      console.error('Profile update error:', error)
      return { success: false, error: 'Network error. Please try again.' }
    }
  }

  const changePassword = async (currentPassword: string, newPassword: string): Promise<{ success: boolean; error?: string }> => {
    if (!authState.token) {
      return { success: false, error: 'Not authenticated' }
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/change-password`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authState.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword
        })
      })

      if (response.ok) {
        return { success: true }
      } else {
        const errorData = await response.json()
        return { success: false, error: errorData.error || 'Password change failed' }
      }
    } catch (error) {
      console.error('Password change error:', error)
      return { success: false, error: 'Network error. Please try again.' }
    }
  }

  return (
    <AuthContext.Provider
      value={{
        ...authState,
        login,
        register,
        logout,
        updateProfile,
        changePassword
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

// Custom hook to use auth context
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

// Helper hook for checking roles
export function useAuthRole(requiredRole?: string): boolean {
  const { user, isAuthenticated } = useAuth()
  
  if (!isAuthenticated || !user) {
    return false
  }
  
  if (!requiredRole) {
    return true // Just check if authenticated
  }
  
  // Admin can access everything
  if (user.role === 'admin') {
    return true
  }
  
  return user.role === requiredRole
}

// Helper hook for checking if user can manage documents
export function useCanManageDocuments(): boolean {
  const { user, isAuthenticated } = useAuth()
  
  if (!isAuthenticated || !user) {
    return false
  }
  
  return user.role === 'admin' || user.role === 'it_security'
} 