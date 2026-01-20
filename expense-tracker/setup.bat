@echo off
echo ========================================
echo    EXPENSE TRACKER SETUP
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% == 0 (
    echo ✅ Go is already installed
    go version
    goto :check_node
)

echo ❌ Go not found. Installing Go...
echo Downloading Go 1.21.5...
powershell -Command "Invoke-WebRequest -Uri 'https://go.dev/dl/go1.21.5.windows-amd64.msi' -OutFile 'go-installer.msi'"

if not exist "go-installer.msi" (
    echo ❌ Download failed. Please install Go manually from: https://golang.org/dl/
    pause
    exit /b 1
)

echo Installing Go...
msiexec /i go-installer.msi /quiet /norestart
del go-installer.msi

echo ⚠️  Please restart command prompt and run setup.bat again
pause
exit /b 0

:check_node
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js not found. Please install from: https://nodejs.org/
    pause
    exit /b 1
)
echo ✅ Node.js is installed

echo.
echo Installing dependencies...
cd backend
set CGO_ENABLED=1
go mod tidy
if %errorlevel% neq 0 (
    echo ❌ Failed to install Go dependencies
    pause
    exit /b 1
)

cd .\frontend
npm install
if %errorlevel% neq 0 (
    echo ❌ Failed to install Node.js dependencies
    pause
    exit /b 1
)

cd ..
echo.
echo ✅ Setup completed! Run start-app.bat to launch
pause