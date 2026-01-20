@echo off
echo Setting up Expense Tracker for first time...
echo.

echo 1. Installing Go dependencies...
cd backend
go mod tidy
if %errorlevel% neq 0 (
    echo Error: Failed to install Go dependencies
    pause
    exit /b 1
)

echo.
echo 2. Checking Python installation...
python --version
if %errorlevel% neq 0 (
    echo Error: Python is not installed or not in PATH
    echo Please install Python from https://python.org
    pause
    exit /b 1
)

echo.
echo 3. Setup complete!
echo.
echo Configuration loaded from .env file:
echo - Google Sheets ID: %GOOGLE_SHEETS_ID%
echo - Service Account: %GOOGLE_SERVICE_ACCOUNT_EMAIL%
echo.
echo Ready to run! Use start-app.bat to launch the application
echo.
pause