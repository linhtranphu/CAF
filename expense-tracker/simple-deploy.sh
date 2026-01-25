#!/bin/bash
# Simple Deploy Script for Expense Tracker
# Uses pre-built images from Docker Hub

set -e

echo "🚀 Expense Tracker - Simple Deploy"
echo "=================================="

# Docker images
BACKEND_IMAGE="linhtranphu/expense-backend:latest"
FRONTEND_IMAGE="linhtranphu/expense-frontend:latest"

# Get GEMINI API Key (optional)
echo "🔑 GEMINI API Key (Optional)"
echo "============================"
echo "You can configure API key later via Settings page"
echo "Get your free API key at: https://aistudio.google.com/app/apikey"
echo ""
read -p "Enter your GEMINI_API_KEY (press Enter to skip): " GEMINI_API_KEY

if [ -z "$GEMINI_API_KEY" ]; then
    echo "⚠️  Skipping API Key - You can add it later in Settings"
    GEMINI_API_KEY=""
else
    echo "✅ API Key received"
fi

# Install Docker if needed
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker $USER
    echo "⚠️  Please logout and login again, then run this script again"
    exit 0
fi

# Check Docker permissions
if ! docker ps &> /dev/null; then
    echo "❌ Docker permission denied! Please logout and login again"
    exit 1
fi

# Create project directory
PROJECT_DIR="$HOME/expense-tracker"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

echo "📁 Project directory: $PROJECT_DIR"

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

echo "🐳 Starting deployment..."

# Pull images
docker pull mongo:7
docker pull "$BACKEND_IMAGE"
docker pull "$FRONTEND_IMAGE"

# Stop existing containers
docker-compose down 2>/dev/null || true

# Start services
docker-compose up -d

echo "⏳ Waiting for services..."

# Wait for backend (60 seconds max)
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
    sleep 1
done

# Wait for frontend (30 seconds max)
for i in {1..30}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend ready!"
        break
    fi
    sleep 1
done

# Get public IP
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")

echo ""
echo "🎉 Deployment Successful!"
echo "========================"
echo ""
echo "🌐 Access URLs:"
echo "Frontend:    http://$PUBLIC_IP:3000"
echo "Backend API: http://$PUBLIC_IP:8081"
echo "Admin Panel: http://$PUBLIC_IP:8081/admin"
echo "Settings:    http://$PUBLIC_IP:8081/settings"
echo ""
echo "📋 Next Steps:"
echo "1. Open http://$PUBLIC_IP:3000 in your browser"
echo "2. Register a new account"
if [ -z "$GEMINI_API_KEY" ]; then
    echo "3. Configure Gemini API key at Settings page"
    echo "4. Add expenses like: 'mua 2 cái bánh 50k'"
    echo "5. View reports at Admin Panel"
else
    echo "3. Add expenses like: 'mua 2 cái bánh 50k'"
    echo "4. View reports at Admin Panel"
fi
echo ""
echo "🛠️  Management Commands:"
echo "Restart:  docker-compose restart"
echo "Stop:     docker-compose down"
echo "Logs:     docker logs expense-backend"
echo "Update:   docker-compose pull && docker-compose up -d"
echo ""
echo "⚠️  Ensure Security Group allows ports 3000 and 8081"