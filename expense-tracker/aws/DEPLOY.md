# 🚀 AWS EC2 Deployment Guide

## Step 1: Setup EC2 Instance

```bash
# SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Run setup
curl -O https://raw.githubusercontent.com/your-repo/expense-tracker/main/aws/setup.sh
chmod +x setup.sh
./setup.sh

# Logout and login again
exit
ssh -i your-key.pem ec2-user@your-ec2-ip
```

## Step 2: Deploy Application

### Option A: Simple Deploy (Recommended)
```bash
curl -O https://raw.githubusercontent.com/your-repo/expense-tracker/main/aws/deploy-simple.sh
chmod +x deploy-simple.sh
./deploy-simple.sh
```

### Option B: Manual Deploy (If Docker Compose fails)
```bash
# Clone repository
git clone https://github.com/your-username/expense-tracker.git
cd expense-tracker

# Run manual deploy
chmod +x aws/deploy-manual.sh
./aws/deploy-manual.sh
```

### Option C: Step by Step Manual
```bash
# 1. Start MongoDB
docker run -d --name expense-mongodb -p 27017:27017 mongo:7

# 2. Fix Go version
cd expense-tracker/backend
sed -i 's/go 1.24/go 1.21/' go.mod

# 3. Build Backend
docker build -t expense-backend .

# 4. Start Backend
docker run -d --name expense-backend -p 8081:8081 \
  -e PORT=8081 \
  -e MONGODB_URI=mongodb://172.17.0.1:27017 \
  -e GEMINI_API_KEY=AIzaSyD_X1AdGqKXQ0EETgC80BDWYt8zKuSyviM \
  -e SESSION_SECRET=expense-tracker-secret-123 \
  expense-backend

# 5. Build Frontend
cd ../frontend
docker build -t expense-frontend .

# 6. Start Frontend
docker run -d --name expense-frontend -p 3000:80 expense-frontend
```

## Step 3: Verify Deployment

```bash
# Check containers
docker ps

# Test endpoints
curl http://localhost:8081/api/health
curl http://localhost:3000

# View logs
docker logs expense-backend
docker logs expense-frontend
```

## Access URLs
- Frontend: http://your-ec2-ip:3000
- Backend: http://your-ec2-ip:8081

## Troubleshooting

```bash
# Restart containers
docker restart expense-mongodb expense-backend expense-frontend

# View logs
docker logs expense-backend

# Check MongoDB connection
docker exec -it expense-mongodb mongo
```