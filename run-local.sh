#!/bin/bash

# Smart Task Orchestrator - Local Development Script
# Uses host Redis, Kafka, MongoDB

set -e

echo "🚀 Starting Smart Task Orchestrator (Local Development Mode)"

# Check if required services are running on host
echo "🔍 Checking host services..."

# Check MongoDB
if ! nc -z localhost 27017 2>/dev/null; then
    echo "❌ MongoDB is not running on localhost:27017"
    echo "   Please start MongoDB first"
    exit 1
fi
echo "✅ MongoDB is running on localhost:27017"

# Check Redis
if ! nc -z localhost 6379 2>/dev/null; then
    echo "❌ Redis is not running on localhost:6379"
    echo "   Please start Redis first"
    exit 1
fi
echo "✅ Redis is running on localhost:6379"

# Check Kafka
if ! nc -z localhost 9092 2>/dev/null; then
    echo "❌ Kafka is not running on localhost:9092"
    echo "   Please start Kafka first"
    exit 1
fi
echo "✅ Kafka is running on localhost:9092"

# Check Docker
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running"
    echo "   Please start Docker first"
    exit 1
fi
echo "✅ Docker is running"

echo ""
echo "📦 Building Go applications..."

# Build backend applications
cd backend

# Build all binaries
echo "🔨 Building API server..."
go build -o bin/api ./cmd/api

echo "🔨 Building Scheduler service..."
go build -o bin/scheduler ./cmd/scheduler

echo "🔨 Building Worker service..."
go build -o bin/worker ./cmd/worker

echo "✅ All binaries built successfully"

# Create logs directory
mkdir -p ../logs

echo ""
echo "🚀 Starting services..."

# Start API server in background
echo "🔧 Starting API server on port 8080..."
./bin/api > ../logs/api.log 2>&1 &
API_PID=$!
echo $API_PID > ../logs/api.pid

# Wait a moment for API to start
sleep 2

# Start Scheduler service in background
echo "⏰ Starting Scheduler service..."
./bin/scheduler > ../logs/scheduler.log 2>&1 &
SCHEDULER_PID=$!
echo $SCHEDULER_PID > ../logs/scheduler.pid

# Start Worker service in background
echo "🔄 Starting Worker service..."
./bin/worker > ../logs/worker.log 2>&1 &
WORKER_PID=$!
echo $WORKER_PID > ../logs/worker.pid

cd ..

# Wait for services to start
echo "⏳ Waiting for services to initialize..."
sleep 3

# Check if services are running
echo "🔍 Checking service health..."

# Check API health
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API Server is running (PID: $API_PID)"
else
    echo "⚠️  API Server might still be starting..."
fi

# Check if processes are still running
if kill -0 $SCHEDULER_PID 2>/dev/null; then
    echo "✅ Scheduler service is running (PID: $SCHEDULER_PID)"
else
    echo "❌ Scheduler service failed to start"
fi

if kill -0 $WORKER_PID 2>/dev/null; then
    echo "✅ Worker service is running (PID: $WORKER_PID)"
else
    echo "❌ Worker service failed to start"
fi

echo ""
echo "🎉 Smart Task Orchestrator is running!"
echo ""
echo "🔧 API Server: http://localhost:8080"
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
echo "   View API logs: tail -f logs/api.log"
echo "   View Scheduler logs: tail -f logs/scheduler.log"
echo "   View Worker logs: tail -f logs/worker.log"
echo "   Stop all: ./stop-local.sh"
echo ""
echo "📊 Process IDs:"
echo "   API: $API_PID"
echo "   Scheduler: $SCHEDULER_PID"
echo "   Worker: $WORKER_PID"
echo ""
echo "🚀 Ready for testing!"