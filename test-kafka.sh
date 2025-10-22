#!/bin/bash

# Test Kafka integration

echo "📨 Testing Kafka Integration"
echo "============================"

# Test message based on your requirements
TEST_MESSAGE='{
    "run_id": "test-run-123",
    "scheduler_id": "test-scheduler-456", 
    "generation": 1,
    "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/ubuntu:latest",
    "command": "echo \"Hello from Smart Task Orchestrator!\" && sleep 2 && echo \"Job completed successfully\"",
    "env": {"TEST_ENV": "development"},
    "triggered_by": "manual_test",
    "run_at": "2025-10-22T19:10:00Z",
    "metadata": {
        "test": true,
        "scheduler_name": "Test Job"
    }
}'

echo ""
echo "📤 Sending test job execution message to Kafka..."
echo "Topic: job_executions"
echo "Message: $TEST_MESSAGE"

# Send message to Kafka
echo "$TEST_MESSAGE" | kafka-console-producer --bootstrap-server localhost:9092 --topic job_executions

echo ""
echo "✅ Message sent to Kafka!"
echo ""
echo "📋 Check worker logs to see if it processes the message:"
echo "   tail -f logs/worker.log"
echo ""
echo "🔍 You should see the worker:"
echo "   1. Consume the message from Kafka"
echo "   2. Create a history record in MongoDB"
echo "   3. Execute the Docker container"
echo "   4. Stream logs to Redis"
echo "   5. Update final status"