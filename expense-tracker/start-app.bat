@echo off
setlocal enabledelayedexpansion
echo Starting Expense Tracker Application...
echo.

REM Check if .env file exists
if not exist "%~dp0.env" (
    echo WARNING: .env file not found!
    echo Copy .env.example to .env and add your API keys.
    echo Continuing with fallback parsing...
    echo.
)

echo Starting Backend Server...
start "Backend" cmd /k "cd /d "%~dp0backend" && go mod tidy && go run cmd/main.go"

echo Waiting for backend to start...
timeout /t 3 /nobreak > nul

echo Starting Frontend Server...
start "Frontend" cmd /k "cd /d "%~dp0frontend" && npm install && npm run dev"

echo.
echo Both servers are starting...
echo Backend API: http://localhost:8081
echo Frontend App: http://localhost:3000
echo Database: MongoDB (localhost:27017)
echo.
echo Press Ctrl+C or close this window to terminate all servers...

REM Wait for user interrupt
:wait
timeout /t 1 /nobreak > nul
goto wait