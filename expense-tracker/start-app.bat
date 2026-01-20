@echo off
echo Starting Expense Tracker Application...
echo.

REM Check if .env file exists
if not exist "%~dp0.env" (
    echo ERROR: .env file not found!
    echo Looking for: %~dp0.env
    echo Please make sure .env file exists in the same directory as this batch file.
    pause
    exit /b 1
)

REM Load environment variables from .env file
for /f "usebackq tokens=1,2 delims==" %%a in ("%~dp0.env") do (
    if not "%%a"=="" if not "%%a:~0,1"=="#" set %%a=%%b
)

echo Configuration loaded:
echo - Google Sheets ID: %GOOGLE_SHEETS_ID%
echo - Service Account: %GOOGLE_SERVICE_ACCOUNT_EMAIL%
echo - Target Sheet: cost
echo.

REM Validate required environment variables
if "%GOOGLE_SHEETS_ID%"=="" (
    echo ERROR: GOOGLE_SHEETS_ID not found in .env file
    pause
    exit /b 1
)
if "%GOOGLE_SERVICE_ACCOUNT_EMAIL%"=="" (
    echo ERROR: GOOGLE_SERVICE_ACCOUNT_EMAIL not found in .env file
    pause
    exit /b 1
)

echo Checking and cleaning ports...
REM Kill processes using port 8081 (backend)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8081') do (
    taskkill /f /pid %%a 2>nul
)

REM Kill processes using port 3001 (frontend)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :3001') do (
    taskkill /f /pid %%a 2>nul
)

echo Starting Backend Server...
start "Backend" cmd /k "cd /d "%~dp0backend" && go mod tidy && go run cmd/main.go"

echo Waiting for backend to start...
timeout /t 5 /nobreak > nul

echo Starting Frontend Server...
start "Frontend" cmd /k "cd /d "%~dp0frontend" && echo Frontend running at http://localhost:3001 && python -m http.server 3001"

echo.
echo Both servers are starting...
echo Backend API: http://localhost:8081
echo Frontend App: http://localhost:3001
echo Data will be saved to sheet: cost
echo.
echo Press any key to close all servers...
pause > nul

echo Closing servers...
REM Kill only processes using our specific ports
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8081') do (
    taskkill /f /pid %%a 2>nul
)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :3001') do (
    taskkill /f /pid %%a 2>nul
)
echo Done.