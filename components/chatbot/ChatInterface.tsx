"use client";

import { useState, useRef, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ChatMessage } from "./ChatMessage";
import { sendChatMessage } from "@/lib/api";
import { ChatMessage as ChatMessageType, ChatMode } from "@/lib/types";
import { Send, Shield, BookOpen, MessageCircle } from "lucide-react";

export function ChatInterface() {
  const [messages, setMessages] = useState<ChatMessageType[]>([]);
  const [currentMessage, setCurrentMessage] = useState("");
  const [chatMode, setChatMode] = useState<ChatMode | null>(null);
  const [isClient, setIsClient] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messageIdCounter = useRef(0);

  // Initialize client-side state after hydration
  useEffect(() => {
    setIsClient(true);
    setMessages([
      {
        id: "welcome-message",
        content:
          "Hello! I'm your IT Security Policy Assistant. I can help you with security onboarding or finding specific policy information. How can I assist you today?",
        role: "assistant",
        timestamp: new Date(),
        type: "general",
      },
    ]);
  }, []);

  const generateMessageId = () => {
    messageIdCounter.current += 1;
    return `msg-${messageIdCounter.current}-${Date.now()}`;
  };

  const chatMutation = useMutation({
    mutationFn: sendChatMessage,
    onSuccess: (response) => {
      const assistantMessage: ChatMessageType = {
        id: generateMessageId(),
        content: response.response,
        role: "assistant",
        timestamp: new Date(),
        type: (response.type as ChatMode) || "general",
        policyFiles: response.policy_files,
      };

      setMessages((prev) => [...prev, assistantMessage]);
    },
    onError: () => {
      const errorMessage: ChatMessageType = {
        id: generateMessageId(),
        content:
          "Sorry, I encountered an error. Please make sure the backend server is running and try again.",
        role: "assistant",
        timestamp: new Date(),
        type: "general",
      };
      setMessages((prev) => [...prev, errorMessage]);
    },
  });

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = async () => {
    if (!currentMessage.trim() || !chatMode) return;

    const userMessage: ChatMessageType = {
      id: generateMessageId(),
      content: currentMessage,
      role: "user",
      timestamp: new Date(),
      type: chatMode,
    };

    setMessages((prev) => [...prev, userMessage]);

    chatMutation.mutate({
      message: currentMessage,
      type: chatMode,
    });

    setCurrentMessage("");
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleModeSelection = (mode: ChatMode) => {
    setChatMode(mode);

    const modeMessage =
      mode === "onboarding"
        ? "Great! I'll help you with IT security onboarding. What specific topic would you like to learn about? (e.g., passwords, VPN, email security, data protection)"
        : "Perfect! I can help you search for specific security policies. What policy information are you looking for?";

    const assistantMessage: ChatMessageType = {
      id: generateMessageId(),
      content: modeMessage,
      role: "assistant",
      timestamp: new Date(),
      type: mode,
    };

    setMessages((prev) => [...prev, assistantMessage]);
  };

  // Don't render messages until client-side hydration is complete
  if (!isClient) {
    return (
      <div className="w-full max-w-7xl mx-auto flex flex-col h-full">
        <Card className="flex-1 flex flex-col h-full shadow-lg">
          <CardHeader className="flex-shrink-0">
            <CardTitle className="flex items-center gap-2">
              <Shield className="w-6 h-6 text-blue-600" />
              IT Security Policy Assistant
            </CardTitle>
          </CardHeader>
          <Separator />
          <CardContent className="flex-1 flex items-center justify-center">
            <div className="text-center text-muted-foreground">Loading...</div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="w-full max-w-7xl mx-auto flex flex-col h-full">
      <Card className="flex-1 flex flex-col h-full shadow-lg">
        <CardHeader className="flex-shrink-0">
          <CardTitle className="flex items-center gap-2">
            <Shield className="w-6 h-6 text-blue-600" />
            IT Security Policy Assistant
          </CardTitle>

          {/* Mode Selection */}
          {!chatMode && (
            <div className="flex gap-2 mt-4">
              <Button
                variant="outline"
                onClick={() => handleModeSelection("onboarding")}
                className="flex items-center gap-2"
              >
                <BookOpen className="w-4 h-4" />
                Security Onboarding
              </Button>
              <Button
                variant="outline"
                onClick={() => handleModeSelection("policy_search")}
                className="flex items-center gap-2"
              >
                <MessageCircle className="w-4 h-4" />
                Policy Search
              </Button>
            </div>
          )}

          {/* Current Mode Display */}
          {chatMode && (
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="flex items-center gap-1">
                {chatMode === "onboarding" ? (
                  <>
                    <BookOpen className="w-3 h-3" />
                    Security Onboarding
                  </>
                ) : (
                  <>
                    <MessageCircle className="w-3 h-3" />
                    Policy Search
                  </>
                )}
              </Badge>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setChatMode(null)}
              >
                Change Mode
              </Button>
            </div>
          )}
        </CardHeader>

        <Separator />

        {/* Messages Area */}
        <CardContent className="flex-1 overflow-y-auto p-4 min-h-0">
          <div className="space-y-4 pb-4">
            {messages.map((message) => (
              <ChatMessage key={message.id} message={message} />
            ))}
            <div ref={messagesEndRef} />
          </div>
        </CardContent>

        <Separator />

        {/* Input Area */}
        <CardContent className="flex-shrink-0 p-4">
          <div className="flex gap-2">
            <Input
              value={currentMessage}
              onChange={(e) => setCurrentMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder={
                chatMode
                  ? "Type your message..."
                  : "Please select a mode first..."
              }
              disabled={!chatMode || chatMutation.isPending}
              className="flex-1"
            />
            <Button
              onClick={handleSendMessage}
              disabled={
                !currentMessage.trim() || !chatMode || chatMutation.isPending
              }
              className="flex items-center gap-2"
            >
              <Send className="w-4 h-4" />
              {chatMutation.isPending ? "Sending..." : "Send"}
            </Button>
          </div>

          {chatMode && (
            <p className="text-xs text-muted-foreground mt-2">
              Mode:{" "}
              {chatMode === "onboarding"
                ? "Security Onboarding"
                : "Policy Search"}
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
