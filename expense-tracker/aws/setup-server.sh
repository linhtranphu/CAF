#!/bin/bash
# EC2 Server Setup Script

echo "🚀 Setting up Expense Tracker on AWS EC2..."

# Update system
sudo yum update -y

# Install Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
sudo yum install -y git

# Create app directory
mkdir -p /home/ec2-user/expense-tracker
cd /home/ec2-user/expense-tracker

echo "✅ Server setup complete!"
echo "Next steps:"
echo "1. Clone your repository"
echo "2. Configure environment variables"
echo "3. Run docker-compose up -d"