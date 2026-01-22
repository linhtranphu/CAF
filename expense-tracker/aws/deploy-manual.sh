#!/bin/bash
set -e

echo "🚀 Manual Docker Deployment..."

# Get GEMINI API Key
read -p "Enter GEMINI_API_KEY: " GEMINI_API_KEY
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required!"
    exit 1
fi

# Stop existing containers
docker stop expense-mongodb expense-backend expense-frontend 2>/dev/null || true
docker rm expense-mongodb expense-backend expense-frontend 2>/dev/null || true

# Start MongoDB
echo "Starting MongoDB..."
docker run -d --name expense-mongodb \
  -p 27017:27017 \
  -v mongodb_data:/data/db \
  public.ecr.aws/docker/library/mongo:7

# Wait for MongoDB
sleep 10

# Build and start Backend
echo "Building Backend..."
cd backend
docker build -t expense-backend .
docker run -d --name expense-backend \
  -p 8081:8081 \
  -e PORT=8081 \
  -e MONGODB_URI=mongodb://$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' expense-mongodb):27017 \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e SESSION_SECRET=expense-tracker-secret-$(date +%s) \
  expense-backend

# Build and start Frontend
echo "Building Frontend..."
cd ../frontend
docker build -t expense-frontend .
docker run -d --name expense-frontend \
  -p 3000:80 \
  expense-frontend

echo "✅ Manual deployment complete!"
echo "Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "Access URLs:"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081"