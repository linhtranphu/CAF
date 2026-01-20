@echo off
echo Starting Expense Tracker Backend...
cd backend
go mod tidy
go run cmd/main.go
pause