#!/bin/bash
set -e

echo "🚀 Deploying Expense Tracker..."

# Get GEMINI API Key
read -p "Enter GEMINI_API_KEY: " GEMINI_API_KEY
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY is required!"
    exit 1
fi
# Clone repository if not exists
if [ ! -d "CAF" ]; then
    git clone https://github.com/linhtranphu/CAF.git
fi

cd CAF/expense-tracker

# Create environment file
cat > .env << EOF
GEMINI_API_KEY=$GEMINI_API_KEY
SESSION_SECRET=expense-tracker-secret-$(date +%s)
EOF

# Fix go.mod version
sed -i 's/go 1.24/go 1.21/' backend/go.mod

# Create simple docker-compose
cat > docker-compose.simple.yml << 'EOF'
services:
  mongodb:
    image: public.ecr.aws/docker/library/mongo:7
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

  backend:
    build: ./backend
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - SESSION_SECRET=${SESSION_SECRET}
    depends_on:
      - mongodb

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

volumes:
  mongodb_data:
EOF

# Deploy
docker-compose -f docker-compose.simple.yml up -d --build

echo "✅ Deployment complete!"
echo "Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo "Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8081"