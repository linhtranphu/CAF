#!/bin/bash
set -e

echo "⏹️  Stopping Expense Tracker..."

# Stop all containers
docker-compose down

echo "✅ All services stopped!"
echo ""
echo "🚀 To start again: ./start-dev.sh or ./restart.sh"