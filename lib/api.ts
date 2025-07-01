import { ChatRequest, ChatResponse, PolicyFile } from './types';

const API_BASE_URL = 'http://localhost:8080/api';

export async function sendChatMessage(request: ChatRequest): Promise<ChatResponse> {
  const response = await fetch(`${API_BASE_URL}/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error('Failed to send chat message');
  }

  return response.json();
}

export async function getAllPolicies(): Promise<PolicyFile[]> {
  const response = await fetch(`${API_BASE_URL}/policies`);

  if (!response.ok) {
    throw new Error('Failed to fetch policies');
  }

  return response.json();
}

export async function checkHealth(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE_URL}/health`);

  if (!response.ok) {
    throw new Error('Backend health check failed');
  }

  return response.json();
} 