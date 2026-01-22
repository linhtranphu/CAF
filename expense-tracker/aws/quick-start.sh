#!/bin/bash
# Quick Start Script for Expense Tracker on EC2
# Usage: curl -sSL https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/quick-start.sh | bash

set -e

echo "🚀 Expense Tracker - Quick Start"
echo "================================"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    print_error "Không chạy script này với sudo/root!"
    exit 1
fi

# Check internet connection
if ! ping -c 1 google.com &> /dev/null; then
    print_error "Không có kết nối internet!"
    exit 1
fi

print_status "Bắt đầu cài đặt Expense Tracker..."

# Step 1: Setup system
print_status "Bước 1: Cài đặt dependencies..."
if ! command -v docker &> /dev/null; then
    curl -sSL https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/setup.sh | bash
    
    print_warning "Docker đã được cài đặt. Vui lòng logout và login lại, sau đó chạy lệnh sau:"
    echo ""
    echo "curl -sSL https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/quick-start.sh | bash"
    echo ""
    exit 0
else
    print_status "Docker đã có sẵn"
fi

# Step 2: Check Docker permissions
if ! docker ps &> /dev/null; then
    print_error "Docker permission denied! Vui lòng logout và login lại sau khi chạy setup.sh"
    exit 1
fi

# Step 3: Get GEMINI API Key
echo ""
echo "🔑 GEMINI API Key"
echo "================="
echo "Bạn cần GEMINI API Key để sử dụng tính năng AI parsing chi phí"
echo "Lấy miễn phí tại: https://makersuite.google.com/app/apikey"
echo ""

if [ -z "$GEMINI_API_KEY" ]; then
    read -p "Nhập GEMINI_API_KEY của bạn: " GEMINI_API_KEY
fi

if [ -z "$GEMINI_API_KEY" ]; then
    print_error "GEMINI_API_KEY là bắt buộc!"
    exit 1
fi

# Step 4: Deploy application
print_status "Bước 2: Deploy ứng dụng..."
export GEMINI_API_KEY="$GEMINI_API_KEY"

# Download and run deploy script
PROJECT_DIR="$HOME/expense-tracker"
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

curl -sSL https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/deploy-ec2.sh -o deploy-ec2.sh
chmod +x deploy-ec2.sh

# Run deploy
./deploy-ec2.sh

print_status "🎉 Quick Start hoàn thành!"