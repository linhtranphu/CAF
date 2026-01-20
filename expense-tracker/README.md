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
- **Deployment**: Docker + AWS ECS/EC2 ready

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- MongoDB (local or Atlas)
- Google Gemini API key

### Local Development
```bash
# Clone and setup
git clone <repository>
cd expense-tracker

# Configure environment
cp backend/.env.example backend/.env
# Edit .env with your Gemini API key

# Start all services (Windows)
./start-app.bat

# Or start manually:
# 1. MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:7

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

## 🌐 Deployment Options

### Option 1: AWS ECS Fargate (Recommended)
```bash
# See detailed guide
cat aws/README.md

# Quick deploy
cd aws
./deploy.sh
```

### Option 2: AWS EC2 (Simple)
```bash
# 1. Launch EC2 instance (t3.micro or larger)
# 2. Install Docker and Docker Compose
sudo yum update -y
sudo yum install -y docker
sudo systemctl start docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 3. Clone and deploy
git clone <your-repo>
cd expense-tracker

# 4. Configure production environment
cp backend/.env.example backend/.env
# Edit .env with production values:
# - GEMINI_API_KEY=your-production-key
# - MONGODB_URI=mongodb://mongodb:27017 (for Docker)
# - SESSION_SECRET=secure-random-string

# 5. Deploy with Docker Compose
docker-compose up -d

# 6. Configure security group
# Allow inbound: 22 (SSH), 80 (HTTP), 443 (HTTPS), 3000 (Frontend), 8081 (Backend)

# 7. Access application
# Frontend: http://your-ec2-ip:3000
# Backend Admin: http://your-ec2-ip:8081
```

### Option 3: Traditional Server
```bash
# Install dependencies
# Go, Node.js, MongoDB, PM2

# Build and run
cd backend && go build -o expense-tracker cmd/main.go
cd frontend && npm run build

# Use PM2 or systemd for process management
pm2 start ecosystem.config.js
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