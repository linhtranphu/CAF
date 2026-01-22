#!/bin/bash
set -e

echo "🚀 Starting Expense Tracker (Development)"

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo "Please copy .env.example to .env and configure your settings"
    exit 1
fi

# Load environment variables
source .env

# Validate required variables
if [ -z "$GEMINI_API_KEY" ] || [ "$GEMINI_API_KEY" = "your-gemini-api-key-here" ]; then
    echo "❌ Please set GEMINI_API_KEY in .env file"
    exit 1
fi

# Start services
echo "📦 Building and starting containers..."
docker-compose up -d --build

echo "⏳ Waiting for services to start..."
sleep 10

# Check health
echo "🔍 Checking service health..."
docker-compose ps

echo "✅ Services started successfully!"
echo ""
echo "🌐 Access URLs:"
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8081"
echo "Admin:    http://localhost:8081/admin"
echo ""
echo "📊 View logs: docker-compose logs -f"
echo "🛑 Stop services: docker-compose down"