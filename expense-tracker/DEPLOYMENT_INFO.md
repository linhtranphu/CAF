# 🚀 Deployment Information

## 📦 Docker Images Ready for EC2

### Backend Image
- **Image**: `linhtranphu/expense-backend:latest`
- **Size**: 175MB
- **Digest**: `sha256:d6f2f48691a9a56ec671d7a2c20092cd4b93323807f33ad51d501b16a42df419`
- **Features**:
  - ✅ AI-powered expense parsing (Gemini API)
  - ✅ Original message tracking
  - ✅ User management with MongoDB
  - ✅ Role-based system (admin/supervisor)
  - ✅ Quantity & unit conversion
  - ✅ Session-based authentication

### Frontend Image
- **Image**: `linhtranphu/expense-frontend:latest`
- **Size**: 208MB
- **Digest**: `sha256:15459ee500abacd0f6f7b7933defbca46272a2be82ebc7312e0c67b5b8fc1af8`
- **Features**:
  - ✅ Vue.js 3 + Vite
  - ✅ Responsive design
  - ✅ Real-time expense tracking
  - ✅ User authentication UI

### Database
- **Image**: `public.ecr.aws/docker/library/mongo:7`
- **Collections**: `expenses`, `users`, `settings`

## 👥 Default Users

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | admin |
| linh | linh123 | supervisor |
| toan | toan123 | supervisor |
| yen | yen123 | supervisor |

## 🔧 Environment Variables

```env
PORT=8081
GEMINI_API_KEY=your-gemini-api-key
MONGODB_URI=mongodb://expense-mongodb:27017
SESSION_SECRET=your-secure-session-secret
```

## 🌐 EC2 Deployment Command

```bash
# Quick deploy with Docker Hub images
curl -O https://raw.githubusercontent.com/linhtranphu/CAF/main/expense-tracker/hub-deploy.sh
chmod +x hub-deploy.sh
./hub-deploy.sh
```

## 📊 Application Endpoints

- **Frontend**: http://your-ec2-ip:3000
- **Backend API**: http://your-ec2-ip:8081
- **Admin Panel**: http://your-ec2-ip:8081/admin

## ✨ New Features Added

1. **Original Message Tracking**: Admin có thể xem tin nhắn gốc
2. **User Management**: Persistent users trong MongoDB
3. **Role System**: Admin vs Supervisor roles
4. **Enhanced Parsing**: Quantity/unit conversion với base units

Images đã sẵn sàng để deploy lên EC2! 🎉