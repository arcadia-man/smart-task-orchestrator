#!/bin/bash

# Simple test script for Smart Task Orchestrator

API_BASE="http://localhost:8080/api"

echo "🧪 Testing Smart Task Orchestrator API..."

# Test 1: Check if API is running
echo "1. Checking API health..."
if curl -s "$API_BASE/jobs" > /dev/null; then
    echo "✅ API is running"
else
    echo "❌ API is not running. Please start it with: ./run.sh start"
    exit 1
fi

# Test 2: Create an immediate job
echo ""
echo "2. Creating immediate job..."
JOB_RESPONSE=$(curl -s -X POST "$API_BASE/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Immediate Job",
    "type": "immediate",
    "payload": {"message": "Hello World", "userId": 123},
    "maxRetries": 3
  }')

if echo "$JOB_RESPONSE" | grep -q '"id"'; then
    echo "✅ Job created successfully"
    JOB_ID=$(echo $JOB_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "   Job ID: $JOB_ID"
else
    echo "❌ Failed to create job"
    echo "   Response: $JOB_RESPONSE"
    exit 1
fi

# Test 3: Get all jobs
echo ""
echo "3. Getting all jobs..."
JOBS_RESPONSE=$(curl -s "$API_BASE/jobs")
JOB_COUNT=$(echo "$JOBS_RESPONSE" | grep -o '"id"' | wc -l)
echo "✅ Found $JOB_COUNT jobs"

# Test 4: Get specific job
if [ ! -z "$JOB_ID" ]; then
    echo ""
    echo "4. Getting job details..."
    JOB_DETAIL=$(curl -s "$API_BASE/jobs/$JOB_ID")
    if echo "$JOB_DETAIL" | grep -q '"status"'; then
        STATUS=$(echo "$JOB_DETAIL" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
        echo "✅ Job status: $STATUS"
    else
        echo "❌ Failed to get job details"
    fi
fi

# Test 5: Create a cron job
echo ""
echo "5. Creating cron job..."
CRON_RESPONSE=$(curl -s -X POST "$API_BASE/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Cron Job",
    "type": "cron",
    "cronExpr": "*/5 * * * *",
    "payload": {"task": "periodic cleanup"},
    "maxRetries": 2
  }')

if echo "$CRON_RESPONSE" | grep -q '"id"'; then
    echo "✅ Cron job created successfully"
else
    echo "❌ Failed to create cron job"
    echo "   Response: $CRON_RESPONSE"
fi

echo ""
echo "🎉 API tests completed!"
echo ""
echo "📊 Access the dashboard: http://localhost:3000"
echo "📋 View all jobs: curl $API_BASE/jobs | jq '.'"