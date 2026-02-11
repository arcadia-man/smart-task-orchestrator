# Smart Task Orchestrator - Complete Requirements & Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Core Requirements](#core-requirements)
4. [Implementation Status](#implementation-status)
5. [API Documentation](#api-documentation)
6. [Deployment Guide](#deployment-guide)
7. [Testing Guide](#testing-guide)

---

## Project Overview

A production-ready distributed job scheduling system with real-time monitoring, RBAC, and Docker container execution.

### Key Features
- ✅ Sub-second job scheduling precision (1-5 seconds)
- ✅ Real-time log streaming with ANSI colors
- ✅ Role-based access control (RBAC)
- ✅ Docker container execution
- ✅ Horizontal scalability
- ✅ Complete admin interface with real API integration

### Technology Stack
- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: React 18+ with TypeScript and Tailwind CSS
- **Database**: MongoDB 7+ for document storage
- **Cache**: Redis 7+ with sharding support
- **Messaging**: Kafka for event streaming
- **Containers**: Docker for job execution

---

## Architecture

### System Architecture
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

### Service Components

#### 1. API Server (`cmd/api`)
- REST endpoints for CRUD operations
- JWT authentication with RBAC
- WebSocket for real-time log streaming
- Admin user seeding and management
- Health checks and monitoring

#### 2. Scheduler Library (`cmd/scheduler`)
- **Precompute Service**: Runs every 5 minutes, calculates next 15 minutes of executions
- **Poller Service**: 1-second precision polling with Redis sharding
- **Generation System**: Prevents stale job execution after configuration changes
- **Duplicate Prevention**: Atomic Redis operations + database constraints

#### 3. Worker Service (`cmd/worker`)
- Kafka consumer for job execution
- Docker container management
- Real-time log streaming with ANSI preservation
- Process lifecycle management
- Concurrent execution with semaphore-based limiting

#### 4. Frontend (`frontend/`)
- React dashboard with real-time updates
- Job creation with advanced scheduling
- Live log viewer with color preservation
- User and permission management
- Responsive design for all screen sizes

---

## Core Requirements

### Functional Requirements

#### 1. Job Scheduling
- **Immediate Jobs**: Execute right away
- **Cron Jobs**: Schedule with cron expressions (e.g., `0 2 * * *`)
- **Interval Jobs**: Periodic execution (e.g., every 300 seconds)
- **Timezone Support**: Configurable timezone for scheduling

#### 2. Distributed Execution
- Multi-worker Docker container execution
- Host Docker integration (not Docker-in-Docker)
- Concurrent job execution with resource limits
- Container lifecycle management (pull, run, monitor, cleanup)

#### 3. Real-time Monitoring
- Live log streaming with ANSI color preservation
- WebSocket + Redis streams for real-time updates
- Job status tracking (pending, running, success, failed)
- Execution history and audit trail

#### 4. RBAC (Role-Based Access Control)
- User authentication with JWT
- Role management with permissions
- Fine-grained per-scheduler access control
- Auto-role creation from UI

#### 5. Exact Timing
- 1-5 second precision for job execution
- Precompute + poller architecture
- Generation-based invalidation
- Distributed coordination with Redis

### Technical Requirements

#### Backend (Go)
- Gin framework for REST API
- MongoDB driver for database operations
- Kafka client for event streaming
- Redis client with sharding support
- Docker SDK for container management
- JWT for authentication
- Bcrypt for password hashing

#### Frontend (React)
- TypeScript for type safety
- Tailwind CSS for styling
- React Query for data fetching
- React Router for navigation
- WebSocket for real-time updates
- Toast notifications for user feedback

#### Database (MongoDB)
- **Collections**:
  - `users` - User accounts and authentication
  - `roles` - Role definitions and permissions
  - `scheduler_definitions` - Job configurations
  - `scheduler_history` - Execution history and logs
  - `scheduler_precompute` - Scheduled execution times
  - `images` - Docker image registry
  - `audit_logs` - System audit trail

#### Cache (Redis)
- Sharded sorted sets for scheduled items
- Real-time log streaming with Redis streams
- Distributed locking for coordination
- Metadata caching for performance

#### Messaging (Kafka)
- `job_executions` topic for job dispatch
- Event-driven architecture
- Reliable message delivery
- Horizontal scaling support

---

## Implementation Status

### ✅ Phase 1: Foundation & Architecture (COMPLETE)

#### Backend Services
- ✅ API Server with REST endpoints
- ✅ JWT authentication with RBAC middleware
- ✅ Scheduler with precompute + poller
- ✅ Worker with Docker execution
- ✅ MongoDB integration with seeding
- ✅ Redis sharding setup
- ✅ Kafka producer/consumer

#### Frontend Components
- ✅ Login with authentication
- ✅ Dashboard with statistics
- ✅ Scheduler management (CRUD)
- ✅ User management with roles
- ✅ Role and permission management
- ✅ Docker image management
- ✅ System logs with filtering
- ✅ Monitoring dashboard

#### Security Features
- ✅ JWT-based authentication
- ✅ Role-based access control
- ✅ Password change on initial login
- ✅ Secure password hashing
- ✅ Route protection (frontend & backend)
- ✅ Token persistence in localStorage

#### Data Models
- ✅ Users with role assignment
- ✅ Roles with permissions
- ✅ Scheduler definitions
- ✅ Execution history
- ✅ Docker images
- ✅ Audit logs

### 🔄 Phase 2: API Integration (COMPLETE)

#### Real API Implementation
- ✅ User Management APIs (CRUD + password reset)
- ✅ Role Management APIs (CRUD + permissions)
- ✅ Scheduler Management APIs (CRUD + run + history)
- ✅ Image Management APIs (CRUD + usage tracking)
- ✅ System Logs APIs (fetch + filter + stats)
- ✅ Monitoring APIs (metrics + services + alerts)
- ✅ Dashboard Stats API

#### Frontend Integration
- ✅ Removed all dummy data
- ✅ Connected to real APIs
- ✅ Loading states for async operations
- ✅ Error handling with user feedback
- ✅ Real-time data updates
- ✅ Form validation
- ✅ Toast notifications

### 🚀 Phase 3: Production Features (IN PROGRESS)

#### Deployment
- ✅ Docker Compose orchestration
- ✅ Multi-stage Docker builds
- ✅ Environment configuration
- ✅ Health checks
- ⏳ Kubernetes manifests
- ⏳ CI/CD pipeline

#### Monitoring & Observability
- ✅ Health check endpoints
- ✅ Service status monitoring
- ⏳ Prometheus metrics
- ⏳ Grafana dashboards
- ⏳ Alerting system

#### Advanced Features
- ⏳ Job dependencies and workflows
- ⏳ Webhook notifications
- ⏳ Job result storage
- ⏳ Advanced scheduling (timezone support)
- ⏳ Job cancellation
- ⏳ Resource quotas and limits

---

## API Documentation

### Authentication Endpoints

#### POST /api/auth/login
Login with username and password.

**Request:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "must_change_password": true,
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "admin",
    "email": "admin@orchestrator.local",
    "roleId": "507f1f77bcf86cd799439012",
    "isInitialLogin": true
  }
}
```

#### POST /api/auth/refresh
Refresh access token using refresh token.

#### POST /api/auth/change-password
Change user password (requires authentication).

**Request:**
```json
{
  "old_password": "admin",
  "new_password": "newpassword123"
}
```

#### GET /api/me
Get current user information (requires authentication).

### Scheduler Endpoints

#### GET /api/schedulers
Get all schedulers with execution information.

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Daily ETL Pipeline",
    "description": "Process daily data",
    "image": "python:3.11-slim",
    "jobType": "cron",
    "cronExpr": "0 2 * * *",
    "command": "python /app/etl.py",
    "status": "active",
    "lastRun": "2025-11-16T02:00:00Z",
    "nextRun": "2025-11-17T02:00:00Z",
    "lastStatus": "success"
  }
]
```

#### POST /api/schedulers
Create new scheduler.

**Request:**
```json
{
  "name": "Daily Backup",
  "description": "Backup database daily",
  "image": "postgres:15",
  "jobType": "cron",
  "cronExpr": "0 3 * * *",
  "command": "pg_dump -U postgres mydb > /backup/db.sql",
  "timezone": "UTC"
}
```

#### GET /api/schedulers/:id
Get scheduler details.

#### PUT /api/schedulers/:id
Update scheduler configuration.

#### DELETE /api/schedulers/:id
Delete scheduler.

#### POST /api/schedulers/:id/run
Run scheduler immediately.

#### GET /api/schedulers/:id/history
Get scheduler execution history.

#### GET /api/dashboard/stats
Get dashboard statistics.

**Response:**
```json
{
  "total": 10,
  "active": 7,
  "paused": 2,
  "inactive": 1
}
```

### User Management Endpoints

#### GET /api/users
Get all users with role information.

#### POST /api/users
Create new user (admin only).

**Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "roleId": "507f1f77bcf86cd799439012",
  "password": "initialpass123"
}
```

#### PUT /api/users/:id
Update user information.

#### DELETE /api/users/:id
Delete user.

#### POST /api/users/:id/reset-password
Reset user password (admin only).

### Role Management Endpoints

#### GET /api/roles
Get all roles with user counts.

#### POST /api/roles
Create new role.

#### PUT /api/roles/:id
Update role.

#### DELETE /api/roles/:id
Delete role.

#### GET /api/roles/permissions
Get available permissions.

### Image Management Endpoints

#### GET /api/images
Get all Docker images with usage statistics.

#### POST /api/images
Add new Docker image.

#### PUT /api/images/:id
Update image information.

#### DELETE /api/images/:id
Delete image.

### System Logs Endpoints

#### GET /api/logs
Get system logs with filtering.

**Query Parameters:**
- `level`: Filter by log level (info, warning, error)
- `source`: Filter by log source
- `search`: Search in log messages
- `limit`: Number of logs to return
- `offset`: Pagination offset

#### GET /api/logs/stats
Get log statistics.

#### GET /api/logs/sources
Get available log sources.

### Monitoring Endpoints

#### GET /api/monitoring
Get complete monitoring data.

#### GET /api/monitoring/metrics
Get system metrics (CPU, memory, disk, network).

#### GET /api/monitoring/services
Get service health status.

#### GET /api/monitoring/alerts
Get system alerts.

---

## Deployment Guide

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

### Quick Start with Docker Compose

1. **Clone the repository:**
```bash
git clone <repository-url>
cd smart-task-orchestrator
```

2. **Start all services:**
```bash
docker-compose up --build
```

3. **Access the application:**
- Frontend: http://localhost:3000
- API: http://localhost:8080
- Health Check: http://localhost:8080/health

4. **Default credentials:**
- Username: `admin`
- Password: `admin` (must change on first login)

### Local Development Setup

#### Backend Development
```bash
cd backend

# Install dependencies
go mod tidy

# Start API server
go run cmd/api/main.go

# Start scheduler (in another terminal)
go run cmd/scheduler/main.go

# Start worker (in another terminal)
go run cmd/worker/main.go
```

#### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### Environment Configuration

Create `.env` file in project root:

```env
# MongoDB
MONGO_URI=mongodb://localhost:27017
DB_NAME=orchestrator

# Redis
REDIS_URL=redis://localhost:6379

# Kafka
KAFKA_BROKER=localhost:9092

# JWT
JWT_SECRET=your-secret-key-change-in-production

# API
PORT=8080
```

### Production Deployment

1. **Update environment variables** in `docker-compose.yml`
2. **Set strong JWT secret**
3. **Configure MongoDB and Redis persistence**
4. **Set up reverse proxy (nginx) for SSL termination**
5. **Configure monitoring and alerting**
6. **Set up backup and disaster recovery**

### Health Checks

All services expose health check endpoints:
- API: `GET /health`
- Scheduler: Internal health monitoring
- Worker: Internal health monitoring

### Monitoring

View service logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f scheduler
docker-compose logs -f worker
docker-compose logs -f frontend
```

### Troubleshooting

#### Services won't start
```bash
# Check service status
docker-compose ps

# Restart specific service
docker-compose restart api

# Rebuild and restart
docker-compose up --build
```

#### Port conflicts
```bash
# Check what's using the ports
lsof -ti:8080 | xargs kill -9  # Kill API
lsof -ti:3000 | xargs kill -9  # Kill Frontend
```

#### Database issues
```bash
# Access MongoDB shell
docker-compose exec mongodb mongosh

# Check database
use orchestrator
db.users.find()
```

---

## Testing Guide

### Manual API Testing

#### Test Authentication
```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}'

# Save the token from response
TOKEN="eyJhbGc..."

# Test authenticated endpoint
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN"
```

#### Test Scheduler Creation
```bash
curl -X POST http://localhost:8080/api/schedulers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Job",
    "description": "Test scheduler",
    "image": "ubuntu:22.04",
    "jobType": "immediate",
    "command": "echo Hello World"
  }'
```

#### Test User Creation
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "roleId": "507f1f77bcf86cd799439012",
    "password": "testpass123"
  }'
```

### Frontend Testing

1. **Login Flow**
   - Navigate to http://localhost:3000
   - Login with admin/admin
   - Verify password change prompt
   - Change password
   - Verify redirect to dashboard

2. **Dashboard**
   - Verify statistics display
   - Check scheduler list
   - Test navigation to other pages

3. **Scheduler Management**
   - Create new scheduler
   - Edit existing scheduler
   - Run scheduler manually
   - View execution history
   - Delete scheduler

4. **User Management**
   - Create new user
   - Assign role
   - Reset password
   - Delete user

5. **System Monitoring**
   - View system logs
   - Check monitoring metrics
   - View service health

### Performance Testing

#### Load Testing with Apache Bench
```bash
# Test API endpoint
ab -n 1000 -c 10 -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/schedulers
```

#### Concurrent Job Execution
```bash
# Create multiple jobs simultaneously
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/schedulers \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"name\": \"Job $i\", \"jobType\": \"immediate\", \"image\": \"ubuntu:22.04\", \"command\": \"echo Job $i\"}" &
done
wait
```

### Integration Testing

#### End-to-End Job Flow
1. Create scheduler via API
2. Verify scheduler appears in database
3. Wait for scheduled execution
4. Verify job runs in Docker
5. Check execution history
6. Verify logs are stored

---

## Performance Targets

- **Scheduling Precision**: 1-5 seconds maximum delay
- **Concurrent Jobs**: 500+ simultaneous executions
- **Log Latency**: <100ms from container to UI
- **API Response**: <200ms for most operations
- **Database Queries**: <50ms for indexed queries
- **WebSocket Latency**: <10ms for log streaming

---

## Security Considerations

### Authentication
- JWT tokens with 1-hour expiration
- Refresh tokens with 7-day expiration
- Secure password hashing with bcrypt
- Forced password change on initial login

### Authorization
- Role-based access control (RBAC)
- Fine-grained permissions per scheduler
- Admin-only operations protected
- Audit logging for security events

### Data Protection
- Passwords never stored in plain text
- Sensitive data encrypted at rest
- HTTPS required in production
- CORS properly configured

### Container Security
- Docker socket access controlled
- Container resource limits enforced
- Network isolation for containers
- Image scanning recommended

---

## Maintenance & Operations

### Backup Strategy
- MongoDB: Daily automated backups
- Redis: Persistence enabled with AOF
- Configuration: Version controlled
- Logs: Rotated and archived

### Monitoring
- Service health checks
- Resource usage monitoring
- Error rate tracking
- Performance metrics

### Scaling
- Horizontal: Add more worker instances
- Vertical: Increase container resources
- Database: MongoDB replica sets
- Cache: Redis cluster mode

---

## Future Enhancements

### Planned Features
- [ ] Job dependencies and workflows
- [ ] Webhook notifications
- [ ] Advanced scheduling with timezone support
- [ ] Job cancellation
- [ ] Resource quotas and limits
- [ ] Multi-tenant support
- [ ] Job templates
- [ ] Integration with external systems

### Performance Improvements
- [ ] Query optimization
- [ ] Caching strategies
- [ ] Connection pooling
- [ ] Batch operations

### Security Enhancements
- [ ] Two-factor authentication
- [ ] API rate limiting
- [ ] IP whitelisting
- [ ] Enhanced audit logging

---

## Support & Troubleshooting

### Common Issues

1. **Cannot connect to MongoDB**
   - Check MongoDB is running: `docker-compose ps`
   - Verify connection string in environment variables
   - Check network connectivity

2. **Jobs not executing**
   - Verify worker service is running
   - Check Kafka connectivity
   - Review worker logs for errors

3. **Frontend not loading**
   - Check API is accessible
   - Verify CORS configuration
   - Check browser console for errors

4. **Authentication failures**
   - Verify JWT secret is configured
   - Check token expiration
   - Review auth middleware logs

### Getting Help

- Check logs: `docker-compose logs -f <service>`
- Review documentation in this file
- Check GitHub issues
- Contact support team

---

## License

MIT License - see LICENSE file for details.

---

## Changelog

### Version 1.0.0 (Current)
- Initial release with complete feature set
- Real API integration
- Production-ready deployment
- Comprehensive documentation

---

**Last Updated**: November 16, 2025
**Version**: 1.0.0
**Status**: Production Ready
