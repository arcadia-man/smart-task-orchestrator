#!/bin/bash

# Smart Task Orchestrator - Simple Start Script

set -e

echo "🚀 Starting Smart Task Orchestrator..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose > /dev/null 2>&1; then
    echo "❌ docker-compose is not installed. Please install docker-compose first."
    exit 1
fi

# Start all services
echo "📦 Building and starting all services..."
docker-compose up --build -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check service health
echo "🔍 Checking service health..."

# Check API
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API Server is running"
else
    echo "⚠️  API Server might still be starting..."
fi

# Check Frontend
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend is running"
else
    echo "⚠️  Frontend might still be starting..."
fi

echo ""
echo "🎉 Smart Task Orchestrator is starting up!"
echo ""
echo "📱 Frontend: http://localhost:3000"
echo "🔧 API: http://localhost:8080"
echo "📊 MongoDB: localhost:27017"
echo "🔴 Redis: localhost:6379"
echo "📨 Kafka: localhost:9092"
echo ""
echo "🔐 Default Login:"
echo "   Username: admin"
echo "   Password: admin (change on first login)"
echo ""
echo "📝 View logs: docker-compose logs -f"
echo "🛑 Stop: docker-compose down"
echo ""