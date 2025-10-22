# Smart Task Orchestrator - Requirements Summary

## 🎯 Core Requirements

### Functional Requirements
1. **Job Scheduling**: Immediate, cron-based, and interval-based job execution
2. **Distributed Execution**: Multi-worker Docker container execution
3. **Real-time Monitoring**: Live log streaming with ANSI color preservation
4. **RBAC**: Role-based access control with fine-grained permissions
5. **Exact Timing**: 1-5 second precision for job execution

### Technical Requirements
1. **Backend**: Go with Gin framework
2. **Frontend**: React with TypeScript and Tailwind CSS
3. **Database**: MongoDB for flexible document storage
4. **Cache**: Redis for fast queuing and coordination
5. **Messaging**: Kafka for reliable event streaming
6. **Containers**: Docker for job execution

## 🏗️ Architecture Components

### Backend Services
- **API Server**: REST endpoints, authentication, WebSocket
- **Scheduler**: Precompute (5min) + Poller (1sec) for exact timing
- **Worker**: Docker container execution with log streaming

### Data Models
- **Users & Roles**: RBAC with auto-role creation
- **Scheduler Definitions**: Job configurations with generation tracking
- **Execution History**: Complete audit trail with logs
- **Permissions**: Fine-grained access control per scheduler

### Key Features
- **Generation-based Invalidation**: Prevents stale job execution
- **Real-time Log Streaming**: WebSocket + Redis streams
- **Docker Integration**: Host Docker execution (not Docker-in-Docker)
- **Horizontal Scaling**: Redis sharding + Kafka partitioning

## 🔐 Security & Authentication

### Initial Setup
- Default admin user: `admin/admin` (must change on first login)
- JWT-based authentication with refresh tokens
- Role auto-creation when typed in UI

### RBAC Structure
- **Users**: Basic user information and role assignment
- **Roles**: Named roles with coarse permissions
- **Permissions**: Fine-grained per-scheduler access control

## 🚀 Deployment Strategy

### Docker Composition
- **Backend Image**: Single image with multiple binaries (api, scheduler, worker)
- **Frontend Image**: Nginx-served React build
- **Infrastructure**: MongoDB, Redis, Kafka via official images

### Production Considerations
- Environment-based configuration
- Persistent volumes for data
- Health checks and monitoring
- Horizontal scaling support

## 📊 Performance Targets

- **Scheduling Precision**: 1-5 seconds maximum delay
- **Concurrent Jobs**: 500+ simultaneous executions
- **Log Latency**: <100ms from container to UI
- **API Response**: <200ms for most operations

## 🔄 Consumer Execution Flow

### Job Lifecycle
1. **Message Consumption**: Kafka job execution events
2. **Validation**: Generation check and scheduler status
3. **Container Execution**: Docker pull, run, and monitor
4. **Log Streaming**: Real-time stdout/stderr to Redis + WebSocket
5. **Completion**: Status update, log persistence, cleanup

### Error Handling
- Container failures: Mark as failed, preserve logs
- Network issues: Retry with exponential backoff
- Worker crashes: Detect and finalize incomplete jobs

This document serves as the single source of truth for the Smart Task Orche