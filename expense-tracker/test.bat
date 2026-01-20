@echo off
echo Testing Google Sheets connection...
echo.

REM Load environment variables from .env file
for /f "tokens=1,2 delims==" %%a in (.env) do (
    if not "%%a"=="" if not "%%a:~0,1"=="#" set %%a=%%b
)

echo Configuration:
echo - Sheets ID: %GOOGLE_SHEETS_ID%
echo - Service Account: %GOOGLE_SERVICE_ACCOUNT_EMAIL%
echo.

cd backend
echo Running connection test...
go run cmd/main.go
pause