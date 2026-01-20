@echo off
echo ========================================
echo    EXPENSE TRACKER QUICK SETUP
echo ========================================
echo.

REM Check Go
go version >nul 2>&1
if %errorlevel% == 0 (
    echo ✅ Go installed
    goto :check_node
)

echo ❌ Go not found
echo Please install Go manually from: https://golang.org/dl/
echo Choose: go1.21.5.windows-amd64.msi
echo Then run this script again
pause
exit /b 1

:check_node
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js not found
    echo Please install from: https://nodejs.org/
    pause
    exit /b 1
)
echo ✅ Node.js installed

echo.
echo Installing dependencies...
cd backend
go mod tidy
cd ..\frontend
npm install
cd ..

echo.
echo ✅ Setup complete! Run start-app.bat
pause