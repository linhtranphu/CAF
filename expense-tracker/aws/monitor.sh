#!/bin/bash
# Monitoring and Maintenance Script

APP_DIR="/home/ec2-user/expense-tracker"
BACKUP_DIR="/home/ec2-user/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Function to backup MongoDB
backup_mongodb() {
    echo "🗄️ Backing up MongoDB..."
    docker exec expense-mongodb-prod mongodump --out /backup/mongodb_$DATE
    cp -r $APP_DIR/mongodb-backup/mongodb_$DATE $BACKUP_DIR/
    echo "✅ MongoDB backup completed: $BACKUP_DIR/mongodb_$DATE"
}

# Function to check service health
check_health() {
    echo "🏥 Checking service health..."
    
    # Backend health
    if curl -f http://localhost:8081/api/health > /dev/null 2>&1; then
        echo "✅ Backend: Healthy"
    else
        echo "❌ Backend: Unhealthy"
        docker-compose -f $APP_DIR/aws/docker-compose.prod.yml restart backend
    fi
    
    # Frontend health
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ Frontend: Healthy"
    else
        echo "❌ Frontend: Unhealthy"
        docker-compose -f $APP_DIR/aws/docker-compose.prod.yml restart frontend
    fi
    
    # MongoDB health
    if docker exec expense-mongodb-prod mongo --eval "db.adminCommand('ismaster')" > /dev/null 2>&1; then
        echo "✅ MongoDB: Healthy"
    else
        echo "❌ MongoDB: Unhealthy"
        docker-compose -f $APP_DIR/aws/docker-compose.prod.yml restart mongodb
    fi
}

# Function to show system stats
show_stats() {
    echo "📊 System Statistics:"
    echo "CPU Usage: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%"
    echo "Memory Usage: $(free | grep Mem | awk '{printf("%.2f%%", $3/$2 * 100.0)}')"
    echo "Disk Usage: $(df -h / | awk 'NR==2{printf "%s", $5}')"
    echo ""
    echo "🐳 Docker Container Status:"
    docker-compose -f $APP_DIR/aws/docker-compose.prod.yml ps
}

# Function to clean old backups (keep last 7 days)
cleanup_backups() {
    echo "🧹 Cleaning old backups..."
    find $BACKUP_DIR -name "mongodb_*" -type d -mtime +7 -exec rm -rf {} \;
    echo "✅ Cleanup completed"
}

# Main menu
case "$1" in
    "backup")
        backup_mongodb
        ;;
    "health")
        check_health
        ;;
    "stats")
        show_stats
        ;;
    "cleanup")
        cleanup_backups
        ;;
    "all")
        check_health
        show_stats
        backup_mongodb
        cleanup_backups
        ;;
    *)
        echo "Usage: $0 {backup|health|stats|cleanup|all}"
        echo ""
        echo "Commands:"
        echo "  backup  - Backup MongoDB database"
        echo "  health  - Check service health and restart if needed"
        echo "  stats   - Show system and container statistics"
        echo "  cleanup - Remove old backups (>7 days)"
        echo "  all     - Run all commands"
        exit 1
        ;;
esac