const axios = require('axios');

const API_URL = "http://localhost:8080/api";
const USER_EMAIL = `tester_node_${Date.now()}@test.com`;

const sandboxCommand = `
apk add --no-cache nodejs npm python3 && \\
cat << 'INNEREOF' > fibonacci.js
const fs = require('fs');
const pythonScript = \`
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
\`;
fs.writeFileSync('script.py', pythonScript);
INNEREOF

cat << 'INNEREOF' > input.csv
start,end
5,15
INNEREOF

cat << 'INNEREOF' > input.json
{"start":5,"end":15}
INNEREOF

node fibonacci.js && \\
python3 script.py input.csv && \\
python3 script.py input.json
`;

async function runTest() {
  try {
    console.log("======================================");
    console.log("🧪 Running Sandbox API E2E Tests (NodeJS)");
    console.log("======================================");

    // 1. Signup
    console.log("1️⃣  Creating test user...");
    await axios.post(`${API_URL}/signup`, {
      name: "Tester",
      email: USER_EMAIL,
      password: "password123"
    });

    // 2. Login
    console.log("2️⃣  Logging in to get auth token...");
    const loginRes = await axios.post(`${API_URL}/login`, {
      email: USER_EMAIL,
      password: "password123"
    });
    const token = loginRes.data.token;
    console.log("✅  Token received!");

    // 3. Create Job
    console.log("3️⃣  Dispatching Sandbox Test Case...");
    const jobRes = await axios.post(`${API_URL}/jobs`, {
      name: "Fibonacci Multi-Format Test",
      type: "one-time",
      image: "alpine:latest",
      command: sandboxCommand
    }, {
      headers: { Authorization: `Bearer ${token}` }
    });
    
    const jobId = jobRes.data.id;
    console.log(`✅  Job Created with ID: ${jobId}`);

    // Wait for container to execute
    console.log("4️⃣  Waiting for Sandbox Execution to finish... (15s)");
    await new Promise(r => setTimeout(r, 15000));

    // 4. Verify MongoDB Execution Logs (using local subprocess hack to query mongo)
    const { execSync } = require('child_process');
    const logsOut = execSync(`docker exec $(docker ps -q -f name=mongodb) mongosh smart_orchestrator --quiet --eval "db.executions.find({job_id: ObjectId('${jobId}')}).toArray()[0].logs"`).toString();

    const expectedCsv = 'OUTPUT_input.csv: {"start": 5, "end": 15, "fib": [5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610]}';
    const expectedJson = 'OUTPUT_input.json: {"start": 5, "end": 15, "fib": [5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610]}';

    if (logsOut.includes(expectedCsv) && logsOut.includes(expectedJson)) {
        console.log(`\n\x1b[32m======================================\x1b[0m`);
        console.log(`\x1b[32m🎉 TEST PASS: OUTPUT MATCHES EXACTLY!\x1b[0m`);
        console.log(`\x1b[32m======================================\x1b[0m`);
        console.log("Execution Output:\n\n" + logsOut);
    } else {
        console.log(`\n\x1b[31m======================================\x1b[0m`);
        console.log(`\x1b[31m❌ TEST FAIL: Output did not match.\x1b[0m`);
        console.log(`\x1b[31m======================================\x1b[0m`);
        console.log("Expected:");
        console.log(expectedCsv);
        console.log("\nActual Logs:");
        console.log(logsOut);
    }

  } catch (error) {
    console.error("❌ Process Failed:", error.response?.data || error.message);
  }
}

runTest();
