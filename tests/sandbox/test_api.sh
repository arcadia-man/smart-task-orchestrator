#!/bin/bash

# Define colors
GREEN='\03 expected expected_output\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api"

echo "======================================"
echo "🧪 Running Sandbox API E2E Tests"
echo "======================================"

# 1. Register a test user
echo "1️⃣ Creating test user..."
SIGNUP_RESP=$(curl -s -X POST "$API_URL/signup" \
  -H "Content-Type: application/json" \
  -d '{"name":"Tester","email":"tester_api_'$(date +%s)'@test.com","password":"password123"}')

# Wait, we can just login if user exists or use the previous one. Let's create a dynamic one.
USER_EMAIL="tester_api_$(date +%s)@test.com"
curl -s -X POST "$API_URL/signup" -H "Content-Type: application/json" -d "{\"name\":\"Tester\",\"email\":\"$USER_EMAIL\",\"password\":\"password123\"}" > /dev/null

# 2. Login to get token
echo "2️⃣ Logging in to get auth token..."
LOGIN_RESP=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$USER_EMAIL\",\"password\":\"password123\"}")

TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to get token. Authentication failed.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Token received!${NC}"

# 3. Define the payload
echo "3️⃣ Dispatching Sandbox Test Case..."

# The payload consists of creating the files in the container, generating the script with node, and running it with python
read -r -d '' COMMAND_SCRIPT << 'EOF'
apk add --no-cache nodejs npm python3 && \
cat << 'INNEREOF' > fibonacci.js
const fs = require('fs');
const pythonScript = `
import sys
import json
import csv

def fib(n):
    if n <= 0: return 0
    elif n == 1: return 1
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    return b

def process_file(filepath):
    try:
        with open(filepath, 'r') as f:
             data = json.load(f)
             return int(data.get('start', 0)), int(data.get('end', 0))
    except: pass
    try:
        with open(filepath, 'r') as f:
             reader = csv.reader(f)
             next(reader)
             row = next(reader)
             return int(row[0]), int(row[1])
    except: pass
    return 0, 0

if len(sys.argv) < 2: sys.exit(1)
start, end = process_file(sys.argv[1])
results = [fib(i) for i in range(start, end + 1)]
print(f"OUTPUT_{sys.argv[1]}: " + json.dumps({"start": start, "end": end, "fib": results}))
`;
fs.writeFileSync('script.py', pythonScript);
INNEREOF
cat << 'INNEREOF' > input.csv
start,end
5,15
INNEREOF
cat << 'INNEREOF' > input.json
{"start":5,"end":15}
INNEREOF

node fibonacci.js && \
python3 script.py input.csv && \
python3 script.py input.json
EOF

# Escape quotes and newlines for JSON payload
JSON_COMMAND=$(echo "$COMMAND_SCRIPT" | tr '\n' ' ' | sed 's/"/\\"/g')

JOB_PAYLOAD='{
  "name": "Fibonacci Multi-Format Test",
  "type": "one-time",
  "image": "alpine:latest",
  "command": "'"$JSON_COMMAND"'"
}'

# 4. Trigger the job
CREATE_JOB_RESP=$(curl -s -X POST "$API_URL/jobs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$JOB_PAYLOAD")

JOB_ID=$(echo $CREATE_JOB_RESP | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$JOB_ID" ]; then
    echo -e "${RED}❌ Failed to create job.${NC}"
    echo "Response: $CREATE_JOB_RESP"
    exit 1
fi
echo -e "${GREEN}✅ Job Created with ID: $JOB_ID${NC}"

# Wait for execution and check logs
echo "4️⃣ Waiting for Sandbox Execution to finish... (Polling executions DB)"

# Since there is no direct endpoint for executions, we will poll the job details endpoint
# But wait, we need to check the executions collection in MongoDB. Let's do it via the API if it's there.
# We don't have a GET /executions endpoint, but wait, we have GET /jobs/:id? It doesn't contain logs.
# Let's write a temporary script to execute python and verify the output directly, or just add the execution logs returned for test if we can't fetch it through API.
echo "⚠️  Note: Getting execution logs currently requires DB access or UI. We will fetch it from mongo directly for validation... "

sleep 10 # Wait for container to pull image and run

EXPECTED_OUTPUT_CSV='OUTPUT_input.csv: {"start": 5, "end": 15, "fib": [5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610]}'
EXPECTED_OUTPUT_JSON='OUTPUT_input.json: {"start": 5, "end": 15, "fib": [5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610]}'

# We can query mongo via docker exec
LOG_OUPUT=$(docker exec $(docker ps -q -f name=mongodb) mongosh smart_orchestrator --quiet --eval "db.executions.find({job_id: ObjectId('$JOB_ID')}).toArray()[0].logs")

if [[ "$LOG_OUPUT" == *"$EXPECTED_OUTPUT_CSV"* ]] && [[ "$LOG_OUPUT" == *"$EXPECTED_OUTPUT_JSON"* ]]; then
    echo -e "\n${GREEN}======================================${NC}"
    echo -e "${GREEN}🎉 TEST PASS: OUTPUT MATCHES EXACTLY!${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo "Execution Output:"
    echo "$LOG_OUPUT"
else
    echo -e "\n${RED}======================================${NC}"
    echo -e "${RED}❌ TEST FAIL: Output did not match.${NC}"
    echo -e "${RED}======================================${NC}"
    echo "Expected:"
    echo "$EXPECTED_OUTPUT_CSV"
    echo "$EXPECTED_OUTPUT_JSON"
    echo -e "\nActual Logs:"
    echo "$LOG_OUPUT"
fi
