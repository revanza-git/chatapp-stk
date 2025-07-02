import { ChatRequest, ChatResponse, PolicyFile } from './types';

// Use environment variable for API URL, fallback to localhost for development
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export async function sendChatMessage(request: ChatRequest): Promise<ChatResponse> {
  const response = await fetch(`${API_BASE_URL}/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to send chat message: ${response.status} ${errorText}`);
  }

  return response.json();
}

export async function getAllPolicies(): Promise<PolicyFile[]> {
  const response = await fetch(`${API_BASE_URL}/policies`);

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