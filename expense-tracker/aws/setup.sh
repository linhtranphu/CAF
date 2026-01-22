#!/bin/bash
set -e

echo "🚀 Setting up Expense Tracker on AWS EC2..."
echo "==========================================="

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

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
else
    print_error "Cannot detect OS"
    exit 1
fi

print_status "Detected OS: $OS"

# Update system
print_status "Updating system packages..."
if [[ "$OS" == *"Amazon Linux"* ]]; then
    sudo yum update -y
    PACKAGE_MANAGER="yum"
elif [[ "$OS" == *"Ubuntu"* ]]; then
    sudo apt update && sudo apt upgrade -y
    PACKAGE_MANAGER="apt"
else
    print_warning "Unsupported OS: $OS. Trying with yum..."
    sudo yum update -y
    PACKAGE_MANAGER="yum"
fi

# Install Docker
if command -v docker &> /dev/null; then
    print_status "Docker already installed"
else
    print_status "Installing Docker..."
    if [ "$PACKAGE_MANAGER" = "yum" ]; then
        sudo yum install -y docker git curl
    else
        sudo apt install -y docker.io git curl
    fi
    
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -a -G docker $USER
fi

# Install Docker Compose
if command -v docker-compose &> /dev/null; then
    print_status "Docker Compose already installed"
else
    print_status "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create docker-compose symlink if needed
if [ ! -f /usr/bin/docker-compose ] && [ -f /usr/local/bin/docker-compose ]; then
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
fi

# Test Docker installation
print_status "Testing Docker installation..."
if docker --version && docker-compose --version; then
    print_status "Docker and Docker Compose installed successfully!"
else
    print_error "Docker installation failed"
    exit 1
fi

# Check if user needs to logout
if ! groups $USER | grep -q docker; then
    print_warning "User $USER is not in docker group yet"
    NEED_LOGOUT=true
else
    print_status "User $USER is already in docker group"
    NEED_LOGOUT=false
fi

# Create project directory
PROJECT_DIR="$HOME/expense-tracker"
mkdir -p "$PROJECT_DIR"
print_status "Created project directory: $PROJECT_DIR"

# Download deploy script
print_status "Downloading deploy script..."
cd "$PROJECT_DIR"
curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/deploy-ec2.sh
chmod +x deploy-ec2.sh

echo ""
print_status "Setup complete!"
echo ""
echo "📋 Next Steps:"
if [ "$NEED_LOGOUT" = true ]; then
    print_warning "1. Logout and login again to apply Docker permissions:"
    echo "   exit"
    echo "   ssh -i your-key.pem ec2-user@your-ec2-ip"
    echo ""
    print_status "2. Then run the deploy script:"
    echo "   cd $PROJECT_DIR"
    echo "   ./deploy-ec2.sh"
else
    print_status "Run the deploy script:"
    echo "   cd $PROJECT_DIR"
    echo "   ./deploy-ec2.sh"
fi

echo ""
echo "🔑 You will need a GEMINI API Key from: https://makersuite.google.com/app/apikey"
echo ""
print_warning "Make sure your Security Group allows inbound traffic on ports 22, 3000, and 8081"