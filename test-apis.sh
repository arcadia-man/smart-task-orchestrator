#!/bin/bash

# Smart Task Orchestrator - API Testing Script

set -e

API_BASE="http://localhost:8080"
TOKEN=""

echo "🧪 Testing Smart Task Orchestrator APIs"
echo "========================================"

# Function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo ""
    echo "📡 Testing: $description"
    echo "   $method $endpoint"
    
    if [ -n "$data" ]; then
        if [ -n "$TOKEN" ]; then
            curl -s -X "$method" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                -d "$data" \
                "$API_BASE$endpoint" | jq . || echo "Response: $(curl -s -X "$method" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$data" "$API_BASE$endpoint")"
        else
            curl -s -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$API_BASE$endpoint" | jq . || echo "Response: $(curl -s -X "$method" -H "Content-Type: application/json" -d "$data" "$API_BASE$endpoint")"
        fi
    else
        if [ -n "$TOKEN" ]; then
            curl -s -X "$method" \
                -H "Authorization: Bearer $TOKEN" \
                "$API_BASE$endpoint" | jq . || echo "Response: $(curl -s -X "$method" -H "Authorization: Bearer $TOKEN" "$API_BASE$endpoint")"
        else
            curl -s -X "$method" \
                "$API_BASE$endpoint" | jq . || echo "Response: $(curl -s -X "$method" "$API_BASE$endpoint")"
        fi
    fi
}

# Test 1: Health Check
api_call "GET" "/health" "" "Health Check"

# Test 2: Login (should work with placeholder)
echo ""
echo "🔐 Testing Authentication..."
LOGIN_DATA='{
    "username": "admin",
    "password": "admin"
}'

LOGIN_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA" \
    "$API_BASE/api/auth/login")

echo "Login Response: $LOGIN_RESPONSE"

# Try to extract token (will fail with placeholder, but that's expected)
if command -v jq > /dev/null 2>&1; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // empty' 2>/dev/null || echo "")
fi

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo "✅ Got access token: ${TOKEN:0:20}..."
else
    echo "⚠️  Login returned placeholder response (expected for now)"
    # Use a mock token for testing other endpoints
    TOKEN="mock-token-for-testing"
fi

# Test 3: Get current user info
api_call "GET" "/api/me" "" "Get Current User Info"

# Test 4: Get all schedulers
api_call "GET" "/api/schedulers" "" "Get All Schedulers"

# Test 5: Create a scheduler (using sample data from requirements)
echo ""
echo "📋 Testing Scheduler Creation..."
SCHEDULER_DATA='{
    "name": "Daily ETL Pipeline",
    "description": "Process daily data from multiple sources",
    "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/data-pipeline:v5",
    "jobType": "cron",
    "cronExpr": "0 2 * * *",
    "command": "bash /opt/run_etl.sh",
    "timezone": "UTC",
    "permissions": [
        {
            "roleName": "admin",
            "roleSee": true,
            "roleExecute": true,
            "roleAlterConfig": true
        }
    ]
}'

api_call "POST" "/api/schedulers" "$SCHEDULER_DATA" "Create Scheduler"

# Test 6: Get scheduler by ID (using mock ID)
api_call "GET" "/api/schedulers/507f1f77bcf86cd799439011" "" "Get Scheduler by ID"

# Test 7: Get scheduler history
api_call "GET" "/api/schedulers/507f1f77bcf86cd799439011/history" "" "Get Scheduler History"

# Test 8: Manual run scheduler
RUN_DATA='{
    "triggered_by": "manual_test"
}'
api_call "POST" "/api/schedulers/507f1f77bcf86cd799439011/run" "$RUN_DATA" "Manual Run Scheduler"

# Test 9: Get all users (admin only)
api_call "GET" "/api/users" "" "Get All Users"

# Test 10: Get all roles
api_call "GET" "/api/roles" "" "Get All Roles"

# Test 11: Get all images
api_call "GET" "/api/images" "" "Get All Images"

# Test 12: Create a new user
USER_DATA='{
    "username": "testuser",
    "email": "test@example.com",
    "password": "testpassword123",
    "roleName": "developer"
}'
api_call "POST" "/api/users" "$USER_DATA" "Create New User"

# Test 13: Create a new role
ROLE_DATA='{
    "roleName": "developer",
    "description": "Software Developer",
    "canCreateTask": true
}'
api_call "POST" "/api/roles" "$ROLE_DATA" "Create New Role"

# Test 14: Create a new image
IMAGE_DATA='{
    "name": "Node.js Runtime",
    "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/node:18",
    "description": "Node.js 18 runtime for JavaScript tasks",
    "version": "18",
    "registryUrl": "123456789012.dkr.ecr.us-east-1.amazonaws.com"
}'
api_call "POST" "/api/images" "$IMAGE_DATA" "Create New Image"

echo ""
echo "🎉 API Testing Complete!"
echo ""
echo "📝 Notes:"
echo "   - Most endpoints return placeholder responses (expected)"
echo "   - This tests the API structure and routing"
echo "   - Next step: implement actual handlers"
echo ""