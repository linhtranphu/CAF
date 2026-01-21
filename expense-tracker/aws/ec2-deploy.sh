#!/bin/bash
set -e

echo "🚀 Deploying Expense Tracker to EC2..."

# Get GEMINI API Key
read -p "Enter GEMINI_API_KEY: " GEMINI_API_KEY
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required!"
    exit 1
fi

# Create project directory
mkdir -p ~/expense-tracker
cd ~/expense-tracker

# Create environment file
cat > .env << EOF
GEMINI_API_KEY=$GEMINI_API_KEY
SESSION_SECRET=expense-tracker-secret-$(date +%s)
EOF

# Create docker-compose file
cat > docker-compose.yml << 'EOF'
services:
  mongodb:
    image: public.ecr.aws/docker/library/mongo:7
    container_name: expense-mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    restart: unless-stopped

  backend:
    image: expense-backend:latest
    container_name: expense-backend
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - SESSION_SECRET=${SESSION_SECRET}
    depends_on:
      - mongodb
    restart: unless-stopped

  frontend:
    image: expense-frontend:latest
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

echo "✅ Configuration created!"
echo "📝 Next steps:"
echo "1. Upload your Docker images or source code"
echo "2. Run: docker-compose up -d"
echo "3. Access: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"