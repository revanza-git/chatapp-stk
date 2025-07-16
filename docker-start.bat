@echo off
echo 🐳 Starting Chat App with Docker...

REM Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker is not running. Please start Docker first.
    pause
    exit /b 1
)

REM Create .env file if it doesn't exist
if not exist .env (
    echo 📝 Creating .env file from example...
    copy .env.example .env
    echo ⚠️  Please review and update the .env file with your configuration.
)

REM Start the services
echo 🚀 Starting all services...
docker-compose up --build -d

REM Wait for services to be ready
echo ⏳ Waiting for services to start...
timeout /t 10 /nobreak >nul

REM Pull Ollama model
echo 📥 Setting up Ollama model...
docker-compose exec ollama ollama pull llama2 || echo ⚠️  Failed to pull model. You can do this manually later.

echo ✅ Setup complete!
echo.
echo 🌐 Frontend: http://localhost:3000
echo 🔧 Backend API: http://localhost:8080
echo 🗄️  Database: localhost:5433
echo 🤖 Ollama: http://localhost:11434
echo.
echo 📊 To view logs: docker-compose logs -f
echo 🛑 To stop: docker-compose down
echo 🔄 To restart: docker-compose restart
pause 