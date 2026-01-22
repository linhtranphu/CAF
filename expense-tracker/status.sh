#!/bin/bash

echo "🔍 Expense Tracker Status"
echo "========================="

# Check container status
docker-compose ps

echo ""
echo "🌐 Access URLs:"
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8081"
echo "Admin:    http://localhost:8081/admin"

echo ""
echo "📊 Quick health check:"
echo -n "Backend API: "
if curl -s http://localhost:8081 > /dev/null 2>&1; then
    echo "✅ Running"
else
    echo "❌ Not responding"
fi

echo -n "Frontend: "
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Running"
else
    echo "❌ Not responding"
fi