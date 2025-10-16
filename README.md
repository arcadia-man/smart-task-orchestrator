# Smart Task Orchestrator

A full-stack job scheduling and execution system built with Go, Kafka, MongoDB, React, and Tailwind CSS.

## Features

- ✅ Job creation and management
- ✅ Retry mechanism with exponential backoff
- ✅ Cron job scheduling
- ✅ Real-time job monitoring dashboard
- ✅ Dead Letter Queue (DLQ) for failed jobs
- ✅ REST API for job operations
- ✅ Kafka-based event pipeline

## Tech Stack

**Backend:**
- Go with Gin framework
- MongoDB for job storage
- Kafka for message queuing
- Cron scheduler for periodic jobs

**Frontend:**
- React with Vite
- Tailwind CSS for styling
- React Query for data fetching
- React Router for navigation

## Prerequisites

1. **Go 1.21+**
2. **Node.js 18+**
3. **Docker & Docker Compose**
4. **Kafka** (running locally on port 9092)

## Quick Start

### Prerequisites
1. **Kafka running on localhost:9092** (you already have this)
2. **Go 1.21+**
3. **Node.js 18+**
4. **Docker**

### Simple Commands

```bash
cd smart-task-orchestrator

# Install frontend dependencies (first time only)
cd frontend && npm install && cd ..

# Start all services
./run.sh start

# Test the API
./test.sh

# Stop all services
./run.sh stop

# Check service status
./run.sh status
```

**That's it!** 🎉

- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **MongoDB**: localhost:27017 (auto-started)
- **Kafka**: localhost:9092 (your existing setup)

## API Endpoints

- `POST /api/jobs` - Create a new job
- `GET /api/jobs` - Get all jobs
- `GET /api/jobs/:id` - Get job by ID
- `POST /api/jobs/:id/retry` - Retry a failed job
- `GET /api/jobs/:id/status` - Get job status

## Job Types

### Immediate Jobs
Execute immediately when created:

```json
{
  "name": "Process User Data",
  "type": "immediate",
  "payload": {"userId": 123},
  "maxRetries": 3
}
```

### Cron Jobs
Execute on a schedule:

```json
{
  "name": "Daily Report",
  "type": "cron",
  "cronExpr": "0 9 * * *",
  "payload": {"reportType": "daily"},
  "maxRetries": 2
}
```

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Frontend  │───▶│  API Server │───▶│   MongoDB   │
│   (React)   │    │    (Go)     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
                   ┌─────────────┐
                   │    Kafka    │
                   │             │
                   └─────────────┘
                           │
                           ▼
                   ┌─────────────┐    ┌─────────────┐
                   │   Worker    │    │  Scheduler  │
                   │ (Consumer)  │    │   (Cron)    │
                   └─────────────┘    └─────────────┘
```

## Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod tidy

# Run API server
go run cmd/api/main.go

# Run worker
go run cmd/worker/main.go

# Run scheduler
go run cmd/scheduler/main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build
```

## Testing the System

1. **Create a Job**: Use the frontend or API to create a new job
2. **Monitor Execution**: Watch the job status change in real-time
3. **Test Retries**: Create a job that will fail to see retry mechanism
4. **Cron Jobs**: Create a cron job and watch it execute periodically

## Configuration

Environment variables (`.env`):

```env
MONGO_URI=mongodb://localhost:27017/orchestrator
KAFKA_BROKER=localhost:9092
DB_NAME=orchestrator
PORT=8080
```

## Troubleshooting

1. **Kafka Connection Issues**: Ensure Kafka is running on localhost:9092
2. **MongoDB Connection**: Check if MongoDB container is running
3. **Port Conflicts**: Make sure ports 8080, 27017, and 3000 are available

## Next Steps

- Add authentication and authorization
- Implement job dependencies
- Add more sophisticated scheduling options
- Implement job result storage
- Add metrics and monitoring
- Implement job cancellation