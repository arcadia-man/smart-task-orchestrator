# Smart Task Orchestrator - Complete Rebuild Plan

## 🎯 Overview
Complete rebuild of the Smart Task Orchestrator with production-ready architecture, advanced scheduling, RBAC, and real-time monitoring.

## 🏗️ New Architecture

### Core Services
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

### Service Breakdown

1. **API Server** (`cmd/api`)
   - REST endpoints for CRUD operations
   - JWT authentication with RBAC
   - WebSocket for real-time log streaming
   - Admin user seeding and management

2. **Scheduler Library** (`cmd/scheduler`) 
   - Precompute service (every 5 minutes)
   - Fast poller (every 1 second)
   - Redis sharding for scalability
   - Generation-based invalidation

3. **Worker Service** (`cmd/worker`)
   - Kafka consumer for job execution
   - Docker container management
   - Real-time log streaming with ANSI preservation
   - Process lifecycle management

4. **Frontend** (`frontend/`)
   - React dashboard with real-time updates
   - Job creation with advanced scheduling
   - Live log viewer with color preservation
   - User and permission management

## 📊 Enhanced Data Models

### MongoDB Collections

```javascript
// scheduler_definition
{
  _id: "uuid",
  name: "string",
  description: "string", 
  image: "string",           // ECR path
  job_type: "cron|interval|immediate",
  cron_expr: "string",       // for cron jobs
  interval_seconds: "number", // for interval jobs
  command: "string",         // multi-line script
  timezone: "string",        // default UTC
  status: "active|inactive",
  generation: "number",      // for invalidation
  created_at: "datetime",
  updated_at: "datetime",
  created_by: "uuid",
  updated_by: "uuid"
}

// scheduler_precompute (durable queue)
{
  _id: "uuid",
  scheduler_id: "uuid",
  run_at: "datetime",
  generation: "number",
  status: "pending|dispatched|canceled|discarded",
  created_at: "datetime"
}

// scheduler_history (execution logs)
{
  _id: "uuid",              // run_id
  scheduler_id: "uuid",
  precompute_id: "uuid",
  executed_by: "uuid",      // user or system
  status: "pending|running|success|failed|discarded",
  start_time: "datetime",
  end_time: "datetime",
  command: "string",        // actual command executed
  exit_code: "number",
  process_id: "string",     // container ID
  log_text: "string",       // full logs
  error_message: "string",
  created_at: "datetime",
  updated_at: "datetime"
}

// Enhanced RBAC
// users, roles, permissions, images as per requirements
// + audit_logs for security tracking
```

### Redis Data Structures

```redis
# Sharded sorted sets for scheduled items
sched:zset:0 -> {member: "precompute:uuid", score: run_at_ms}
sched:zset:1 -> {member: "precompute:uuid", score: run_at_ms}

# Optional metadata cache
sched:meta:uuid -> {scheduler_id, run_at_ms, generation}

# Per-scheduler invalidation index
sched:ids:scheduler_uuid -> [precompute_id1, precompute_id2, ...]

# Real-time log streaming
logs:run_uuid -> Redis Stream for WebSocket
```

## 🔧 Implementation Phases

### Phase 1: Core Infrastructure
- [ ] Enhanced data models with proper indexing
- [ ] JWT authentication with RBAC middleware
- [ ] Redis sharding setup
- [ ] Kafka topic configuration
- [ ] Docker multi-stage builds

### Phase 2: Scheduler Library
- [ ] Precompute service with 15-minute lookahead
- [ ] Fast poller with 1-second precision
- [ ] Generation-based invalidation logic
- [ ] Distributed locking for coordination
- [ ] Metrics and monitoring hooks

### Phase 3: Worker & Execution
- [ ] Kafka consumer with proper error handling
- [ ] Docker container lifecycle management
- [ ] Real-time log streaming with ANSI preservation
- [ ] Process monitoring and cleanup
- [ ] Resource limit enforcement

### Phase 4: Frontend Enhancement
- [ ] Advanced job creation form
- [ ] Real-time dashboard with WebSocket
- [ ] Live log viewer with color support
- [ ] User and permission management UI
- [ ] Job history and analytics

### Phase 5: Production Features
- [ ] Comprehensive monitoring and alerting
- [ ] Backup and disaster recovery
- [ ] Performance optimization
- [ ] Security hardening
- [ ] Documentation and deployment guides

## 🚀 Key Improvements Over Current System

1. **Precision Scheduling**: 1-second accuracy vs current minute-level
2. **Scalability**: Redis sharding + horizontal scaling
3. **Security**: Full RBAC with audit logging
4. **Reliability**: Generation-based invalidation prevents stale jobs
5. **Observability**: Real-time metrics and comprehensive logging
6. **User Experience**: Live log streaming with color preservation
7. **Production Ready**: Proper error handling, monitoring, and deployment

## 📋 Acceptance Criteria

- [ ] Jobs execute within 1-5 seconds of scheduled time
- [ ] System handles 1000+ concurrent jobs
- [ ] Zero duplicate job executions
- [ ] Real-time log streaming with ANSI colors
- [ ] Complete RBAC with fine-grained permissions
- [ ] Graceful handling of configuration changes
- [ ] Comprehensive monitoring and alerting
- [ ] Production-ready Docker deployment

## 🔄 Migration Strategy

1. **Parallel Deployment**: Run new system alongside current
2. **Data Migration**: Export/import job definitions
3. **Gradual Cutover**: Move jobs incrementally
4. **Validation**: Compare execution results
5. **Full Cutover**: Switch DNS/load balancer
6. **Cleanup**: Remove old system after validation

This rebuild will transform the current basic job orchestrator into a production-ready, enterprise-grade distributed scheduler capable of handling thousands of jobs with sub-second precision.