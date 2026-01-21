#!/bin/bash
# Automated Deployment Script for AWS EC2

set -e

echo "🚀 Starting Expense Tracker Deployment..."

# Variables
APP_DIR="/home/ec2-user/expense-tracker"
REPO_URL="https://github.com/your-username/expense-tracker.git"  # Update this
BRANCH="main"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as ec2-user
if [ "$USER" != "ec2-user" ]; then
    log_error "Please run this script as ec2-user"
    exit 1
fi

# Create app directory
log_info "Creating application directory..."
mkdir -p $APP_DIR
cd $APP_DIR

# Clone or update repository
if [ -d ".git" ]; then
    log_info "Updating existing repository..."
    git pull origin $BRANCH
else
    log_info "Cloning repository..."
    git clone -b $BRANCH $REPO_URL .
fi

# Check if .env exists
if [ ! -f "backend/.env" ]; then
    log_warn "backend/.env not found. Creating from template..."
    cp aws/.env.production backend/.env
    log_error "Please edit backend/.env with your production values:"
    log_error "- GEMINI_API_KEY"
    log_error "- SESSION_SECRET"
    log_error "Then run this script again."
    exit 1
fi

# Stop existing containers
log_info "Stopping existing containers..."
docker-compose -f aws/docker-compose.prod.yml down || true

# Build and start containers
log_info "Building and starting containers..."
docker-compose -f aws/docker-compose.prod.yml up -d --build

# Wait for services to be ready
log_info "Waiting for services to start..."
sleep 30

# Health check
log_info "Performing health checks..."
if curl -f http://localhost:8081/api/health > /dev/null 2>&1; then
    log_info "✅ Backend is healthy"
else
    log_error "❌ Backend health check failed"
fi

if curl -f http://localhost:3000 > /dev/null 2>&1; then
    log_info "✅ Frontend is healthy"
else
    log_error "❌ Frontend health check failed"
fi

# Show status
log_info "Deployment complete! 🎉"
echo ""
echo "📊 Service Status:"
docker-compose -f aws/docker-compose.prod.yml ps
echo ""
echo "🌐 Access URLs:"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081"
echo ""
echo "📝 Useful Commands:"
echo "View logs:    docker-compose -f aws/docker-compose.prod.yml logs -f"
echo "Restart:      docker-compose -f aws/docker-compose.prod.yml restart"
echo "Stop:         docker-compose -f aws/docker-compose.prod.yml down"