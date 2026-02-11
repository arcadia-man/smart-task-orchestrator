# Smart Task Orchestrator — Implementation & Developer Prompt

> This document combines the original requirements (`new requirement.md`), the stored decisions you approved, and a concrete, developer-ready implementation prompt for building the scheduler library and surrounding system during a "vibe coding" session today. Use this as the single-source-of-truth for coding, tests, and acceptance criteria.

---

## Contents

1. Original requirements (embedded)
2. Decisions & stored points (what we saved)
3. System architecture (components)
4. Data models (DB + Redis keys)
5. Precompute & Poller design (every 5m / 15m lookahead variant)
6. Scheduler library scope + public API (what to implement today)
7. Worker / Container execution & logging
8. Invalidation, de-dup & concurrency controls
9. Operational guidelines & scaling
10. Observability, testing & acceptance criteria
11. Vibe-coding plan for today + final prompt for the developer (you)

---

# 1. Embedded original requirement.md

(Full content from `new requirement.md` as provided by the user)

```
# 📘 Smart Task Orchestrator: Distributed Job Scheduler

## 🔧 Objective

Design and implement a **distributed job scheduling system** where tasks can be submitted, managed, executed, and monitored at scale. It should support real-time and scheduled executions across distributed worker nodes using a modern tech stack.

---

## 📌 Functional Requirements

1. **Job Scheduling**

   * Users can submit jobs:

     * For **immediate execution**
     * At **specific times** (cron-style)
     * At **intervals**
   * Jobs include executable commands inside Docker containers.

2. **Distributed Execution**

   * Tasks are dispatched to multiple worker nodes in a distributed setup.
   * Workers pull jobs and execute them concurrently.

3. **Monitoring and Reporting**

   * Real-time job execution tracking.
   * System health monitoring.
   * Resource utilization reporting (optional advanced feature).

---

## 📌 Non-Functional Requirements

1. **Reliability**

   * Jobs must execute **accurately and on time**.
   * System should have high availability.

2. **Performance**

   * Minimal latency and overhead in job dispatch and execution.

3. **Scalability**

   * Horizontally scalable system components:

     * Backend
     * Worker nodes
     * Kafka/Redis

4. **Fault Tolerance**

   * Retry mechanisms on job failure.
   * Node failure detection and recovery.

5. **Security**

   * Authenticated access.
   * Role-based access control (RBAC).
   * Secure job execution outside orchestrator.

---

## 💻 Tech Stack

* **Backend**: Go
* **Frontend**: ReactJS, TailwindCSS
* **Databases**: MongoDB
* **Messaging Queue**: Kafka
* **Cache & Coordination**: Redis
* **Container Runtime**: Docker

---

## 🧩 System Design Goals

* **Database Initialization**: On startup, the application should:

  * Perform DB migrations (collections, indexes).
  * Verify and create necessary Kafka topics.
  * Setup Redis keys or streams.

* **YAML Template / Docker Compose File**

  * Users must be able to:

    * Pull from GitHub
    * Edit config with credentials (MongoDB, Redis, Kafka)
    * Run the system anywhere
  * This should be **clearly documented in README.md**

* **UI Scalability vs Backend**

  * Frontend is **not required to scale much**.
  * Backend services must be **decoupled and scalable**.
  * Stateless Go services with horizontal scaling.

---

## 🔐 Authentication & Authorization

### Authentication Rules

1. **Initial Admin User**

   * On first launch, one default **admin user** should be seeded with:

     * `username`: "admin"
     * `password`: "admin"

2. **User Management (by Admin only)**

   * Admin can create users.
   * Roles should be **auto-populated** from the `roles` collection.
   * If a new role is typed, it should be **created automatically**.

3. **Initial Login**

   * All users must **change their password** on first login (admin too).

### Database Schema

#### 1. `users`

| Field          | Type     | Notes              |
| -------------- | -------- | ------------------ |
| userid         | UUID     | Primary Key        |
| username       | string   | Unique             |
| password       | string   | Encrypted          |
| email          | string   |                    |
| isInitialLogin | boolean  | Default `true`     |
| role_id        | FK       | References `roles` |
| active         | boolean  |                    |
| created_at     | datetime |                    |
| updated_at     | datetime |                    |

#### 2. `roles`

| Field           | Type     | Notes                |
| --------------- | -------- | -------------------- |
| id              | UUID     | Primary Key          |
| role_name       | string   | Unique               |
| active          | boolean  | Default `true`       |
| can_create_task | boolean  | Role permission flag |
| created_at      | datetime |                      |

---

## 🗓 Scheduler (Tasks) Requirements

1. **Task Creation (UI Form)**

   * Form with multiple sections:

     * Scheduler Definition
     * Permissions

2. **Execution Logic**

   * Use provided Docker image (hosted on AWS ECR).
   * Spin a new Docker container **on host system** (not Docker-in-Docker).
   * Run the specified shell command inside the container.

3. **Logging**

   * Stream terminal logs from running container.
   * Maintain:

     * Color codes
     * Line formatting
     * Stream logs to history table (poll or stream)

4. **Scheduling Logic**

   * Central scheduler checks every second for jobs ready to run.
   * Tolerate **1-5 second delay at scale**.

---

## 🧾 Scheduler Schema

### 1. `scheduler_definition`

| Field            | Type     | Notes                           |
| ---------------- | -------- | ------------------------------- |
| id               | UUID     | Primary Key                     |
| description      | string   |                                 |
| image            | string   | Docker image name or ECR path   |
| jobType          | enum     | `immediate`, `cron`, `interval` |
| command          | string   | Shell command                   |
| status           | string   | `active`, `inactive`, etc.      |
| created_at       | datetime |                                 |
| updated_at       | datetime |                                 |
| last_modified_by | string   | User who last modified          |

---

### 2. `scheduler_history`

| Field            | Type     | Notes                                        |
| ---------------- | -------- | -------------------------------------------- |
| id               | UUID     |                                              |
| scheduler_id     | FK       | Links to `scheduler_definition`              |
| current_status   | boolean  | Paused, running, failed, pending.            |
| log              | text     | Full log output                              |
| process_id       | string   | PID of spawned container                     |
| execution_time   | datetime | When it actually ran                         |
| command_executed | string   | Final command run (may differ from original) |
| executed_by      | string   | Username                                     |
| status           | string   | `running`, `success`, `failed`               |
| start_time       | datetime |                                              |
| end_time         | datetime |                                              |

---

### 3. `permissions`

| Field           | Type    | Notes                                 |
| --------------- | ------- | ------------------------------------- |
| id              | UUID    |                                       |
| role_id         | FK      | Role from `roles` table               |
| scheduler_id    | FK      | Scheduler from `scheduler_definition` |
| roleSee         | boolean | View permission                       |
| roleExecute     | boolean | Execution permission                  |
| roleAlterConfig | boolean | Update permission                     |

---

### 4. `images`

| Field        | Type   | Notes            |
| ------------ | ------ | ---------------- |
| id           | UUID   |                  |
| image        | string | Docker image URI |
| name         | string | Logical name     |
| description  | string | Optional         |
| version      | string | E.g., `v1.2.3`   |
| registry_url | string | AWS ECR or other |

---

I am doing vibe coding I need the command to give the cursor. Store the above data and answer or let’s discuss below questions that can be added to in requirement. When I will say I am fully prepared according to our discussion then you give me the actual command. I between I will give feed back please make sure that you mention in final prompt need to give the cursor.

— chat gpt please here help me how it should be docker image for backend and frontend I don’t know different devops how they will handle. Making two image just is fine or not 
— chat gpt please check this schema if you find any wrong or missing then tell me

— chat gpt suggest additional functional which I have missed and similar kind of application contains

— chat gpt how go will maintain the exact time operation give some idea what similar application are using to achieve this. Give one ideal to achieve very close to it by mentioned tech stack.
```

---

# 2. Decisions & stored points (what we saved)

> These are the items you explicitly accepted and asked to be stored in memory. They must be included in the final implementation and the developer prompt.

* You accepted the suggested schema updates and confirmed the system will include collections/tables: **users, roles, scheduler_history, permissions, images**.
* Workflow decisions that were stored:

  1. Users create tasks via a form with grouped sections: **Scheduler definition** and **Permission**.
  2. Based on configuration, tasks should start execution automatically.
  3. Docker images will be hosted in **AWS ECR**; system will spin up Docker containers on the host machine (not Docker-in-Docker). For dev, it should run locally.
  4. Terminal output must be patched into logs and polled/streamed into `scheduler_history`; preserve ANSI color codes, line spacing and formatting.
  5. The main scheduling library should only handle the scheduling algorithm (make tasks ready at the correct time). A **1–5 second delay** is acceptable at very high scale.
  6. Job dependencies and parameters will be handled in the job's command script (multiline script). Retry policy is postponed for now — future retries will only be for container spawn failures, not code-level failures. Alerting, multi-tenancy, and audit logging are desired. UI enhancements accepted.
* You asked that the final prompt include instruction to **give the cursor** when you say “I am fully prepared”.
* You chose the hybrid **precompute + short-poller** approach (details below): precompute every **5 minutes** with a **15 minute lookahead** (improved decision), poller every **1 second** to pop due items and publish to Kafka.
* Duplicate/overlap prevention: producer-side unique insertion (`scheduler_id`, `run_at`, `generation`) and consumer-side discard if an instance is already running (mark as `discarded` but keep audit if needed).

---

# 3. System architecture (components)

High-level services & responsibilities:

1. **API / Admin (Go)**

   * CRUD for `scheduler_definition`, users, roles, permissions, images
   * Auth & RBAC enforcement
   * Hooks to increment `generation` on update and trigger invalidation

2. **Scheduler Library (Go)** — the core library you will implement today

   * Periodic poller (1s) orchestrator position: evaluates precomputed queue entries, acquires locks if needed and publishes execution events to Kafka
   * Precompute integration hook (not necessarily inside the library) to load precomputed items
   * Public API for listing active schedules per shard

3. **Precompute Service (Go)**

   * Runs every **5 minutes**: computes next run times for lookahead (15 minutes)
   * Writes to `scheduler_precompute` (durable) and Redis ZSET(s)
   * Inserts idempotently (unique constraint scheduler_id+run_at+generation)

4. **Poller (Go)**

   * Runs every **1 second** per-shard
   * Atomically pops due items from Redis (Lua script)
   * Validates generation & publishes to Kafka
   * Marks `scheduler_precompute` as dispatched or canceled based on generation

5. **Dispatcher / Worker(s)**

   * Consume Kafka `job_execution` events
   * Pull image from ECR, `docker run -it --rm` on host, run commands
   * Attach TTY to preserve ANSI and stream logs to `scheduler_history` + Redis streams for UI

6. **Logging & Streaming**

   * Redis streams or WebSocket server for real-time log streaming to UI
   * `scheduler_history` stores full log text and metadata

7. **Coordination & Locking**

   * Redis for distributed locking (shard locks/Redlock pattern)
   * Sharding by `hash(scheduler_id) % N` for parallelism

8. **Messaging**

   * Kafka topic `job_executions` (partition key: `scheduler_id`) for ordering & consumer affinity

9. **Monitoring & Alerting**

   * Prometheus metrics on scheduler latency, popped items/s, consumer lag, container failures
   * Alerting to Slack/email on failures and backpressure

---

# 4. Data models (DB + Redis keys)

## MongoDB Collections / SQL Tables

* `scheduler_definition` (id, description, image, jobType, command, status, generation, created_at, updated_at, last_modified_by)
* `scheduler_precompute` (id, scheduler_id, run_at TIMESTAMP, generation INT, status ENUM('pending','dispatched','canceled'), created_at)
* `scheduler_history` (id, scheduler_id, run_id, status, log, start_time, end_time, exit_code, executed_by, process_id)
* `users`, `roles`, `permissions`, `images` as specified in original doc

## Redis keys & patterns

* Sharded ZSETs: `sched:zset:<shard>` — members: `precompute:<preid>` or compact `scheduler:<id>|run_at_ms|gen` ; score: `run_at_ms`
* Optional meta hash: `sched:meta:<precompute_id>` = `{scheduler_id, run_at_ms, generation}`
* Per-scheduler index (optional): `sched:ids:<scheduler_id>` (list of precompute IDs for fast invalidation)
