# Smart Task Orchestrator

A production-ready distributed job scheduling system with real-time monitoring, RBAC, and Docker container execution.

## 🚀 Quick Start

```bash
# Clone and start
git clone <your-repo>
cd smart-task-orchestrator

# Start all services
docker-compose up --build

# Access the application
# Frontend: http://localhost:3000
# API: http://localhost:8080
```

## 🏗️ Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Frontend  │───▶│  API Server │───▶│   MongoDB   │
│   (React)   │    │    (Go)     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
                   ┌─────────────┐    ┌─────────────┐
                   │    Redis    │    │    Kafka    │
                   │  (Sharded)  │    │             │
                   └─────────────┘    └─────────────┘
                           │                   │
                           ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐
                   │  Scheduler  │    │   Worker    │
                   │  Library    │    │ (Container) │
                   └─────────────┘    └─────────────┘
```

## 🔧 Services

- **API Server**: REST endpoints, authentication, WebSocket for logs
- **Scheduler**: Precompute + poller for exact timing (1-5s precision)
- **Worker**: Docker container execution with real-time log streaming
- **Frontend**: React dashboard with live monitoring

## 🔐 Default Credentials

- **Username**: `admin`
- **Password**: `admin` (must change on first login)

## 📊 Features

- ✅ Sub-second job scheduling precision
- ✅ Real-time log streaming with ANSI colors
- ✅ Role-based access control (RBAC)
- ✅ Docker container execution
- ✅ Horizontal scalability
- ✅ Production-ready deployment

## 🛠️ Development

```bash
# Backend development
cd backend
go mod tidy
go run cmd/api/main.go

# Frontend development  
cd frontend
npm install
npm run dev
```

## 📝 Configuration

Environment variables in `docker-compose.yml`:

- `MONGO_URI`: MongoDB connection string
- `REDIS_URL`: Redis connection string  
- `KAFKA_BROKER`: Kafka broker address
- `JWT_SECRET`: JWT signing secret (change in production)

## 🔍 Monitoring

- **Logs**: `docker-compose logs -f <service>`
- **Health**: API endpoints at `/health`
- **Metrics**: Prometheus endpoints at `/metrics`

## 📚 API Documentation

- `POST /api/auth/login` - User authentication
- `GET /api/schedulers` - List all schedulers
- `POST /api/schedulers` - Create new scheduler
- `POST /api/schedulers/:id/run` - Manual job execution
- `GET /ws/logs/:runId` - WebSocket log streaming

## 🚀 Production Deployment

1. Update environment variables in `docker-compose.yml`
2. Set strong JWT secret
3. Configure proper MongoDB and Redis persistence
4. Set up reverse proxy (nginx) for SSL termination
5. Configure monitoring and alerting

## 🔧 Troubleshooting

```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs -f api
docker-compose logs -f scheduler
docker-compose logs -f worker

# Restart services
docker-compose restart <service>
```

## 📄 License

MIT License - see LICENSE file for details.