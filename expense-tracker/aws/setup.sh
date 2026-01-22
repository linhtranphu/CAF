#!/bin/bash
set -e

echo "🚀 Setting up Expense Tracker on AWS EC2..."

# Update system
sudo yum update -y

# Install Docker
sudo yum install -y docker git
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose v2
sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

echo "✅ Setup complete!"
echo "⚠️  Please logout and login again for Docker permissions"
echo "Then run: ./deploy.sh"