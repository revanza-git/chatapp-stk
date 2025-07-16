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

# Pull Ollama model
echo "📥 Setting up Ollama model..."
docker-compose exec ollama ollama pull llama2 || echo "⚠️  Failed to pull model. You can do this manually later."

echo "✅ Setup complete!"
echo ""
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🗄️  Database: localhost:5433"
echo "🤖 Ollama: http://localhost:11434"
echo ""
echo "📊 To view logs: docker-compose logs -f"
echo "🛑 To stop: docker-compose down"
echo "🔄 To restart: docker-compose restart" 