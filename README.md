# Smart Task Orchestrator 🚀

A high-performance, scalable task orchestration platform built with Go, MongoDB, Kafka, and Docker.

## ✨ Features
- **Secure Sandbox**: Execute untrusted code in isolated Docker containers.
- **Job Types**:
  - **One-time**: Immediate execution.
  - **Cron**: Scheduled tasks.
  - **Sandbox**: Auto-scaling container pools (Min/Max).
- **Premium UI**: Modern dashboard with real-time monitoring and holistic data visualization.
- **Event Driven**: Scalable worker architecture powered by Kafka.
- **Auth**: Secure JWT-based authentication and API Key management.

## 🛠 Tech Stack
- **Backend**: Go (Gin, Kafka-go, Mongo Driver, Docker SDK)
- **Frontend**: React (Vite, Framer Motion, Recharts)
- **Infrastructure**: MongoDB, Kafka, Docker & Docker Compose

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)
- Node.js 18+ (for frontend development)

### Running with Docker Compose
```bash
docker-compose up --build
```

### Local Development

#### 1. Backend
```bash
cd backend
go mod tidy
go run cmd/api/main.go
# In separate terminals:
go run cmd/worker/main.go
go run cmd/scheduler/main.go
```

#### 2. Frontend
```bash
cd frontend
npm install
npm run dev
```

## 🏗 Architecture
See [architecture_plan.md](./architecture_plan.md) for detailed design.

## 📄 License
Open Source. Crafted with ☕ by Pritam Kumar Maurya.