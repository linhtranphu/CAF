#!/bin/bash
set -e

echo "🚀 Deploying Expense Tracker to EC2 (Amazon Linux)..."

# 1. Check/Install Dependencies (Amazon Linux)
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker & Git..."
    sudo yum update -y
    sudo yum install -y docker git
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -a -G docker ec2-user
    
    # Install Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    echo "⚠️  Docker installed. Please logout and login again to apply group permissions."
    exit 1
fi

# 2. Configuration
read -p "Enter GEMINI_API_KEY: " GEMINI_API_KEY
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required!"
    exit 1
fi

# 3. Setup Project Directory
mkdir -p ~/expense-tracker
cd ~/expense-tracker

# 4. Clone Source Code (if needed)
if [ ! -d "backend" ]; then
    echo "⬇️  Cloning repository..."
    # Clone source code to build images locally
    git clone https://github.com/linhtranphu/CAF.git temp_repo
    
    if [ -d "temp_repo/expense-tracker" ]; then
        cp -r temp_repo/expense-tracker/* .
    else
        cp -r temp_repo/* .
    fi
    rm -rf temp_repo
fi

# 5. Create environment file
cat > .env << EOF
GEMINI_API_KEY=$GEMINI_API_KEY
SESSION_SECRET=$(openssl rand -hex 32)
PORT=8081
MONGODB_URI=mongodb://mongodb:27017
EOF

# 6. Create docker-compose.yml
cat > docker-compose.yml << 'EOF'
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
      timeout: 10s
      retries: 5

  backend:
    build: ./backend
    container_name: expense-backend
    ports:
      - "8081:8081"
    env_file: .env
    depends_on:
      mongodb:
        condition: service_healthy
    restart: unless-stopped

  frontend:
    build: ./frontend
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

echo "🏗️  Building and Starting services..."
docker-compose up -d --build

# 7. Display Access Info
TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)
PUBLIC_IP=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s http://169.254.169.254/latest/meta-data/public-ipv4 || curl -s http://169.254.169.254/latest/meta-data/public-ipv4)

echo "✅ Deployment Complete!"
echo "Frontend: http://${PUBLIC_IP}:3000"
echo "Backend:  http://${PUBLIC_IP}:8081"