@echo off
echo Starting Expense Tracker Frontend...
cd frontend
echo Frontend is running at http://localhost:3000
python -m http.server 3000
pause