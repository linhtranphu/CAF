@echo off
echo Terminating Expense Tracker servers...

REM Kill backend (port 8081)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8081') do (
    echo Killing backend process %%a
    taskkill /f /pid %%a 2>nul
)

REM Kill frontend (port 3000)  
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :3000') do (
    echo Killing frontend process %%a
    taskkill /f /pid %%a 2>nul
)

REM Kill any remaining Go processes
taskkill /f /im go.exe 2>nul
taskkill /f /im main.exe 2>nul

REM Kill any remaining Node processes
taskkill /f /im node.exe 2>nul

echo All servers terminated.
pause