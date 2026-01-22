#!/bin/bash
set -e

echo "🚀 Expense Tracker - EC2 Deploy Script"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running on EC2
if ! curl -s --max-time 3 http://169.254.169.254/latest/meta-data/instance-id > /dev/null 2>&1; then
    print_warning "Not running on EC2, using localhost for public IP"
    PUBLIC_IP="localhost"
else
    PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
fi

# Check Docker installation
if ! command -v docker &> /dev/null; then
    print_error "Docker not found! Please run setup.sh first"
    echo "curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/setup.sh"
    echo "chmod +x setup.sh && ./setup.sh"
    exit 1
fi

# Check Docker permissions
if ! docker ps &> /dev/null; then
    print_error "Docker permission denied! Please logout and login again after setup.sh"
    exit 1
fi

# Get GEMINI API Key
if [ -z "$GEMINI_API_KEY" ]; then
    echo ""
    echo "🔑 Bạn cần GEMINI API Key để sử dụng AI parsing"
    echo "   Lấy tại: https://makersuite.google.com/app/apikey"
    echo ""
    read -p "Nhập GEMINI_API_KEY: " GEMINI_API_KEY
    
    if [ -z "$GEMINI_API_KEY" ]; then
        print_error "GEMINI_API_KEY là bắt buộc!"
        exit 1
    fi
fi

# Setup project directory
PROJECT_DIR="$HOME/expense-tracker"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

print_status "Setting up project directory: $PROJECT_DIR"

# Clone or update repository
if [ ! -d "CAF" ]; then
    print_status "Cloning repository..."
    git clone https://github.com/linhtranphu/CAF.git
else
    print_status "Updating repository..."
    cd CAF
    git pull origin main
    cd ..
fi

cd CAF/expense-tracker

# Stop existing containers
print_status "Stopping existing containers..."
docker stop expense-mongodb expense-backend expense-frontend 2>/dev/null || true
docker rm expense-mongodb expense-backend expense-frontend 2>/dev/null || true

# Create Docker network
docker network create expense-network 2>/dev/null || true

# Create environment file
print_status "Creating environment configuration..."
cat > .env << EOF
GEMINI_API_KEY=$GEMINI_API_KEY
SESSION_SECRET=$(openssl rand -hex 32)
PORT=8081
MONGODB_URI=mongodb://expense-mongodb:27017
EOF

# Build images
print_status "Building Docker images..."
echo "Building backend..."
docker build -t expense-backend ./backend

echo "Building frontend..."
docker build -t expense-frontend ./frontend

# Start MongoDB
print_status "Starting MongoDB..."
docker run -d --name expense-mongodb \
  --network expense-network \
  -p 27017:27017 \
  -v mongodb_data:/data/db \
  --restart unless-stopped \
  mongo:7

# Wait for MongoDB
print_status "Waiting for MongoDB to be ready..."
for i in {1..30}; do
  if docker exec expense-mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    print_status "MongoDB is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    print_error "MongoDB failed to start"
    docker logs expense-mongodb --tail 10
    exit 1
  fi
  echo "Waiting for MongoDB... ($i/30)"
  sleep 2
done

# Start backend
print_status "Starting backend..."
docker run -d --name expense-backend \
  --network expense-network \
  -p 8081:8081 \
  --env-file .env \
  --restart unless-stopped \
  expense-backend

# Wait for backend
print_status "Waiting for backend to be ready..."
for i in {1..60}; do
  if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    print_status "Backend is ready!"
    break
  fi
  if [ $i -eq 60 ]; then
    print_error "Backend failed to start"
    echo "Backend logs:"
    docker logs expense-backend --tail 20
    exit 1
  fi
  echo "Waiting for backend... ($i/60)"
  sleep 2
done

# Start frontend
print_status "Starting frontend..."
docker run -d --name expense-frontend \
  --network expense-network \
  -p 3000:80 \
  --restart unless-stopped \
  expense-frontend

# Wait for frontend
print_status "Waiting for frontend to be ready..."
for i in {1..30}; do
  if curl -s http://localhost:3000 > /dev/null 2>&1; then
    print_status "Frontend is ready!"
    break
  fi
  echo "Waiting for frontend... ($i/30)"
  sleep 2
done

# Health checks
echo ""
echo "🏥 Health Checks"
echo "================"
BACKEND_HEALTH=$(curl -s http://localhost:8081/health 2>/dev/null || echo "FAILED")
API_HEALTH=$(curl -s http://localhost:8081/api/health 2>/dev/null || echo "FAILED")
FRONTEND_STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 2>/dev/null || echo "FAILED")

echo "Backend Health: $BACKEND_HEALTH"
echo "API Health: $API_HEALTH"
echo "Frontend Status: $FRONTEND_STATUS"

# Show running containers
echo ""
echo "🐳 Running Containers"
echo "===================="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Show access URLs
echo ""
echo "🌐 Access URLs"
echo "=============="
echo "Frontend:    http://$PUBLIC_IP:3000"
echo "Backend API: http://$PUBLIC_IP:8081"
echo "Admin Panel: http://$PUBLIC_IP:8081/admin"
echo "Health:      http://$PUBLIC_IP:8081/health"

# Show useful commands
echo ""
echo "🛠️  Useful Commands"
echo "=================="
echo "View logs:       docker logs expense-backend"
echo "Restart all:     docker restart expense-mongodb expense-backend expense-frontend"
echo "Stop all:        docker stop expense-mongodb expense-backend expense-frontend"
echo "Update app:      cd $PROJECT_DIR/CAF/expense-tracker && git pull && docker-compose up -d --build"

# Final status
echo ""
if [ "$BACKEND_HEALTH" = "OK" ] && [ "$FRONTEND_STATUS" = "200" ]; then
    print_status "🎉 Deployment successful!"
    echo ""
    echo "Bạn có thể truy cập ứng dụng tại: http://$PUBLIC_IP:3000"
    echo ""
    echo "Để test:"
    echo "1. Mở http://$PUBLIC_IP:3000"
    echo "2. Đăng ký tài khoản mới"
    echo "3. Thêm chi phí: 'ăn trưa 50k'"
    echo "4. Xem báo cáo tại: http://$PUBLIC_IP:8081/admin"
else
    print_error "Deployment có vấn đề, kiểm tra logs:"
    echo "docker logs expense-backend"
    echo "docker logs expense-frontend"
fi

echo ""
print_warning "Lưu ý: Đảm bảo Security Group cho phép traffic trên port 3000 và 8081"