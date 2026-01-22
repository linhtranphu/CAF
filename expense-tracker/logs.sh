#!/bin/bash

echo "📊 Viewing Expense Tracker logs..."
echo "Press Ctrl+C to exit"
echo ""

# Follow logs from all services
docker-compose logs -f