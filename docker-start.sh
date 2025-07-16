#!/bin/bash

echo "🐳 Starting Chat App with Docker..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from example..."
    cp .env.example .env
    echo "⚠️  Please review and update the .env file with your configuration."
fi

# Start the services
echo "🚀 Starting all services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Skip Ollama model setup since AI is disabled
echo "ℹ️  AI features are currently disabled for faster deployment"
echo "💡 To enable AI: Edit docker-compose.yml and set AI_ENABLED=true"

echo "✅ Setup complete!"
echo ""
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🗄️  Database: localhost:5433"
echo "🤖 AI Status: DISABLED (for faster deployment)"
echo ""
echo "📊 To view logs: docker-compose logs -f"
echo "🛑 To stop: docker-compose down"
echo "🔄 To restart: docker-compose restart"
echo "🤖 To enable AI: See DOCKER.md for instructions" 