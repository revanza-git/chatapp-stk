import { 
  ChatRequest, 
  ChatResponse, 
  PolicyFile, 
  CreateDocumentRequest, 
  UpdateDocumentRequest, 
  DocumentSearchParams, 
  DocumentSearchResponse,
  DocumentStats 
} from './types';

// Use environment variable for API URL, fallback to localhost for development
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Helper function to get authentication headers
function getAuthHeaders(): Record<string, string> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return headers;
}

export async function sendChatMessage(request: ChatRequest): Promise<ChatResponse> {
  const response = await fetch(`${API_BASE_URL}/chat`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to send chat message: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function getAllPolicies(): Promise<PolicyFile[]> {
  const response = await fetch(`${API_BASE_URL}/policies`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch policies: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function checkHealth(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE_URL}/health`);

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Backend health check failed: ${response.status} ${errorText}`);
  }

  return response.json();
}

// Document Management API Functions

export async function getAllDocuments(params?: {
  type?: 'policy' | 'onboarding';
  category?: string;
  active?: boolean;
}): Promise<PolicyFile[]> {
  const searchParams = new URLSearchParams();
  
  if (params?.type) searchParams.append('type', params.type);
  if (params?.category) searchParams.append('category', params.category);
  if (params?.active !== undefined) searchParams.append('active', params.active.toString());

  const url = `${API_BASE_URL}/documents${searchParams.toString() ? `?${searchParams.toString()}` : ''}`;
  const response = await fetch(url, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch documents: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function getDocumentById(id: number): Promise<PolicyFile> {
  const response = await fetch(`${API_BASE_URL}/documents/${id}`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch document: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function createDocument(document: CreateDocumentRequest): Promise<PolicyFile> {
  const response = await fetch(`${API_BASE_URL}/documents`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(document),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to create document: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function updateDocument(id: number, updates: UpdateDocumentRequest): Promise<PolicyFile> {
  const response = await fetch(`${API_BASE_URL}/documents/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: JSON.stringify(updates),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to update document: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function deleteDocument(id: number): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/documents/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to delete document: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function downloadDocument(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/documents/${id}/download`, {
    headers: getAuthHeaders(),
  });

  // Debug logging for response
  console.log('Response status:', response.status);
  console.log('Response headers:', [...response.headers.entries()]);

  if (!response.ok) {
    const errorText = await response.text();
    console.error('Download failed:', response.status, errorText);
    throw new Error(`Failed to download document: ${response.status} ${errorText}`);
  }

  // Get filename from Content-Disposition header with improved parsing
  const contentDisposition = response.headers.get('Content-Disposition');
  let filename = `document_${id}`; // Default fallback
  
  console.log('Content-Disposition header:', contentDisposition); // Debug logging
  
  if (contentDisposition) {
    // Try multiple patterns for filename extraction
    let filenameMatch = contentDisposition.match(/filename="([^"]+)"/);
    if (!filenameMatch) {
      // Try without quotes
      filenameMatch = contentDisposition.match(/filename=([^;]+)/);
    }
    if (!filenameMatch) {
      // Try filename* (RFC 5987 format)
      filenameMatch = contentDisposition.match(/filename\*=UTF-8''([^;]+)/);
      if (filenameMatch) {
        // Decode URI component
        try {
          filenameMatch[1] = decodeURIComponent(filenameMatch[1]);
        } catch (e) {
          console.warn('Failed to decode filename:', e);
        }
      }
    }
    
    if (filenameMatch && filenameMatch[1]) {
      filename = filenameMatch[1].trim();
      console.log('Extracted filename:', filename); // Debug logging
    } else {
      console.warn('Could not extract filename from Content-Disposition:', contentDisposition);
    }
  } else {
    console.warn('No Content-Disposition header found');
  }
  
  // If filename still doesn't have extension, try to determine from Content-Type
  if (!filename.includes('.')) {
    const contentType = response.headers.get('Content-Type');
    let extension = '';
    
    console.log('Content-Type:', contentType); // Debug logging
    
    if (contentType) {
      if (contentType.includes('pdf')) {
        extension = '.pdf';
      } else if (contentType.includes('wordprocessingml') || contentType.includes('docx')) {
        extension = '.docx';
      } else if (contentType.includes('text/plain')) {
        extension = '.txt';
      } else if (contentType.includes('markdown')) {
        extension = '.md';
      }
    }
    
    if (extension) {
      filename += extension;
      console.log('Added extension, final filename:', filename); // Debug logging
    }
  }

  console.log('Final download filename:', filename); // Debug logging

  // Create blob and download
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

export async function searchDocuments(params: DocumentSearchParams): Promise<DocumentSearchResponse> {
  const searchParams = new URLSearchParams();
  
  if (params.q) searchParams.append('q', params.q);
  if (params.type) searchParams.append('type', params.type);
  if (params.category) searchParams.append('category', params.category);

  const response = await fetch(`${API_BASE_URL}/documents/search?${searchParams.toString()}`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to search documents: ${response.status} ${errorText}`);
  }

  return response.json();
}

// Helper function to calculate document statistics
export async function getDocumentStats(): Promise<DocumentStats> {
  const documents = await getAllDocuments({ active: true });
  
  const stats: DocumentStats = {
    total: documents.length,
    policies: documents.filter(doc => doc.document_type === 'policy').length,
    onboardingDocs: documents.filter(doc => doc.document_type === 'onboarding').length,
    categories: {},
  };

  // Count documents by category
  documents.forEach(doc => {
    stats.categories[doc.category] = (stats.categories[doc.category] || 0) + 1;
  });

  return stats;
}

// User Management API Functions (Admin only)

export interface User {
  id: number;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  last_login?: string;
  created_at: string;
}

export async function getAllUsers(): Promise<User[]> {
  const response = await fetch(`${API_BASE_URL}/users`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch users: ${response.status} ${errorText}`);
  }

  const data = await response.json();
  return data.users;
}

export async function updateUser(id: number, updates: {
  first_name?: string;
  last_name?: string;
  email?: string;
  is_active?: boolean;
}): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/users/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: JSON.stringify(updates),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to update user: ${response.status} ${errorText}`);
  }

  const data = await response.json();
  return data.user;
}

export async function updateUserRole(id: number, role: string): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/users/${id}/role`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: JSON.stringify({ role }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to update user role: ${response.status} ${errorText}`);
  }

  const data = await response.json();
  return data.user;
}

export async function deleteUser(id: number): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE_URL}/users/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to delete user: ${response.status} ${errorText}`);
  }

  return response.json();
}

// Audit Log API Functions (Admin only)

export interface AuditLog {
  id: number;
  user_id: number;
  user?: User;
  action: string;
  resource_type: string;
  resource_id?: number;
  resource_name?: string;
  details?: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface AuditLogFilter {
  page?: number;
  limit?: number;
  action?: string;
  resource_type?: string;
  user_id?: string;
  from?: string;
  to?: string;
}

export interface AuditLogResponse {
  audit_logs: AuditLog[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    pages: number;
  };
}

export async function getAuditLogs(filters?: AuditLogFilter): Promise<AuditLogResponse> {
  const searchParams = new URLSearchParams();
  
  if (filters?.page) searchParams.append('page', filters.page.toString());
  if (filters?.limit) searchParams.append('limit', filters.limit.toString());
  if (filters?.action) searchParams.append('action', filters.action);
  if (filters?.resource_type) searchParams.append('resource_type', filters.resource_type);
  if (filters?.user_id) searchParams.append('user_id', filters.user_id);
  if (filters?.from) searchParams.append('from', filters.from);
  if (filters?.to) searchParams.append('to', filters.to);

  const url = `${API_BASE_URL}/audit-logs${searchParams.toString() ? `?${searchParams.toString()}` : ''}`;
  const response = await fetch(url, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to fetch audit logs: ${response.status} ${errorText}`);
  }

  return response.json();
}

// File Upload API functions

export interface FileUploadResponse {
  success: boolean;
  message: string;
  file_name?: string;
  file_type?: string;
  file_size?: number;
  file_path?: string;
  extracted_text?: string;
  error?: string;
}

export interface SupportedFileType {
  extension: string;
  mime_types: string[];
  description: string;
  max_size_mb: number;
}

export interface SupportedFileTypesResponse {
  supported_types: SupportedFileType[];
  max_file_size: string;
  accepted_extensions: string[];
}

// Upload file and extract text
export async function uploadFile(file: File): Promise<FileUploadResponse> {
  const formData = new FormData();
  formData.append('file', file);

  // Get auth token but don't include Content-Type header for file uploads
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  const headers: Record<string, string> = {};
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: 'POST',
    headers,
    body: formData,
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Upload failed');
  }

  return response.json();
}

// Get supported file types
export async function getSupportedFileTypes(): Promise<SupportedFileTypesResponse> {
  const response = await fetch(`${API_BASE_URL}/upload/supported-types`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to get supported file types: ${response.status} ${errorText}`);
  }

  return response.json();
} 