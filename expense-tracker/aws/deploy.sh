#!/bin/bash

# AWS Deployment Script for Expense Tracker

set -e

# Configuration
AWS_REGION="us-east-1"
ECR_REPOSITORY="expense-tracker"
ECS_CLUSTER="expense-tracker-cluster"
ECS_SERVICE="expense-tracker-service"

# Get AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY}"

echo "🚀 Starting deployment to AWS..."

# 1. Build and push backend image
echo "📦 Building backend image..."
cd backend
docker build -t ${ECR_REPOSITORY}-backend .

# Login to ECR
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_URI}

# Tag and push
docker tag ${ECR_REPOSITORY}-backend:latest ${ECR_URI}-backend:latest
docker push ${ECR_URI}-backend:latest

# 2. Build and push frontend image
echo "📦 Building frontend image..."
cd ../frontend
docker build -t ${ECR_REPOSITORY}-frontend .
docker tag ${ECR_REPOSITORY}-frontend:latest ${ECR_URI}-frontend:latest
docker push ${ECR_URI}-frontend:latest

# 3. Update ECS service
echo "🔄 Updating ECS service..."
cd ../aws
aws ecs update-service \
  --cluster ${ECS_CLUSTER} \
  --service ${ECS_SERVICE} \
  --force-new-deployment \
  --region ${AWS_REGION}

echo "✅ Deployment completed!"
echo "🌐 Backend: https://your-alb-url.amazonaws.com"
echo "🌐 Frontend: https://your-cloudfront-url.cloudfront.net"