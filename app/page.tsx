'use client'

import { ChatInterface } from "@/components/chatbot/ChatInterface";
import { Shield, Settings } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import ProtectedRoute from "@/components/auth/ProtectedRoute";
import UserProfile from "@/components/auth/UserProfile";

export default function Home() {
  return (
    <ProtectedRoute>
      <HomeContent />
    </ProtectedRoute>
  );
}

function HomeContent() {
  return (
    <main className="h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex flex-col overflow-hidden">
      <div className="flex flex-col h-full px-4 py-4">
        {/* Compact Header */}
        <div className="flex-shrink-0 mb-4">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              <Shield className="w-6 h-6 text-blue-600" />
              <h1 className="text-2xl font-bold text-gray-900">
                IT Security Policy Chatbot
              </h1>
            </div>
            <div className="flex items-center gap-2">
              <UserProfile variant="compact" />
              <Link href="/dashboard">
                <Button variant="outline" size="sm">
                  <Settings className="w-4 h-4 mr-2" />
                  Dashboard
                </Button>
              </Link>
            </div>
          </div>
          <p className="text-sm text-gray-600 text-center">
            Get instant help with IT security policies and onboarding
          </p>
        </div>

        {/* Chat Interface - Now takes most of the screen */}
        <div className="flex-1 min-h-0 px-2">
          <ChatInterface />
        </div>

        {/* Compact Footer */}
        <footer className="text-center mt-2 text-xs text-gray-500 flex-shrink-0">
          <p>
            Powered by AI • Always verify critical security information with your IT team
          </p>
        </footer>
      </div>
    </main>
  );
}


