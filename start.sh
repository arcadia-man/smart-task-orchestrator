#!/bin/bash

# Smart Task Orchestrator - Complete Start Script

set -e

echo "🚀 Starting Smart Task Orchestrator..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker compose is available
if ! docker compose version > /dev/null 2>&1; then
    echo "❌ docker compose is not available. Please install Docker Compose first."
    exit 1
fi

# Build and start all services
echo "📦 Building and starting all services..."
docker compose up --build -d

echo "⏳ Waiting for services to initialize..."
sleep 15

# Check service health
echo "🔍 Checking service health..."

# Function to check service health
check_service() {
    local service_name=$1
    local url=$2
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -f "$url" > /dev/null 2>&1; then
            echo "✅ $service_name is running"
            return 0
        fi
        echo "⏳ Waiting for $service_name... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    echo "⚠️  $service_name is not responding after $max_attempts attempts"
    return 1
}

# Check API
check_service "API Server" "http://localhost:8080/health"

# Check Frontend
check_service "Frontend" "http://localhost:3000"

# Check infrastructure services
echo "🔍 Checking infrastructure services..."

# MongoDB
if docker compose exec -T mongodb mongosh --eval "db.runCommand('ping')" > /dev/null 2>&1; then
    echo "✅ MongoDB is running"
else
    echo "⚠️  MongoDB might still be starting..."
fi

# Redis
if docker compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is running"
else
    echo "⚠️  Redis might still be starting..."
fi

# Kafka
if docker compose exec -T kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
    echo "✅ Kafka is running"
else
    echo "⚠️  Kafka might still be starting..."
fi

echo ""
echo "🎉 Smart Task Orchestrator is ready!"
echo ""
echo "📱 Frontend: http://localhost:3000"
echo "🔧 API: http://localhost:8080"
echo "📊 Health Check: http://localhost:8080/health"
echo "📊 MongoDB: localhost:27017"
echo "🔴 Redis: localhost:6379"
echo "📨 Kafka: localhost:9092"
echo ""
echo "🔐 Default Login:"
echo "   Username: admin"
echo "   Password: admin (change on first login)"
echo ""
echo "📝 Useful Commands:"
echo "   View logs: docker compose logs -f [service]"
echo "   Stop all: docker compose down"
echo "   Restart: docker compose restart [service]"
echo "   Shell access: docker compose exec [service] /bin/sh"
echo ""
echo "🚀 Services Status:"
docker compose ps
echo ""