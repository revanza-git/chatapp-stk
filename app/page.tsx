import { ChatInterface } from "@/components/chatbot/ChatInterface";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Shield, Users, Search, FileText } from "lucide-react";

export default function Home() {
  return (
    <main className="h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 flex flex-col overflow-hidden">
      <div className="container mx-auto px-4 py-8 flex flex-col h-full">
        {/* Header */}
        <div className="text-center mb-8 flex-shrink-0">
          <div className="flex items-center justify-center gap-2 mb-4">
            <Shield className="w-8 h-8 text-blue-600" />
            <h1 className="text-3xl font-bold text-gray-900">
              IT Security Policy Chatbot
            </h1>
          </div>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Get instant help with IT security policies and onboarding. Ask
            questions, search policies, and learn about security best practices.
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-6 mb-8 flex-shrink-0">
          <Card>
            <CardHeader className="text-center">
              <Users className="w-8 h-8 text-blue-600 mx-auto mb-2" />
              <CardTitle className="text-lg">Security Onboarding</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>
                Get guided assistance with security awareness topics like
                passwords, VPN setup, email security, and data protection
                policies.
              </CardDescription>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="text-center">
              <Search className="w-8 h-8 text-green-600 mx-auto mb-2" />
              <CardTitle className="text-lg">Policy Search</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>
                Quickly find and access specific security policies, procedures,
                and guidelines relevant to your role and responsibilities.
              </CardDescription>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="text-center">
              <FileText className="w-8 h-8 text-purple-600 mx-auto mb-2" />
              <CardTitle className="text-lg">Instant Answers</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>
                Get immediate responses with relevant policy documents, best
                practices, and actionable security guidance.
              </CardDescription>
            </CardContent>
          </Card>
        </div>

        {/* Chat Interface */}
        <div className="flex justify-center flex-1 min-h-0">
          <ChatInterface />
        </div>

        {/* Footer */}
        <footer className="text-center mt-4 text-sm text-gray-500 flex-shrink-0">
          <p>
            Powered by AI • Always verify critical security information with
            your IT team
          </p>
        </footer>
      </div>
    </main>
  );
}
