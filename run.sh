#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

start_services() {
    print_info "Starting Smart Task Orchestrator..."
    
    # Start MongoDB and Redis
    print_info "Starting MongoDB and Redis..."
    docker compose up -d mongodb redis
    
    # Build and start backend services
    print_info "Building and starting backend services..."
    docker compose up -d --build api scheduler worker
    
    # Build and start frontend
    print_info "Building and starting frontend..."
    docker compose up -d --build frontend
    
    print_success "Services started!"
    echo ""
    echo "Access:"
    echo "  Frontend: http://localhost:3000"
    echo "  API: http://localhost:8080"
    echo ""
    echo "Login: admin / 12345"
}

stop_services() {
    print_info "Stopping services..."
    docker compose down
    print_success "Stopped"
}

show_status() {
    docker compose ps
}

show_logs() {
    if [ -z "$1" ]; then
        docker compose logs -f
    else
        docker compose logs -f "$1"
    fi
}

case "${1:-}" in
    start) start_services ;;
    stop) stop_services ;;
    restart) stop_services && sleep 2 && start_services ;;
    status) show_status ;;
    logs) show_logs "${2:-}" ;;
    *)
        echo "Usage: ./run.sh [start|stop|restart|status|logs]"
        echo ""
        echo "Examples:"
        echo "  ./run.sh start"
        echo "  ./run.sh logs api"
        echo "  ./run.sh status"
        ;;
esac
