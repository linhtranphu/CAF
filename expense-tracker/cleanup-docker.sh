#!/bin/bash

echo "🧹 Cleaning up unnecessary files for Docker..."

# Remove Windows batch files
rm -f cleanup.bat
rm -f run-setup.bat  
rm -f setup-simple.bat
rm -f setup.bat
rm -f start-app.bat
rm -f start-backend.bat
rm -f start-frontend.bat
rm -f test.bat

# Remove installer
rm -f go-installer.msi

# Remove unnecessary backend files
rm -f backend/package-lock.json

# Remove empty directories
rmdir backend/infrastructure/database/ 2>/dev/null || true
rmdir backend/infrastructure/sheets/ 2>/dev/null || true

# Remove duplicate docker-compose files
rm -f aws/docker-compose.prod.yml
rm -f aws/Dockerfile

echo "✅ Cleanup completed!"
echo ""
echo "📁 Files kept for Docker:"
echo "- docker-compose.yml"
echo "- docker-compose.prod.yml" 
echo "- .env"
echo "- backend/Dockerfile"
echo "- frontend/Dockerfile"
echo "- start-dev.sh"
echo "- start-prod.sh"