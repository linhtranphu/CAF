# 🐳 Container Deploy - Expense Tracker

## 🚀 Quick Deploy (1 lệnh)

```bash
# SSH vào EC2 Ubuntu và chạy:
curl -sSL https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/container-deploy.sh | bash
```

## 📋 Manual Deploy

### 1. Setup Docker
```bash
# Install Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER

# Logout và login lại
exit
ssh -i your-key.pem ubuntu@your-ec2-ip
```

### 2. Deploy với Docker Compose
```bash
# Tạo project directory
mkdir ~/expense-tracker && cd ~/expense-tracker

# Tạo docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'
services:
  mongodb:
    image: mongo:7
    ports: ["27017:27017"]
    volumes: [mongodb_data:/data/db]
    restart: unless-stopped

  backend:
    image: linhtranphu/expense-backend:latest
    ports: ["8081:8081"]
    environment:
      - MONGODB_URI=mongodb://mongodb:27017
      - GEMINI_API_KEY=YOUR_GEMINI_API_KEY_HERE
      - SESSION_SECRET=your-secret-key
    depends_on: [mongodb]
    restart: unless-stopped

  frontend:
    image: linhtranphu/expense-frontend:latest
    ports: ["3000:80"]
    depends_on: [backend]
    restart: unless-stopped

volumes:
  mongodb_data:
EOF

# Start services
docker-compose up -d
```

### 3. Kiểm tra
```bash
# Check containers
docker ps

# Test endpoints
curl http://localhost:8081/health
curl http://localhost:3000

# View logs
docker logs expense-backend
```

## 🌐 Access URLs

- **Frontend**: `http://your-ec2-ip:3000`
- **Backend**: `http://your-ec2-ip:8081`
- **Admin**: `http://your-ec2-ip:8081/admin`

## 🔧 Management Commands

```bash
# Update containers
docker-compose pull && docker-compose up -d

# Restart services
docker-compose restart

# Stop all
docker-compose down

# View logs
docker-compose logs -f

# Backup MongoDB
docker exec expense-mongodb mongodump --out /tmp/backup
```

## 🛡️ Security Group Rules

```
Type        Port    Source
SSH         22      Your IP
HTTP        80      0.0.0.0/0
Custom      3000    0.0.0.0/0
Custom      8081    0.0.0.0/0
```

## 🔑 Environment Variables

Cần GEMINI API Key từ: https://makersuite.google.com/app/apikey

## 📊 Monitoring

```bash
# Health check script
cat > health.sh << 'EOF'
#!/bin/bash
echo "Backend: $(curl -s http://localhost:8081/health || echo 'FAILED')"
echo "Frontend: $(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 || echo 'FAILED')"
docker ps --format "table {{.Names}}\t{{.Status}}"
EOF

chmod +x health.sh
./health.sh
```