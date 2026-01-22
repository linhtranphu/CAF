#!/bin/bash
set -e

echo "🔄 Restarting Expense Tracker..."

# Stop all containers
echo "⏹️  Stopping containers..."
docker-compose down

# Start all containers
echo "🚀 Starting containers..."
docker-compose up -d

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 10

# Check status
echo "🔍 Checking service status..."
docker-compose ps

echo ""
echo "✅ Restart completed!"
echo ""
echo "🌐 Access URLs:"
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8081"
echo "Admin:    http://localhost:8081/admin"
echo ""
echo "📊 View logs: docker-compose logs -f"