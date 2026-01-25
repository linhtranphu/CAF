#!/bin/bash
# Auto build and push to Docker Hub

set -e

# Config
DOCKER_USERNAME="linhtranphu"  # Thay bằng username Docker Hub của bạn
BACKEND_IMAGE="$DOCKER_USERNAME/expense-backend"
FRONTEND_IMAGE="$DOCKER_USERNAME/expense-frontend"
TAG="latest"

echo "🐳 Building and pushing Docker images"
echo "===================================="

# Check if logged in
if ! docker info | grep -q "Username"; then
    echo "❌ Chưa login Docker Hub. Chạy: docker login"
    exit 1
fi

# Get current directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📁 Project directory: $SCRIPT_DIR"

# Tag existing images
echo "🏷️  Tagging images..."
docker tag expense-tracker-backend:latest "$BACKEND_IMAGE:$TAG"
echo "✅ Backend tagged: $BACKEND_IMAGE:$TAG"

docker tag expense-tracker-frontend:latest "$FRONTEND_IMAGE:$TAG"
echo "✅ Frontend tagged: $FRONTEND_IMAGE:$TAG"

# Push images
echo "📤 Pushing to Docker Hub..."
docker push "$BACKEND_IMAGE:$TAG"
echo "✅ Backend pushed: $BACKEND_IMAGE:$TAG"

docker push "$FRONTEND_IMAGE:$TAG"
echo "✅ Frontend pushed: $FRONTEND_IMAGE:$TAG"

# Create updated docker-compose.yml
echo "📝 Creating docker-compose.yml with new images..."
cd "$SCRIPT_DIR"

cat > docker-compose.hub.yml << EOF
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
    image: $BACKEND_IMAGE:$TAG
    container_name: expense-backend
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=\${GEMINI_API_KEY}
      - SESSION_SECRET=\${SESSION_SECRET}
    depends_on:
      - mongodb
    restart: unless-stopped

  frontend:
    image: $FRONTEND_IMAGE:$TAG
    container_name: expense-frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

echo "✅ Created docker-compose.hub.yml"

echo ""
echo "🎉 Images pushed successfully!"
echo "=============================="
echo "Backend:  $BACKEND_IMAGE:$TAG"
echo "Frontend: $FRONTEND_IMAGE:$TAG"
echo ""
echo "📋 Next steps:"
echo "1. Update deploy scripts với images mới"
echo "2. Test deploy: docker-compose -f docker-compose.hub.yml up -d"
echo "3. Commit và push changes lên Git"