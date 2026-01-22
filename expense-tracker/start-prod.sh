#!/bin/bash
set -e

echo "🚀 Starting Expense Tracker (Production)"

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    exit 1
fi

# Load environment variables
source .env

# Validate required variables
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required"
    exit 1
fi

if [ -z "$SESSION_SECRET" ] || [ "$SESSION_SECRET" = "expense-tracker-secret-change-in-production" ]; then
    echo "❌ Please set a secure SESSION_SECRET in production"
    exit 1
fi

# Start production services
echo "📦 Building and starting production containers..."
docker-compose -f docker-compose.prod.yml up -d --build

echo "⏳ Waiting for services to start..."
sleep 15

# Check health
echo "🔍 Checking service health..."
docker-compose -f docker-compose.prod.yml ps

echo "✅ Production services started successfully!"
echo ""
echo "🌐 Access URLs:"
echo "Frontend: http://localhost"
echo "Backend:  Internal only"
echo ""
echo "📊 View logs: docker-compose -f docker-compose.prod.yml logs -f"
echo "🛑 Stop services: docker-compose -f docker-compose.prod.yml down"