# 💰 Expense Tracker

Modern expense tracking application with AI-powered message parsing and simple authentication.

## ✨ Features
- **Simple Login**: Username/password authentication (no OAuth complexity)
- **AI Message Parsing**: Natural language expense input using Google Gemini
- **Real-time Tracking**: Add expenses via intuitive web interface
- **Admin Dashboard**: View expenses and summaries grouped by person
- **Soft Delete**: Expenses marked as deleted for audit trail
- **MongoDB Storage**: Scalable NoSQL database with ObjectID
- **Responsive Design**: Works seamlessly on desktop and mobile
- **Session Management**: Secure cookie-based authentication

## 🏗️ Architecture
- **Frontend**: Vue.js 3 + Vite (Port 3000)
- **Backend**: Go + Gin framework (Port 8081)
- **Database**: MongoDB (Port 27017)
- **AI**: Google Gemini 2.5 Flash Lite API
- **Authentication**: Session-based with secure cookies
- **Deployment**: Docker + AWS EC2 ready

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Node.js 18+
- MongoDB (local or Atlas)
- Google Gemini API key

### Local Development
```bash
# Clone and setup
git clone https://github.com/linhtranphu/CAF.git
cd CAF/expense-tracker

# Configure environment
cp backend/.env.example backend/.env
# Edit .env with your Gemini API key

# Start all services (Windows)
./start-app.bat

# Or start manually:
# 1. MongoDB
docker run -d -p 27017:27017 --name mongodb public.ecr.aws/docker/library/mongo:7

# 2. Backend
cd backend
go mod tidy
go run cmd/main.go

# 3. Frontend
cd frontend
npm install
npm run dev
```

### Docker Deployment
```bash
# Local development
docker-compose up -d

# Production build
docker-compose -f docker-compose.prod.yml up -d
```

## 🔧 Configuration

### Environment Variables (.env)
```env
PORT=8081
GEMINI_API_KEY=your-gemini-api-key-here
MONGODB_URI=mongodb://localhost:27017
SESSION_SECRET=your-secure-session-secret-change-in-production
```

### Demo Users
| Username | Password | Role |
|----------|----------|------|
| admin    | admin123 | Admin |
| linh     | linh123  | User  |
| toan     | toan123  | User  |

## 📱 Usage

### For Users (Frontend - Port 3000)
1. Visit `http://localhost:3000`
2. Login with demo credentials
3. Add expenses using natural language:
   - "ăn trưa 50k"
   - "mua xăng 200 nghìn"
   - "cọc nhà 34 triệu"
4. Logout when done

### For Admins (Backend - Port 8081)
1. Visit `http://localhost:8081`
2. Login with admin credentials
3. View expense dashboard with:
   - Total transactions count
   - Summary by person with grand total
   - Detailed expense table
   - Soft delete functionality

## 🌐 AWS EC2 Deployment

### Step 1: Create EC2 Instance
1. **AWS Console** → **EC2** → **Launch Instance**
2. Configuration:
   - **AMI**: Amazon Linux 2023 AMI 2023.10.20260105.0 x86_64 HVM kernel-6.1
   - **Instance Type**: t3.micro (free tier)
   - **Key Pair**: Create or select key pair
   - **Security Group**: 
     - SSH (22) - 0.0.0.0/0
     - HTTP (80) - 0.0.0.0/0
     - Custom TCP (3000) - 0.0.0.0/0
     - Custom TCP (8081) - 0.0.0.0/0
3. **Launch Instance**

### Step 2: Connect to EC2
```bash
# SSH to EC2 (Windows - use PowerShell or Git Bash)
ssh -i "your-key.pem" ec2-user@your-ec2-public-ip
```

### Step 3: Install Docker
```bash
# Update system
sudo yum update -y

# Install Docker and Git
sudo yum install -y docker git

# Start Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -a -G docker ec2-user

# Logout and login again
exit
ssh -i "your-key.pem" ec2-user@your-ec2-public-ip
```

### Step 4: Clone and Deploy
```bash
# Clone repository
git clone https://github.com/linhtranphu/CAF.git
cd CAF/expense-tracker

# Run deploy script
chmod +x aws/deploy-simple.sh
./aws/deploy-simple.sh

# Enter GEMINI_API_KEY when prompted
```

### Step 5: Verify Deployment
```bash
# Check containers
docker ps

# Check logs if needed
docker logs expense-backend
docker logs expense-frontend
docker logs expense-mongodb

# Test endpoints
curl http://localhost:8081
curl http://localhost:3000
```

### Step 6: Access Application
- **Frontend**: `http://your-ec2-public-ip:3000`
- **Backend API**: `http://your-ec2-public-ip:8081`
- **Admin Panel**: `http://your-ec2-public-ip:8081/admin`

### Useful Commands
```bash
# Restart all containers
docker restart expense-mongodb expense-backend expense-frontend

# View logs real-time
docker logs -f expense-backend

# Stop and remove containers
docker stop expense-mongodb expense-backend expense-frontend
docker rm expense-mongodb expense-backend expense-frontend
```

## 🛠️ Tech Stack

### Backend
- **Go 1.21** - High performance, compiled language
- **Gin** - Fast HTTP web framework
- **MongoDB Driver** - Official Go driver
- **Gin Sessions** - Cookie-based session management
- **CORS** - Cross-origin resource sharing
- **Google Gemini API** - AI message parsing

### Frontend
- **Vue.js 3** - Progressive JavaScript framework
- **Vite** - Fast build tool and dev server
- **Fetch API** - HTTP client for backend communication
- **CSS3** - Modern styling with gradients and animations

### Database
- **MongoDB** - Document-based NoSQL database
- **ObjectID** - Native MongoDB identifiers
- **Soft Delete** - Status-based deletion for audit

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **AWS ECS** - Container orchestration service
- **AWS EC2** - Virtual private servers
- **Nginx** - Reverse proxy and static file serving

## 🔐 Security Features
- Session-based authentication with secure cookies
- CORS protection with specific origins
- Environment variable configuration
- Soft delete for data integrity
- Input validation and sanitization
- MongoDB injection protection

## 📊 Monitoring & Logs
- Structured logging with request/response details
- MongoDB operation logging
- Session management logging
- Error handling with proper HTTP status codes

## 🚨 Production Checklist
- [ ] Change SESSION_SECRET to secure random string
- [ ] Use production MongoDB (Atlas recommended)
- [ ] Configure proper CORS origins
- [ ] Set up HTTPS with SSL certificates
- [ ] Configure proper security groups/firewall
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy for MongoDB
- [ ] Use environment-specific .env files

## 🤝 Contributing
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

## 📄 License
MIT License - see LICENSE file for details

## 🆘 Troubleshooting

### Common Issues
1. **MongoDB Connection**: Ensure MongoDB is running and URI is correct
2. **Gemini API**: Verify API key is valid and has quota
3. **CORS Errors**: Check allowed origins in router.go
4. **Session Issues**: Verify SESSION_SECRET is set
5. **Port Conflicts**: Ensure ports 3000, 8081, 27017 are available

### Debug Commands
```bash
# Check MongoDB connection
mongo mongodb://localhost:27017

# View backend logs
docker-compose logs backend

# Check running containers
docker ps

# Test API endpoints
curl http://localhost:8081/api/health
```

## 🛠️ Tech Stack

### Backend
- **Go 1.23** - High performance, compiled language
- **Gin** - Fast HTTP web framework
- **MongoDB Driver** - Official Go driver
- **Gin Sessions** - Cookie-based session management
- **CORS** - Cross-origin resource sharing
- **Google Gemini API** - AI message parsing

### Frontend
- **Vue.js 3** - Progressive JavaScript framework
- **Vite** - Fast build tool and dev server
- **Fetch API** - HTTP client for backend communication
- **CSS3** - Modern styling with gradients and animations

### Database
- **MongoDB** - Document-based NoSQL database
- **ObjectID** - Native MongoDB identifiers
- **Soft Delete** - Status-based deletion for audit

### DevOps
- **Docker** - Containerization with Amazon Linux 2023
- **AWS ECR** - Container registry for optimized images
- **AWS EC2** - Virtual private servers
- **Nginx** - Reverse proxy and static file serving

## 🔐 Security Features
- Session-based authentication with secure cookies
- CORS protection with specific origins
- Environment variable configuration
- Soft delete for data integrity
- Input validation and sanitization
- MongoDB injection protection

## 📈 Monitoring & Logs
- Structured logging with request/response details
- MongoDB operation logging
- Session management logging
- Error handling with proper HTTP status codes

## 🚨 Production Checklist
- [ ] Change SESSION_SECRET to secure random string
- [ ] Use production MongoDB (Atlas recommended)
- [ ] Configure proper CORS origins
- [ ] Set up HTTPS with SSL certificates
- [ ] Configure proper security groups/firewall
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy for MongoDB
- [ ] Use environment-specific .env files
- [ ] Set up Elastic IP for EC2
- [ ] Configure domain name with Route 53

## 🤝 Contributing
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

## 📄 License
MIT License - see LICENSE file for details

## 🆘 Troubleshooting

### Common Issues
1. **MongoDB Connection**: Ensure MongoDB is running and URI is correct
2. **Gemini API**: Verify API key is valid and has quota
3. **CORS Errors**: Check allowed origins in router.go
4. **Session Issues**: Verify SESSION_SECRET is set
5. **Port Conflicts**: Ensure ports 3000, 8081, 27017 are available
6. **Docker Build**: Check Go version compatibility (requires 1.23+)

### Debug Commands
```bash
# Check MongoDB connection
docker exec -it expense-mongodb mongosh

# View backend logs
docker logs expense-backend

# Check running containers
docker ps

# Test API endpoints
curl http://localhost:8081/api/health

# Restart containers
docker restart expense-mongodb expense-backend expense-frontend

# Check EC2 metadata
curl http://169.254.169.254/latest/meta-data/public-ipv4
```