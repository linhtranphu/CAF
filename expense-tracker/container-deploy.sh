#!/bin/bash
# Modern Container Deploy for Expense Tracker
# Ubuntu 22.04 LTS on AWS EC2

set -e

echo "🚀 Expense Tracker - Container Deploy"
echo "===================================="

# Get GEMINI API Key
if [ -z "$GEMINI_API_KEY" ]; then
    echo "🔑 Cần GEMINI API Key từ: https://makersuite.google.com/app/apikey"
    read -p "Nhập GEMINI_API_KEY: " GEMINI_API_KEY
fi

# Install Docker if needed
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    
    echo "⚠️  Docker installed. Please logout and login again, then run this script again"
    exit 0
fi

# Check Docker permissions
if ! docker ps &> /dev/null; then
    echo "❌ Docker permission denied! Please logout and login again"
    exit 1
fi

echo "✅ Docker ready"

# Create project directory
mkdir -p ~/expense-tracker
cd ~/expense-tracker

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
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    image: linhtranphu/expense-backend:latest
    container_name: expense-backend
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - SESSION_SECRET=$(openssl rand -hex 32)
    depends_on:
      mongodb:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: curl -f http://localhost:8081/health || exit 1
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: linhtranphu/expense-frontend:latest
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped
    healthcheck:
      test: curl -f http://localhost:80 || exit 1
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mongodb_data:

networks:
  default:
    name: expense-network
EOF

# Create .env file
cat > .env << EOF
GEMINI_API_KEY=${GEMINI_API_KEY}
SESSION_SECRET=$(openssl rand -hex 32)
PORT=8081
MONGODB_URI=mongodb://mongodb:27017
EOF

echo "🐳 Starting containers..."

# Pull latest images
docker pull mongo:7
docker pull linhtranphu/expense-backend:latest
docker pull linhtranphu/expense-frontend:latest

# Stop existing containers
docker-compose down 2>/dev/null || true

# Start services
docker-compose up -d

echo "⏳ Waiting for services to be ready..."

# Wait for backend
for i in {1..60}; do
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo "✅ Backend ready!"
        break
    fi
    if [ $i -eq 60 ]; then
        echo "❌ Backend failed to start"
        docker logs expense-backend --tail 20
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

# Health check
BACKEND_HEALTH=$(curl -s http://localhost:8081/health 2>/dev/null || echo "FAILED")
FRONTEND_STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 2>/dev/null || echo "FAILED")

echo ""
echo "🎉 Deploy Complete!"
echo "=================="
echo "Frontend:    http://$PUBLIC_IP:3000"
echo "Backend API: http://$PUBLIC_IP:8081"
echo "Admin Panel: http://$PUBLIC_IP:8081/admin"
echo ""
echo "📊 Health Status:"
echo "Backend: $BACKEND_HEALTH"
echo "Frontend: $FRONTEND_STATUS"
echo ""
echo "🐳 Running Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""
echo "📝 Useful Commands:"
echo "View logs:    docker logs expense-backend"
echo "Restart:      docker-compose restart"
echo "Stop all:     docker-compose down"
echo "Update:       docker-compose pull && docker-compose up -d"
echo ""
echo "⚠️  Make sure Security Group allows ports 3000 and 8081"