#!/bin/bash
set -e

echo "🐳 Docker Compose Deployment..."

# Check if .env.prod exists
if [ ! -f ".env.prod" ]; then
    echo "❌ .env.prod file not found!"
    echo "Please create .env.prod with GEMINI_API_KEY"
    exit 1
fi

# Load environment variables
export $(cat .env.prod | grep -v '^#' | xargs)

# Validate required variables
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required in .env.prod!"
    exit 1
fi

# Generate random secret if not set
if [ -z "$RANDOM_SECRET" ]; then
    export RANDOM_SECRET=$(date +%s | sha256sum | head -c 32)
    echo "Generated random secret: $RANDOM_SECRET"
fi

# Pull latest code
echo "📥 Pulling latest code..."
git pull origin main

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down

# Remove old images to force rebuild
echo "🗑️ Removing old images..."
docker rmi expense-tracker-backend expense-tracker-frontend 2>/dev/null || true

# Build and start services
echo "🔨 Building and starting services..."
docker-compose -f docker-compose.prod.yml up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services..."
sleep 15

# Health checks
echo ""
echo "=== Health Checks ==="
echo "Backend: $(curl -s http://localhost:8081/health | jq -r '.status' 2>/dev/null || echo 'FAILED')"
echo "Frontend: $(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 2>/dev/null || echo 'FAILED')"

# Show status
echo ""
echo "=== Container Status ==="
docker-compose -f docker-compose.prod.yml ps

# Show logs
echo ""
echo "=== Recent Logs ==="
docker-compose -f docker-compose.prod.yml logs --tail=5

echo ""
echo "✅ Deployment complete!"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo 'localhost'):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo 'localhost'):8081"