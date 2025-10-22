#!/bin/bash

# Smart Task Orchestrator - Frontend Development Script

set -e

echo "🚀 Starting Smart Task Orchestrator Frontend"

# Check if Node.js is available
if ! command -v node > /dev/null 2>&1; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if npm is available
if ! command -v npm > /dev/null 2>&1; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

echo "✅ Node.js version: $(node --version)"
echo "✅ npm version: $(npm --version)"

# Navigate to frontend directory
cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    npm install
else
    echo "✅ Dependencies already installed"
fi

# Check if API server is running
echo ""
echo "🔍 Checking if API server is running..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API Server is running on http://localhost:8080"
else
    echo "⚠️  API Server is not running. Please start it first:"
    echo "   ./run-local.sh"
    echo ""
    echo "🔄 Starting frontend anyway (will show connection errors until API is started)..."
fi

echo ""
echo "🚀 Starting frontend development server..."
echo "📱 Frontend will be available at: http://localhost:3000"
echo "🔧 API proxy configured to: http://localhost:8080"
echo ""
echo "📝 Press Ctrl+C to stop the frontend server"
echo ""
echo "🔍 Debug mode enabled - check browser console for detailed logs"
echo "   - Login flow: 🔐 LOGIN messages"
echo "   - Password change: 🔐 PASSWORD messages"
echo "   - Modal state: 🔐 MODAL messages"
echo "   - Auth state: 🔐 AUTH messages"
echo ""

# Start the development server
npm run dev