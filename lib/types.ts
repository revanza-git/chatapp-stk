export interface ChatMessage {
  id: string;
  content: string;
  role: 'user' | 'assistant';
  timestamp: Date;
  type?: 'onboarding' | 'policy_search' | 'general';
  policyFiles?: PolicyFile[];
}

export interface PolicyFile {
  name: string;
  content: string;
  category: string;
  last_updated: string;
}

export interface ChatRequest {
  message: string;
  type: 'onboarding' | 'policy_search';
}

export interface ChatResponse {
  response: string;
  type: string;
  policy_files?: PolicyFile[];
}

export type ChatMode = 'onboarding' | 'policy_search'; 