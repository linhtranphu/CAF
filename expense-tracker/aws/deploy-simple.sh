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

# Stop existing containers
docker stop expense-mongodb expense-backend expense-frontend 2>/dev/null || true
docker rm expense-mongodb expense-backend expense-frontend 2>/dev/null || true

# Build images
echo "Building backend..."
docker build -t expense-backend ./backend

echo "Building frontend..."
docker build -t expense-frontend ./frontend

# Start MongoDB
echo "Starting MongoDB..."
docker run -d --name expense-mongodb -p 27017:27017 public.ecr.aws/docker/library/mongo:7

# Wait for MongoDB
sleep 10

# Start backend
echo "Starting backend..."
docker run -d --name expense-backend -p 8081:8081 \
  -e PORT=8081 \
  -e MONGODB_URI=mongodb://172.17.0.1:27017 \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e SESSION_SECRET=expense-tracker-secret-$(date +%s) \
  expense-backend

# Start frontend
echo "Starting frontend..."
docker run -d --name expense-frontend -p 3000:80 expense-frontend

echo "✅ Deployment complete!"
echo "Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "Access URLs:"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081"