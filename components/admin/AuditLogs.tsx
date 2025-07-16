'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { 
  Clock, 
  User, 
  FileText, 
  Settings,
  Search,
  Filter,
  ChevronLeft,
  ChevronRight,
  Calendar,
  Shield,
  Trash2,
  Edit,
  Eye,
  LogIn
} from 'lucide-react'
import { getAuditLogs, type AuditLog, type AuditLogFilter } from '@/lib/api'

const ACTION_ICONS = {
  CREATE: <FileText className="w-4 h-4 text-green-600" />,
  UPDATE: <Edit className="w-4 h-4 text-blue-600" />,
  DELETE: <Trash2 className="w-4 h-4 text-red-600" />,
  VIEW: <Eye className="w-4 h-4 text-gray-600" />,
  LOGIN: <LogIn className="w-4 h-4 text-purple-600" />,
  LOGOUT: <LogIn className="w-4 h-4 text-gray-600" />,
}

const ACTION_COLORS = {
  CREATE: 'bg-green-100 text-green-800',
  UPDATE: 'bg-blue-100 text-blue-800',
  DELETE: 'bg-red-100 text-red-800',
  VIEW: 'bg-gray-100 text-gray-800',
  LOGIN: 'bg-purple-100 text-purple-800',
  LOGOUT: 'bg-gray-100 text-gray-800',
}

const RESOURCE_ICONS = {
  USER: <User className="w-4 h-4" />,
  DOCUMENT: <FileText className="w-4 h-4" />,
  SYSTEM: <Settings className="w-4 h-4" />,
}

export default function AuditLogs() {
  const [filters, setFilters] = useState<AuditLogFilter>({
    page: 1,
    limit: 25,
  })
  const [searchQuery, setSearchQuery] = useState('')

  // Fetch audit logs
  const { data, isLoading, error } = useQuery({
    queryKey: ['audit-logs', filters],
    queryFn: () => getAuditLogs(filters),
  })

  const handleFilterChange = (key: keyof AuditLogFilter, value: any) => {
    setFilters(prev => ({ ...prev, [key]: value, page: 1 }))
  }

  const handlePageChange = (newPage: number) => {
    setFilters(prev => ({ ...prev, page: newPage }))
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  const getActionBadge = (action: string) => {
    const color = ACTION_COLORS[action as keyof typeof ACTION_COLORS] || 'bg-gray-100 text-gray-800'
    const icon = ACTION_ICONS[action as keyof typeof ACTION_ICONS] || <Settings className="w-4 h-4" />
    
    return (
      <Badge className={`${color} border-0 flex items-center gap-1`}>
        {icon}
        {action}
      </Badge>
    )
  }

  const getResourceIcon = (resourceType: string) => {
    return RESOURCE_ICONS[resourceType as keyof typeof RESOURCE_ICONS] || <Settings className="w-4 h-4" />
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-48 mb-6"></div>
          <div className="space-y-3">
            {[...Array(10)].map((_, i) => (
              <div key={i} className="h-20 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <Clock className="w-12 h-12 text-red-400 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-2">Failed to load audit logs</h3>
        <p className="text-red-600">{error.message}</p>
      </div>
    )
  }

  const auditLogs = data?.audit_logs || []
  const pagination = data?.pagination

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Clock className="w-6 h-6 text-blue-600" />
          <h2 className="text-2xl font-bold text-gray-900">Audit Logs</h2>
        </div>
        <div className="text-sm text-gray-600">
          {pagination && (
            <span>
              Showing {(pagination.page - 1) * pagination.limit + 1} - {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} entries
            </span>
          )}
        </div>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <label className="text-sm font-medium">Action</label>
              <select
                value={filters.action || ''}
                onChange={(e) => handleFilterChange('action', e.target.value || undefined)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="">All Actions</option>
                <option value="CREATE">Create</option>
                <option value="UPDATE">Update</option>
                <option value="DELETE">Delete</option>
                <option value="VIEW">View</option>
                <option value="LOGIN">Login</option>
                <option value="LOGOUT">Logout</option>
              </select>
            </div>

            <div>
              <label className="text-sm font-medium">Resource Type</label>
              <select
                value={filters.resource_type || ''}
                onChange={(e) => handleFilterChange('resource_type', e.target.value || undefined)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value="">All Resources</option>
                <option value="USER">User</option>
                <option value="DOCUMENT">Document</option>
                <option value="SYSTEM">System</option>
              </select>
            </div>

            <div>
              <label className="text-sm font-medium">From Date</label>
              <Input
                type="date"
                value={filters.from || ''}
                onChange={(e) => handleFilterChange('from', e.target.value || undefined)}
              />
            </div>

            <div>
              <label className="text-sm font-medium">To Date</label>
              <Input
                type="date"
                value={filters.to || ''}
                onChange={(e) => handleFilterChange('to', e.target.value || undefined)}
              />
            </div>
          </div>

          <div className="flex gap-2 mt-4">
            <Button 
              variant="outline" 
              onClick={() => setFilters({ page: 1, limit: 25 })}
            >
              Clear Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Audit Logs Table */}
      <Card>
        <CardContent className="p-0">
          <div className="space-y-0">
            {auditLogs.length > 0 ? (
              auditLogs.map((log) => (
                <div key={log.id} className="p-4 border-b hover:bg-gray-50">
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-4">
                      <div className="flex items-center gap-2">
                        {getResourceIcon(log.resource_type)}
                        <div className="flex flex-col gap-1">
                          {getActionBadge(log.action)}
                          <span className="text-xs text-gray-500">
                            {formatDate(log.created_at)}
                          </span>
                        </div>
                      </div>
                      
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="font-medium">
                            {log.user?.first_name} {log.user?.last_name}
                          </span>
                          <span className="text-sm text-gray-600">
                            (@{log.user?.username})
                          </span>
                          <Badge variant="outline" className="text-xs">
                            {log.user?.role}
                          </Badge>
                        </div>
                        
                        <div className="text-sm text-gray-600 mb-2">
                          <span className="font-medium">{log.resource_type}</span>
                          {log.resource_name && (
                            <span> • {log.resource_name}</span>
                          )}
                          {log.resource_id && (
                            <span className="text-xs text-gray-500"> (ID: {log.resource_id})</span>
                          )}
                        </div>
                        
                        {log.details && (
                          <p className="text-sm text-gray-700">{log.details}</p>
                        )}
                        
                        <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                          {log.ip_address && (
                            <span>IP: {log.ip_address}</span>
                          )}
                          {log.user_agent && (
                            <span className="truncate max-w-xs">
                              UA: {log.user_agent.split(' ')[0]}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-12">
                <Clock className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No audit logs found</h3>
                <p className="text-gray-500">
                  {Object.values(filters).some(v => v !== undefined && v !== 1 && v !== 25) 
                    ? "Try adjusting your filters to see more results."
                    : "System activity will appear here as users interact with the platform."
                  }
                </p>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Pagination */}
      {pagination && pagination.pages > 1 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-gray-600">
            Page {pagination.page} of {pagination.pages}
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => handlePageChange(pagination.page - 1)}
              disabled={pagination.page <= 1}
            >
              <ChevronLeft className="w-4 h-4" />
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => handlePageChange(pagination.page + 1)}
              disabled={pagination.page >= pagination.pages}
            >
              Next
              <ChevronRight className="w-4 h-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  )
} 