# 🚀 Hướng dẫn Deploy Expense Tracker lên AWS EC2

## Bước 1: Tạo EC2 Instance

### 1.1 Tạo EC2 Instance
- **AMI**: Amazon Linux 2023 (hoặc Amazon Linux 2)
- **Instance Type**: t3.micro (Free tier) hoặc t3.small
- **Storage**: 20GB gp3
- **Security Group**: Tạo mới với các rules sau:

```
Type        Protocol    Port Range    Source
SSH         TCP         22           Your IP
HTTP        TCP         80           0.0.0.0/0
Custom TCP  TCP         3000         0.0.0.0/0  (Frontend)
Custom TCP  TCP         8081         0.0.0.0/0  (Backend API)
```

### 1.2 Tạo Key Pair
- Tạo key pair mới hoặc sử dụng existing
- Download file `.pem` và lưu an toàn

## Bước 2: Kết nối và Setup EC2

### 2.1 SSH vào EC2
```bash
# Thay your-key.pem và your-ec2-ip
chmod 400 your-key.pem
ssh -i your-key.pem ec2-user@your-ec2-ip
```

### 2.2 Setup Dependencies
```bash
# Chạy setup script
curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/setup.sh
chmod +x setup.sh
./setup.sh

# Logout và login lại để apply Docker permissions
exit
ssh -i your-key.pem ec2-user@your-ec2-ip
```

## Bước 3: Deploy Application

### Phương án A: Deploy Tự động (Khuyến nghị)
```bash
# Download và chạy deploy script
curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/aws/deploy-simple.sh
chmod +x deploy-simple.sh
./deploy-simple.sh

# Nhập GEMINI_API_KEY khi được hỏi
```

### Phương án B: Deploy Thủ công
```bash
# 1. Clone repository
git clone https://github.com/linhtranphu/CAF.git
cd CAF/expense-tracker

# 2. Tạo .env file
cat > .env << EOF
GEMINI_API_KEY=your_gemini_api_key_here
SESSION_SECRET=$(openssl rand -hex 32)
PORT=8081
MONGODB_URI=mongodb://mongodb:27017
EOF

# 3. Deploy với Docker Compose
docker-compose up -d --build
```

## Bước 4: Kiểm tra Deployment

### 4.1 Kiểm tra containers
```bash
docker ps
```

### 4.2 Test endpoints
```bash
# Backend health
curl http://localhost:8081/health
curl http://localhost:8081/api/health

# Frontend
curl http://localhost:3000
```

### 4.3 Xem logs
```bash
docker logs expense-backend
docker logs expense-frontend
docker logs expense-mongodb
```

## Bước 5: Truy cập ứng dụng

- **Frontend**: `http://your-ec2-public-ip:3000`
- **Backend API**: `http://your-ec2-public-ip:8081`
- **Admin Panel**: `http://your-ec2-public-ip:8081/admin`

## Troubleshooting

### Lỗi thường gặp:

1. **"Failed to fetch" ở frontend**
   - Kiểm tra Security Group có mở port 8081
   - Kiểm tra backend có chạy: `docker logs expense-backend`

2. **Backend không start**
   ```bash
   # Kiểm tra logs
   docker logs expense-backend
   
   # Restart container
   docker restart expense-backend
   ```

3. **MongoDB connection failed**
   ```bash
   # Kiểm tra MongoDB
   docker logs expense-mongodb
   docker exec -it expense-mongodb mongosh
   ```

4. **Port đã được sử dụng**
   ```bash
   # Stop tất cả containers
   docker stop $(docker ps -q)
   docker rm $(docker ps -aq)
   
   # Chạy lại deploy
   ./deploy-simple.sh
   ```

### Commands hữu ích:

```bash
# Restart tất cả services
docker restart expense-mongodb expense-backend expense-frontend

# Xem resource usage
docker stats

# Cleanup
docker system prune -f

# Update application
cd CAF/expense-tracker
git pull origin main
docker-compose up -d --build
```

## Bảo mật Production

### 1. Sử dụng HTTPS
```bash
# Install Certbot
sudo yum install -y certbot

# Get SSL certificate (cần domain name)
sudo certbot certonly --standalone -d your-domain.com
```

### 2. Firewall rules
```bash
# Chỉ cho phép SSH từ IP cụ thể
# Sử dụng Load Balancer cho HTTPS
# Đặt backend ở private subnet
```

### 3. Environment variables
```bash
# Không commit .env file
# Sử dụng AWS Secrets Manager
# Rotate SESSION_SECRET định kỳ
```

## Monitoring

### 1. Health checks
```bash
# Tạo health check script
cat > health-check.sh << 'EOF'
#!/bin/bash
echo "=== Health Check $(date) ==="
echo "Backend: $(curl -s http://localhost:8081/health || echo 'FAILED')"
echo "Frontend: $(curl -s -o /dev/null -w '%{http_code}' http://localhost:3000 || echo 'FAILED')"
echo "Containers: $(docker ps --format 'table {{.Names}}\t{{.Status}}')"
EOF

chmod +x health-check.sh
```

### 2. Logs rotation
```bash
# Setup logrotate cho Docker logs
sudo tee /etc/logrotate.d/docker << EOF
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size=1M
    missingok
    delaycompress
    copytruncate
}
EOF
```

## Backup

### 1. MongoDB backup
```bash
# Backup script
cat > backup.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec expense-mongodb mongodump --out /tmp/backup_$DATE
docker cp expense-mongodb:/tmp/backup_$DATE ./backup_$DATE
tar -czf backup_$DATE.tar.gz backup_$DATE
rm -rf backup_$DATE
echo "Backup created: backup_$DATE.tar.gz"
EOF

chmod +x backup.sh
```

### 2. Automated backup với cron
```bash
# Add to crontab
crontab -e

# Backup hàng ngày lúc 2AM
0 2 * * * /home/ec2-user/backup.sh
```