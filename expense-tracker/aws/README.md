# Expense Tracker - AWS Deployment Guide

## 🏗️ Architecture
- **Frontend**: Vue.js SPA hosted on S3 + CloudFront
- **Backend**: Go API on ECS Fargate
- **Database**: MongoDB Atlas or EC2
- **Load Balancer**: Application Load Balancer
- **Secrets**: AWS Secrets Manager

## 📋 Prerequisites
1. AWS CLI configured
2. Docker installed
3. MongoDB Atlas cluster or EC2 MongoDB

## 🚀 Deployment Steps

### 1. Setup AWS Resources
```bash
# Create ECR repositories
aws ecr create-repository --repository-name expense-tracker-backend
aws ecr create-repository --repository-name expense-tracker-frontend

# Create ECS cluster
aws ecs create-cluster --cluster-name expense-tracker-cluster

# Create secrets in Secrets Manager
aws secretsmanager create-secret --name expense-tracker/gemini-api-key --secret-string "your-gemini-key"
aws secretsmanager create-secret --name expense-tracker/session-secret --secret-string "your-session-secret"
```

### 2. Update Configuration
1. Update `ecs-task-definition.json` with your AWS account ID
2. Update MongoDB URI in task definition
3. Update CORS origins in backend code for production domain

### 3. Deploy
```bash
chmod +x deploy.sh
./deploy.sh
```

### 4. Setup Load Balancer
1. Create Application Load Balancer
2. Create target group for ECS service
3. Configure health checks on `/api/health`

### 5. Setup CloudFront (Frontend)
1. Upload frontend build to S3
2. Create CloudFront distribution
3. Configure custom domain (optional)

## 🔧 Environment Variables
- `PORT`: 8081
- `MONGODB_URI`: MongoDB connection string
- `GEMINI_API_KEY`: Google Gemini API key
- `SESSION_SECRET`: Secure session secret

## 🔐 Security
- Use AWS Secrets Manager for sensitive data
- Enable HTTPS with ACM certificates
- Configure security groups properly
- Use IAM roles with minimal permissions

## 📊 Monitoring
- CloudWatch logs for ECS tasks
- CloudWatch metrics for performance
- ALB access logs
- MongoDB Atlas monitoring

## 💰 Cost Optimization
- Use Fargate Spot for non-production
- Configure auto-scaling
- Use CloudFront caching
- Monitor with AWS Cost Explorer