#!/bin/bash
set -e

echo "🚀 Deploying Expense Tracker..."

# Get GEMINI API Key
read -p "Enter GEMINI_API_KEY: " GEMINI_API_KEY
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required!"
    exit 1
fi

# Clone repository if not exists
if [ ! -d "CAF" ]; then
    git clone https://github.com/linhtranphu/CAF.git
fi

cd CAF/expense-tracker

# Pull latest code
echo "Pulling latest code..."
cd ..
git pull origin main
cd expense-tracker

# Stop existing containers
docker stop expense-mongodb expense-backend expense-frontend 2>/dev/null || true
docker rm expense-mongodb expense-backend expense-frontend 2>/dev/null || true

# Remove old images to force rebuild
docker rmi expense-backend expense-frontend 2>/dev/null || true

# Create network if not exists
docker network create expense-network 2>/dev/null || true

# Build images with no cache
echo "Building backend..."
docker build --no-cache -t expense-backend ./backend

echo "Building frontend..."
docker build --no-cache -t expense-frontend ./frontend

# Start MongoDB
echo "Starting MongoDB..."
docker run -d --name expense-mongodb \
  --network expense-network \
  -p 27017:27017 \
  public.ecr.aws/docker/library/mongo:7

# Wait for MongoDB
sleep 10

# Start backend
echo "Starting backend..."
docker run -d --name expense-backend \
  --network expense-network \
  -p 8081:8081 \
  -e PORT=8081 \
  -e MONGODB_URI=mongodb://expense-mongodb:27017 \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e SESSION_SECRET=expense-tracker-secret-$(date +%s) \
  --restart unless-stopped \
  expense-backend

# Wait for backend to be ready
echo "Waiting for backend to be ready..."
for i in {1..30}; do
  if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    echo "✅ Backend is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "❌ Backend failed to start, checking logs..."
    docker logs expense-backend --tail 10
    exit 1
  fi
  echo "Waiting for backend... ($i/30)"
  sleep 2
done

# Start frontend
echo "Starting frontend..."
docker run -d --name expense-frontend \
  --network expense-network \
  -p 3000:80 \
  --restart unless-stopped \
  expense-frontend

# Wait for frontend to be ready
echo "Waiting for frontend to be ready..."
for i in {1..15}; do
  if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend is ready!"
    break
  fi
  echo "Waiting for frontend... ($i/15)"
  sleep 2
done

echo "✅ Deployment complete!"
echo "Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "=== Health Checks ==="
echo "Backend health: $(curl -s http://localhost:8081/health 2>/dev/null || echo 'FAILED')"
echo "API health: $(curl -s http://localhost:8081/api/health 2>/dev/null || echo 'FAILED')"
echo "Frontend status: $(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 2>/dev/null || echo 'FAILED')"

echo ""
echo "=== Backend Logs ==="
docker logs expense-backend --tail 5

echo ""
echo "Access URLs:"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081"
echo "Admin:    http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081/admin"

echo ""
echo "=== Troubleshooting ==="
echo "If frontend shows 'Failed to fetch', check Security Group allows ports 3000 and 8081"
echo "Check logs: docker logs expense-backend"
echo "Check logs: docker logs expense-frontend"