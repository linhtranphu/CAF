#!/bin/bash
# Deploy using Docker Hub images

set -e

echo "🚀 Expense Tracker - Docker Hub Deploy"
echo "======================================"

# Config - Sử dụng images đã build
BACKEND_IMAGE="linhtranphu/expense-backend:latest"
FRONTEND_IMAGE="linhtranphu/expense-frontend:latest"

# Get GEMINI API Key - Always prompt for input
echo "🔑 GEMINI API Key Setup"
echo "====================="
echo "Bạn cần GEMINI API Key để sử dụng tính năng AI parsing"
echo "Lấy miễn phí tại: https://makersuite.google.com/app/apikey"
echo ""

while [ -z "$GEMINI_API_KEY" ]; do
    read -p "Nhập GEMINI_API_KEY của bạn: " GEMINI_API_KEY
    if [ -z "$GEMINI_API_KEY" ]; then
        echo "⚠️  API Key không được để trống!"
    fi
done

echo "✅ API Key đã nhận"

# Install Docker if needed
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker $USER
    echo "⚠️  Logout và login lại, sau đó chạy script này lần nữa"
    exit 0
fi

if ! docker ps &> /dev/null; then
    echo "❌ Docker permission denied! Logout và login lại"
    exit 1
fi

echo "✅ Docker ready"

# Setup project
PROJECT_DIR="$HOME/expense-tracker"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

# Create docker-compose.yml
cat > docker-compose.yml << EOF
version: '3.8'

services:
  mongodb:
    image: mongo:7
    container_name: expense-mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    restart: unless-stopped

  backend:
    image: $BACKEND_IMAGE
    container_name: expense-backend
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=$GEMINI_API_KEY
      - SESSION_SECRET=$(openssl rand -hex 32)
    depends_on:
      - mongodb
    restart: unless-stopped

  frontend:
    image: $FRONTEND_IMAGE
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

echo "🐳 Starting services..."

# Pull images
docker pull mongo:7
docker pull "$BACKEND_IMAGE"
docker pull "$FRONTEND_IMAGE"

# Stop existing
docker-compose down 2>/dev/null || true

# Start services
docker-compose up -d

echo "⏳ Waiting for services..."

# Wait for backend
for i in {1..60}; do
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo "✅ Backend ready!"
        break
    fi
    if [ $i -eq 60 ]; then
        echo "❌ Backend timeout"
        docker logs expense-backend --tail 10
        exit 1
    fi
    sleep 2
done

# Wait for frontend
for i in {1..30}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend ready!"
        break
    fi
    sleep 2
done

# Get public IP
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")

echo ""
echo "🎉 Deploy Complete!"
echo "=================="
echo "Frontend:    http://$PUBLIC_IP:3000"
echo "Backend API: http://$PUBLIC_IP:8081"
echo "Admin Panel: http://$PUBLIC_IP:8081/admin"
echo ""
echo "🐳 Running Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"