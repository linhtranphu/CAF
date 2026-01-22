#!/bin/bash

case "$1" in
  start)
    echo "🚀 Starting services..."
    docker-compose -f docker-compose.prod.yml up -d
    ;;
  stop)
    echo "🛑 Stopping services..."
    docker-compose -f docker-compose.prod.yml down
    ;;
  restart)
    echo "🔄 Restarting services..."
    docker-compose -f docker-compose.prod.yml restart
    ;;
  logs)
    echo "📋 Showing logs..."
    docker-compose -f docker-compose.prod.yml logs -f --tail=50
    ;;
  status)
    echo "📊 Container status:"
    docker-compose -f docker-compose.prod.yml ps
    ;;
  clean)
    echo "🧹 Cleaning up..."
    docker-compose -f docker-compose.prod.yml down -v
    docker system prune -f
    ;;
  update)
    echo "🔄 Updating application..."
    git pull origin main
    docker-compose -f docker-compose.prod.yml up -d --build
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|logs|status|clean|update}"
    echo ""
    echo "Commands:"
    echo "  start   - Start all services"
    echo "  stop    - Stop all services"
    echo "  restart - Restart all services"
    echo "  logs    - Show live logs"
    echo "  status  - Show container status"
    echo "  clean   - Stop and remove all containers/volumes"
    echo "  update  - Pull latest code and rebuild"
    exit 1
    ;;
esac