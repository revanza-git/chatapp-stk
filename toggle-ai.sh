#!/bin/bash

# AI Feature Toggle Script for Chat App

show_usage() {
    echo "🤖 AI Feature Toggle Script"
    echo ""
    echo "Usage: $0 [enable|disable|status]"
    echo ""
    echo "Commands:"
    echo "  enable   - Enable AI features (Ollama + full AI responses)"
    echo "  disable  - Disable AI features (Mock responses only)"
    echo "  status   - Show current AI feature status"
    echo ""
    echo "Examples:"
    echo "  $0 enable   # Enable AI features"
    echo "  $0 disable  # Disable AI features"
    echo "  $0 status   # Check current status"
}

check_status() {
    if grep -q "AI_ENABLED=true" docker-compose.yml; then
        echo "🤖 AI Features: ENABLED"
        echo "   - Ollama service: Active"
        echo "   - AI responses: Enabled"
    else
        echo "🤖 AI Features: DISABLED"
        echo "   - Ollama service: Commented out"
        echo "   - AI responses: Mock responses only"
    fi
}

enable_ai() {
    echo "🚀 Enabling AI features..."
    
    # Update AI_ENABLED to true
    sed -i 's/AI_ENABLED=false/AI_ENABLED=true/' docker-compose.yml
    sed -i 's/AI_ENABLED=false/AI_ENABLED=true/' .env.example
    
    # Uncomment Ollama service
    sed -i 's/^  # ollama:/  ollama:/' docker-compose.yml
    sed -i 's/^  #   /      /' docker-compose.yml
    
    # Add ollama dependency back to backend
    sed -i 's/# Removed ollama dependency since AI is disabled/ollama:\n        condition: service_started/' docker-compose.yml
    
    echo "✅ AI features enabled!"
    echo "💡 Run 'docker-compose up --build -d' to apply changes"
    echo "📥 After startup, run: docker-compose exec ollama ollama pull llama2"
}

disable_ai() {
    echo "🛑 Disabling AI features..."
    
    # Update AI_ENABLED to false
    sed -i 's/AI_ENABLED=true/AI_ENABLED=false/' docker-compose.yml
    sed -i 's/AI_ENABLED=true/AI_ENABLED=false/' .env.example
    
    # Comment out Ollama service
    sed -i 's/^  ollama:/  # ollama:/' docker-compose.yml
    sed -i 's/^      /  #   /' docker-compose.yml
    
    # Remove ollama dependency from backend
    sed -i 's/ollama:\n        condition: service_started/# Removed ollama dependency since AI is disabled/' docker-compose.yml
    
    echo "✅ AI features disabled!"
    echo "💡 Run 'docker-compose up --build -d' to apply changes"
}

# Main script logic
case "$1" in
    "enable")
        enable_ai
        ;;
    "disable")
        disable_ai
        ;;
    "status")
        check_status
        ;;
    *)
        show_usage
        exit 1
        ;;
esac 