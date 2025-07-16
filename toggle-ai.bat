@echo off
REM AI Feature Toggle Script for Chat App (Windows)

if "%1"=="enable" goto enable_ai
if "%1"=="disable" goto disable_ai
if "%1"=="status" goto check_status
goto show_usage

:show_usage
echo 🤖 AI Feature Toggle Script
echo.
echo Usage: %0 [enable^|disable^|status]
echo.
echo Commands:
echo   enable   - Enable AI features (Ollama + full AI responses)
echo   disable  - Disable AI features (Mock responses only)
echo   status   - Show current AI feature status
echo.
echo Examples:
echo   %0 enable   # Enable AI features
echo   %0 disable  # Disable AI features
echo   %0 status   # Check current status
goto end

:check_status
findstr /C:"AI_ENABLED=true" docker-compose.yml >nul
if %errorlevel%==0 (
    echo 🤖 AI Features: ENABLED
    echo    - Ollama service: Active
    echo    - AI responses: Enabled
) else (
    echo 🤖 AI Features: DISABLED
    echo    - Ollama service: Commented out
    echo    - AI responses: Mock responses only
)
goto end

:enable_ai
echo 🚀 Enabling AI features...

REM Update AI_ENABLED to true
powershell -Command "(Get-Content docker-compose.yml) -replace 'AI_ENABLED=false', 'AI_ENABLED=true' | Set-Content docker-compose.yml"
powershell -Command "(Get-Content .env.example) -replace 'AI_ENABLED=false', 'AI_ENABLED=true' | Set-Content .env.example"

echo ✅ AI features enabled!
echo 💡 Manually uncomment the ollama service in docker-compose.yml
echo 💡 Run 'docker-compose up --build -d' to apply changes
echo 📥 After startup, run: docker-compose exec ollama ollama pull llama2
goto end

:disable_ai
echo 🛑 Disabling AI features...

REM Update AI_ENABLED to false
powershell -Command "(Get-Content docker-compose.yml) -replace 'AI_ENABLED=true', 'AI_ENABLED=false' | Set-Content docker-compose.yml"
powershell -Command "(Get-Content .env.example) -replace 'AI_ENABLED=true', 'AI_ENABLED=false' | Set-Content .env.example"

echo ✅ AI features disabled!
echo 💡 Manually comment out the ollama service in docker-compose.yml
echo 💡 Run 'docker-compose up --build -d' to apply changes
goto end

:end
pause 