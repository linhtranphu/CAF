#!/bin/bash
# Interactive Container Deploy for Expense Tracker

set -e

echo "🚀 Expense Tracker - Interactive Deploy"
echo "======================================="

# Check if running in terminal
if [ ! -t 0 ]; then
    echo "❌ Script này cần chạy trực tiếp, không qua pipe"
    echo "📥 Download và chạy:"
    echo "curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/interactive-deploy.sh"
    echo "chmod +x interactive-deploy.sh"
    echo "./interactive-deploy.sh"
    exit 1
fi

# Get GEMINI API Key
echo "🔑 GEMINI API Key Setup"
echo "======================="
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
    echo ""
    echo "📦 Docker Installation"
    echo "====================="
    read -p "Docker chưa được cài đặt. Cài đặt ngay? (y/N): " install_docker
    
    if [[ $install_docker =~ ^[Yy]$ ]]; then
        echo "Đang cài đặt Docker..."
        curl -fsSL https://get.docker.com | sudo sh
        sudo usermod -aG docker $USER
        
        echo "✅ Docker đã cài đặt"
        echo "⚠️  Vui lòng logout và login lại, sau đó chạy script này lần nữa"
        exit 0
    else
        echo "❌ Cần Docker để tiếp tục"
        exit 1
    fi
fi

# Check Docker permissions
if ! docker ps &> /dev/null; then
    echo "❌ Docker permission denied!"
    echo "Chạy: sudo usermod -aG docker $USER"
    echo "Sau đó logout và login lại"
    exit 1
fi

echo "✅ Docker ready"

# Create project directory
PROJECT_DIR="$HOME/expense-tracker"
echo ""
echo "📁 Project Setup"
echo "==============="
echo "Tạo project tại: $PROJECT_DIR"

mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

# Create docker-compose.yml
echo "Tạo docker-compose.yml..."
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
      - mongodb
    restart: unless-stopped

  frontend:
    image: linhtranphu/expense-frontend:latest
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

echo "✅ Configuration created"

# Deploy confirmation
echo ""
echo "🚀 Ready to Deploy"
echo "=================="
echo "Services sẽ được deploy:"
echo "- MongoDB (port 27017)"
echo "- Backend API (port 8081)"
echo "- Frontend Web (port 3000)"
echo ""
read -p "Tiếp tục deploy? (Y/n): " confirm_deploy

if [[ $confirm_deploy =~ ^[Nn]$ ]]; then
    echo "❌ Deploy bị hủy"
    exit 0
fi

echo ""
echo "🐳 Starting deployment..."

# Pull images
echo "Downloading Docker images..."
docker pull mongo:7
docker pull linhtranphu/expense-backend:latest
docker pull linhtranphu/expense-frontend:latest

# Stop existing containers
docker-compose down 2>/dev/null || true

# Start services
echo "Starting services..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services..."

# Wait for backend
for i in {1..60}; do
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo "✅ Backend ready!"
        break
    fi
    if [ $i -eq 60 ]; then
        echo "❌ Backend timeout"
        echo "Logs:"
        docker logs expense-backend --tail 10
        exit 1
    fi
    printf "."
    sleep 2
done

# Wait for frontend
for i in {1..30}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend ready!"
        break
    fi
    printf "."
    sleep 2
done

# Get public IP
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")

echo ""
echo "🎉 Deploy Successful!"
echo "===================="
echo ""
echo "🌐 Access URLs:"
echo "Frontend:    http://$PUBLIC_IP:3000"
echo "Backend API: http://$PUBLIC_IP:8081"
echo "Admin Panel: http://$PUBLIC_IP:8081/admin"
echo ""
echo "📊 Health Check:"
BACKEND_HEALTH=$(curl -s http://localhost:8081/health 2>/dev/null || echo "FAILED")
FRONTEND_STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 2>/dev/null || echo "FAILED")
echo "Backend: $BACKEND_HEALTH"
echo "Frontend: HTTP $FRONTEND_STATUS"
echo ""
echo "🐳 Running Containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""
echo "📝 Management Commands:"
echo "View logs:    docker logs expense-backend"
echo "Restart:      docker-compose restart"
echo "Stop all:     docker-compose down"
echo "Update:       docker-compose pull && docker-compose up -d"
echo ""
echo "⚠️  Đảm bảo Security Group mở port 3000 và 8081"
echo ""
echo "🎯 Next Steps:"
echo "1. Mở http://$PUBLIC_IP:3000 trong browser"
echo "2. Đăng ký tài khoản mới"
echo "3. Thêm chi phí: 'ăn trưa 50k'"
echo "4. Xem báo cáo tại Admin Panel"