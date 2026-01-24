#!/bin/bash
# Quick build and push script

set -e

echo "🏗️  Building and pushing Docker images..."

# Go to project root
cd /workspaces/CAF/expense-tracker

# Build backend
echo "Building backend..."
cd backend
docker build -t linhtranphu/expense-backend:latest .

# Build frontend
echo "Building frontend..."
cd ../frontend
docker build -t linhtranphu/expense-frontend:latest .

# Push both
echo "Pushing images..."
docker push linhtranphu/expense-backend:latest
docker push linhtranphu/expense-frontend:latest

echo "✅ Done! Images pushed to Docker Hub"