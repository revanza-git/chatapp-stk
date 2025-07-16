'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/lib/auth'

interface ProtectedRouteProps {
  children: React.ReactNode
  requiredRole?: string
  fallbackComponent?: React.ReactNode
}

export default function ProtectedRoute({ 
  children, 
  requiredRole, 
  fallbackComponent 
}: ProtectedRouteProps) {
  const { isAuthenticated, user, isLoading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, isLoading, router])

  // Show loading state while auth is being checked
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Redirecting to login...</p>
        </div>
      </div>
    )
  }

  // Check role requirements if specified
  if (requiredRole && user) {
    const hasPermission = user.role === 'admin' || user.role === requiredRole
    
    if (!hasPermission) {
      if (fallbackComponent) {
        return <>{fallbackComponent}</>
      }
      
      return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center max-w-md">
            <div className="text-6xl mb-4">🔒</div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Access Denied</h1>
            <p className="text-gray-600 mb-4">
              You don't have permission to access this page. Required role: {requiredRole}
            </p>
            <p className="text-sm text-gray-500">
              Your current role: {user.role}
            </p>
          </div>
        </div>
      )
    }
  }

  // User is authenticated and has required permissions
  return <>{children}</>
}

// Specialized components for common use cases
export function AdminRoute({ children, fallbackComponent }: { children: React.ReactNode, fallbackComponent?: React.ReactNode }) {
  return (
    <ProtectedRoute requiredRole="admin" fallbackComponent={fallbackComponent}>
      {children}
    </ProtectedRoute>
  )
}

export function ITSecurityRoute({ children, fallbackComponent }: { children: React.ReactNode, fallbackComponent?: React.ReactNode }) {
  return (
    <ProtectedRoute requiredRole="it_security" fallbackComponent={fallbackComponent}>
      {children}
    </ProtectedRoute>
  )
}

export function DocumentManagementRoute({ children, fallbackComponent }: { children: React.ReactNode, fallbackComponent?: React.ReactNode }) {
  const { user, isAuthenticated } = useAuth()
  
  if (!isAuthenticated || !user) {
    return (
      <ProtectedRoute>
        {children}
      </ProtectedRoute>
    )
  }
  
  const canManageDocuments = user.role === 'admin' || user.role === 'it_security'
  
  if (!canManageDocuments) {
    if (fallbackComponent) {
      return <>{fallbackComponent}</>
    }
    
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center max-w-md">
          <div className="text-6xl mb-4">📋</div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Document Management</h1>
          <p className="text-gray-600 mb-4">
            Only administrators and IT security personnel can manage documents.
          </p>
          <p className="text-sm text-gray-500">
            Your current role: {user.role}
          </p>
        </div>
      </div>
    )
  }
  
  return <>{children}</>
} 