import { ChatMessage as ChatMessageType } from "@/lib/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Bot, User, FileText, Calendar } from "lucide-react";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";

interface ChatMessageProps {
  message: ChatMessageType;
}

export function ChatMessage({ message }: ChatMessageProps) {
  const [isClient, setIsClient] = useState(false);
  const isUser = message.role === "user";

  useEffect(() => {
    setIsClient(true);
  }, []);

  return (
    <div className={cn("flex gap-3", isUser ? "justify-end" : "justify-start")}>
      {!isUser && (
        <div className="flex-shrink-0">
          <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center">
            <Bot className="w-4 h-4 text-primary-foreground" />
          </div>
        </div>
      )}

      <div className={cn("max-w-[80%]", isUser && "order-first")}>
        <Card
          className={cn(
            isUser ? "bg-primary text-primary-foreground" : "bg-muted"
          )}
        >
          <CardContent className={cn("p-3", isUser && "py-0 px-3")}>
            <p className="text-sm">{message.content}</p>

            {message.type && !isUser && (
              <Badge variant="default" className="mt-2">
                {message.type}
              </Badge>
            )}
          </CardContent>
        </Card>

        {/* Policy Files Display */}
        {message.policyFiles && message.policyFiles.length > 0 && (
          <div className="mt-3 space-y-2">
            {message.policyFiles.map((policy, index) => (
              <Card key={index} className="border-l-4 border-l-blue-500">
                <CardContent className="p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <FileText className="w-4 h-4 text-blue-600" />
                    <h4 className="font-medium text-sm">{policy.name}</h4>
                    <Badge variant="outline" className="text-xs">
                      {policy.category}
                    </Badge>
                  </div>
                  <p className="text-sm text-muted-foreground mb-2">
                    {policy.content}
                  </p>
                  <div className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Calendar className="w-3 h-3" />
                    Last updated: {policy.last_updated}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {isClient && (
          <p className="text-xs text-muted-foreground mt-1">
            {message.timestamp.toLocaleTimeString()}
          </p>
        )}
      </div>

      {isUser && (
        <div className="flex-shrink-0">
          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
            <User className="w-4 h-4" />
          </div>
        </div>
      )}
    </div>
  );
}
