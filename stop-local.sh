#!/bin/bash

# Smart Task Orchestrator - Stop Local Services

echo "🛑 Stopping Smart Task Orchestrator services..."

# Function to stop a service
stop_service() {
    local service_name=$1
    local pid_file="logs/${service_name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            echo "🛑 Stopping $service_name (PID: $pid)..."
            kill "$pid"
            sleep 1
            
            # Force kill if still running
            if kill -0 "$pid" 2>/dev/null; then
                echo "⚠️  Force killing $service_name..."
                kill -9 "$pid"
            fi
        else
            echo "⚠️  $service_name was not running"
        fi
        rm -f "$pid_file"
    else
        echo "⚠️  No PID file found for $service_name"
    fi
}

# Stop all services
stop_service "api"
stop_service "scheduler"
stop_service "worker"

# Kill any remaining processes on port 8080
echo "🔍 Checking for processes on port 8080..."
if lsof -ti:8080 > /dev/null 2>&1; then
    echo "🛑 Killing processes on port 8080..."
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
fi

echo "✅ All services stopped"