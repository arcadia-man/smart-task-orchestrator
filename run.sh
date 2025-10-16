#!/bin/bash

# Smart Task Orchestrator - Simple Run Script

set -e  # Exit on any error

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a port is in use
check_port() {
    local port=$1
    if lsof -ti:$port > /dev/null 2>&1; then
        return 0  # Port is in use
    else
        return 1  # Port is free
    fi
}

# Function to wait for service to be ready
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1
    
    print_status "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to start within $max_attempts seconds"
    return 1
}

# Function to start services
start_services() {
    print_status "Starting Smart Task Orchestrator..."
    
    # Check prerequisites
    print_status "Checking prerequisites..."
    
    # Check if Kafka is running
    if ! nc -z localhost 9092 2>/dev/null; then
        print_error "Kafka is not running on localhost:9092"
        print_error "Please start Kafka first:"
        print_error "1. Start Zookeeper: bin/zookeeper-server-start.sh config/zookeeper.properties"
        print_error "2. Start Kafka: bin/kafka-server-start.sh config/server.properties"
        exit 1
    fi
    print_success "Kafka is running"
    
    # Create Kafka topics if they don't exist
    print_status "Creating Kafka topics..."
    kafka-topics --create --topic jobs.execute --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 --if-not-exists > /dev/null 2>&1 || true
    kafka-topics --create --topic jobs.failed --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1 --if-not-exists > /dev/null 2>&1 || true
    print_success "Kafka topics ready"
    
    # Start MongoDB
    print_status "Starting MongoDB..."
    if ! docker ps | grep -q orchestrator-mongo; then
        docker run -d --name orchestrator-mongo -p 27017:27017 mongo:7 > /dev/null 2>&1
        sleep 5
    fi
    print_success "MongoDB is running"
    
    # Check if ports are available
    if check_port 8080; then
        print_error "Port 8080 is already in use. Please stop the service using that port."
        exit 1
    fi
    
    if check_port 3000; then
        print_error "Port 3000 is already in use. Please stop the service using that port."
        exit 1
    fi
    
    # Start backend services
    print_status "Starting backend services..."
    
    cd backend
    
    # Start API server
    print_status "Starting API server..."
    nohup go run cmd/api/main.go > ../logs/api.log 2>&1 &
    API_PID=$!
    echo $API_PID > ../pids/api.pid
    
    # Wait for API to be ready
    if ! wait_for_service "http://localhost:8080/api/jobs" "API Server"; then
        print_error "Failed to start API server"
        exit 1
    fi
    
    # Start worker
    print_status "Starting worker..."
    nohup go run cmd/worker/main.go > ../logs/worker.log 2>&1 &
    WORKER_PID=$!
    echo $WORKER_PID > ../pids/worker.pid
    
    # Start scheduler
    print_status "Starting scheduler..."
    nohup go run cmd/scheduler/main.go > ../logs/scheduler.log 2>&1 &
    SCHEDULER_PID=$!
    echo $SCHEDULER_PID > ../pids/scheduler.pid
    
    cd ..
    
    # Start frontend
    print_status "Starting frontend..."
    cd frontend
    nohup npm run dev > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > ../pids/frontend.pid
    cd ..
    
    # Wait for frontend to be ready
    if ! wait_for_service "http://localhost:3000" "Frontend"; then
        print_warning "Frontend may take a moment to start"
    fi
    
    print_success "All services started successfully!"
    echo ""
    echo "📊 Services:"
    echo "  - API Server: http://localhost:8080"
    echo "  - Frontend: http://localhost:3000"
    echo "  - MongoDB: localhost:27017"
    echo "  - Kafka: localhost:9092"
    echo ""
    echo "📋 Test the API:"
    echo "  curl http://localhost:8080/api/jobs"
    echo ""
    echo "🛑 To stop all services:"
    echo "  ./run.sh stop"
    echo ""
    echo "📝 View logs:"
    echo "  tail -f logs/api.log"
    echo "  tail -f logs/worker.log"
    echo "  tail -f logs/scheduler.log"
    echo "  tail -f logs/frontend.log"
}

# Function to stop services
stop_services() {
    print_status "Stopping Smart Task Orchestrator..."
    
    # Stop processes using PID files
    for service in api worker scheduler frontend; do
        if [ -f "pids/$service.pid" ]; then
            PID=$(cat "pids/$service.pid")
            if kill -0 $PID 2>/dev/null; then
                print_status "Stopping $service (PID: $PID)..."
                kill $PID
                rm -f "pids/$service.pid"
            else
                print_warning "$service was not running"
                rm -f "pids/$service.pid"
            fi
        fi
    done
    
    # Stop MongoDB container
    print_status "Stopping MongoDB..."
    docker stop orchestrator-mongo > /dev/null 2>&1 || true
    docker rm orchestrator-mongo > /dev/null 2>&1 || true
    
    # Kill any remaining processes
    pkill -f "go run cmd/api/main.go" > /dev/null 2>&1 || true
    pkill -f "go run cmd/worker/main.go" > /dev/null 2>&1 || true
    pkill -f "go run cmd/scheduler/main.go" > /dev/null 2>&1 || true
    pkill -f "npm run dev" > /dev/null 2>&1 || true
    
    print_success "All services stopped"
}

# Function to show status
show_status() {
    print_status "Service Status:"
    echo ""
    
    # Check API
    if curl -s http://localhost:8080/api/jobs > /dev/null 2>&1; then
        print_success "API Server: Running (http://localhost:8080)"
    else
        print_error "API Server: Not running"
    fi
    
    # Check Frontend
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        print_success "Frontend: Running (http://localhost:3000)"
    else
        print_error "Frontend: Not running"
    fi
    
    # Check MongoDB
    if docker ps | grep -q orchestrator-mongo; then
        print_success "MongoDB: Running (localhost:27017)"
    else
        print_error "MongoDB: Not running"
    fi
    
    # Check Kafka
    if nc -z localhost 9092 2>/dev/null; then
        print_success "Kafka: Running (localhost:9092)"
    else
        print_error "Kafka: Not running"
    fi
}

# Create necessary directories
mkdir -p logs pids

# Main script logic
case "${1:-start}" in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        stop_services
        sleep 2
        start_services
        ;;
    status)
        show_status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        echo ""
        echo "Commands:"
        echo "  start   - Start all services (default)"
        echo "  stop    - Stop all services"
        echo "  restart - Restart all services"
        echo "  status  - Show service status"
        exit 1
        ;;
esac