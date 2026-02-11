
# 🧠 CONSUMER EXECUTION FLOW — DOCKER SPIN-UP, HISTORY MANAGEMENT & LOG STREAMING

---

## 1️⃣ Overview

The **consumer** is a background service that listens to job execution messages (from Kafka or Redis queue).
Each message means:

> “Run this job now — using this Docker image, command, and metadata.”

When the consumer receives such a message, it must:

1. Record the **start** of execution in `scheduler_history`.
2. Pull and run the **Docker container** for that job.
3. Stream the **container logs (stdout/stderr)** line-by-line (with ANSI color codes preserved).
4. Send each log line to:

   * The **WebSocket/Log Gateway** (for real-time UI updates).
   * The **DB or blob store** (for persistence after job finishes).
5. When the job ends, mark the **final status** in history.
6. Handle cleanup, retries (if any), and resource release safely.

---

## 2️⃣ Components Involved

| Component                         | Role                                                        |
| --------------------------------- | ----------------------------------------------------------- |
| **Kafka Consumer**                | Reads job execution events from `scheduler_run_topic`.      |
| **Docker Engine API**             | Runs containers on host (not inside another container).     |
| **MongoDB (`scheduler_history`)** | Persists job status and metadata.                           |
| **Redis Stream**                  | Used for log streaming buffer (before WebSocket broadcast). |
| **WebSocket Gateway**             | Pushes live logs from Redis to frontend.                    |
| **Blob Store (S3/minio)**         | Stores large completed logs (optional).                     |

---

## 3️⃣ Job Message Structure

Example Kafka message published by scheduler library:

```json
{
  "run_id": "run-789",
  "scheduler_id": "sched-123",
  "generation": 4,
  "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/data-pipeline:v5",
  "command": "bash /opt/run_etl.sh",
  "env": { "REGION": "us-east-1" },
  "triggered_by": "user-999",
  "run_at": "2025-10-22T10:00:00Z"
}
```

---

## 4️⃣ Step-by-Step Execution Flow (Explained)

### Step 1: Message Consumption

* The consumer continuously polls Kafka for new messages.
* Each message represents one job execution.
* When a message arrives:

  1. Validate required fields.
  2. Check `generation` against current `scheduler_definition.generation` in DB.

     * If mismatch → mark as **discarded** in history (don’t execute).
  3. Insert a new entry in `scheduler_history`:

     ```text
     run_id: "run-789"
     scheduler_id: "sched-123"
     start_time: now()
     status: "running"
     executed_by: "user-999"
     generation: 4
     ```
  4. Proceed to container spin-up.

---

### Step 2: Docker Container Spin-Up

**Key concept:**
The consumer does **not** run inside Docker-in-Docker.
It connects to the host’s Docker daemon using the **Docker SDK** (or CLI interface) — same as `docker run`.

**Behavior:**

1. Pull image (if not cached):

   * Authenticate with AWS ECR using IAM role or credentials.
   * Command equivalent:
     `docker pull 123456789012.dkr.ecr.us-east-1.amazonaws.com/data-pipeline:v5`

2. Create container:

   * Assign name = `sched-<id>-run-<run_id>`.
   * Mount ephemeral volume for logs (optional).
   * Set environment vars (if any).
   * Command = the multi-line script (executed via `/bin/sh -c "<command>"`).
   * Attach both stdout and stderr to log stream.

3. Start container.

   * Record container ID in `scheduler_history.process_id`.

4. Immediately start reading the container logs stream (next step).

---

### Step 3: Capture Logs in Real-Time

When the container starts, the consumer **attaches** to its stdout and stderr streams.

Each log line:

* Contains raw **ANSI escape sequences** (for color, bold, etc.).
* Must be preserved exactly (so the UI can render true terminal formatting).

The consumer will:

1. Read stream line by line (non-blocking).

2. For each line:

   * Publish it to Redis Stream `logs:<run_id>`, e.g.

     ```json
     { "type": "log", "line": "\u001b[32mSuccess!\u001b[0m" }
     ```
   * Also append to an in-memory buffer (for final upload or DB storage).
   * Optionally patch a few lines into the `scheduler_history.log_text` field (truncated preview).

3. The **WebSocket Gateway** subscribes to these Redis streams.
   Whenever a new entry appears in `logs:<run_id>`, it forwards it to all connected clients in real-time.

UI result: the terminal on the screen updates live — preserving colors, spaces, and structure.

---

### Step 4: Job Completion & Finalization

When the container stops:

1. Capture the **exit code** (e.g., 0 = success, >0 = failure).

   * Equivalent to `docker inspect <container_id> --format '{{.State.ExitCode}}'`.

2. Mark final status in `scheduler_history`:

   * If exit code 0 → `status = "success"`.
   * Else → `status = "failed"`.
   * Add `end_time = now()`, `exit_code`, `error_message` if any.

3. Upload full logs:

   * Combine buffered lines into a text blob.
   * Upload to S3 or Minio bucket (`logs/<run_id>.log`).
   * Save `log_blob_id` = URL/path in history.

4. Send closing log message to Redis Stream:

   ```json
   { "type": "end", "status": "success", "timestamp": "2025-10-22T10:05:13Z" }
   ```

   The frontend sees this → closes the log stream.

5. Cleanup container:

   * `docker rm <container_id>` to free disk.
   * Optionally delete pulled image (if ephemeral system).

---

### Step 5: Error Handling & Recovery (Human Explanation)

| Failure Type                         | Handling Strategy                                                                                                                  |
| ------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------- |
| **Container pull fails (ECR issue)** | Log error, set `status="failed"`, `error_message="Image pull failed"`.                                                             |
| **Command fails (exit code > 0)**    | Mark as `failed`, but **do not retry** (as per your decision).                                                                     |
| **Container runtime error**          | Capture exception, cleanup, mark `failed`.                                                                                         |
| **Redis log stream unavailable**     | Keep local buffer and push later; UI may temporarily miss live logs.                                                               |
| **Worker crash mid-run**             | On restart, detect unfinished runs (status=running but container exited) → finalize status as `failed` with reason `worker_crash`. |

---

### Step 6: How Logs Are Streamed to Frontend

Flow summary:

1. Consumer → Redis Stream
2. WebSocket Gateway → UI

**Redis Stream (fast, ordered):**
Each run has its own stream key: `logs:run-789`.
Each line added as entry with incremental ID.

**WebSocket Gateway:**

* Subscribes to stream using `XREAD BLOCK 0`.
* Forwards new log lines to frontend over WebSocket:

  ```json
  { "type": "log", "line": "\u001b[33mProcessing file 1...\u001b[0m" }
  ```
* When final message (`type=end`) arrives → closes stream.

**Frontend (React):**

* Opens WebSocket connection on job start.
* Renders logs inside `<pre>` block.
* Uses client-side `ansi-to-html` to colorize lines.
* Auto-scrolls with each new log.

---

### Step 7: History Maintenance & Data Retention

* Every job run adds one entry in `scheduler_history`.
* Daily background cleanup or archival process:

  * If logs older than N days (e.g., 30) → delete from DB, keep only S3 file.
  * Compress large logs before upload.
* Maintain indexes on `scheduler_id` + `start_time` for quick filtering.

Retention example:

| Data Type                    | Retain For | Storage |
| ---------------------------- | ---------- | ------- |
| Recent history (status only) | 30 days    | MongoDB |
| Logs                         | 90 days    | S3      |
| Audit events                 | 1 year     | MongoDB |
| Metadata                     | permanent  | MongoDB |

---

### Step 8: Concurrency, Parallelism & Scaling

* Each consumer instance can run multiple containers concurrently (configurable, e.g., `MAX_CONCURRENT_JOBS=5`).
* Maintain an internal queue with worker slots.
* To scale horizontally:

  * Deploy multiple consumer pods/nodes.
  * Partition Kafka topic by scheduler_id (so same job always handled by same partition → no duplicates).
* For 500 active jobs (your current scale), one or two consumers are sufficient.

---

### Step 9: Example Timeline

Let’s walk through a real case:

| Time              | Action                                                       | Actor              |
| ----------------- | ------------------------------------------------------------ | ------------------ |
| 10:00:00          | Scheduler library publishes job run message                  | scheduler          |
| 10:00:01          | Consumer receives message, records history(status=running)   | consumer           |
| 10:00:02          | Docker container starts, logs begin                          | consumer           |
| 10:00:03–10:05:10 | Consumer streams logs to Redis → WebSocket → UI              | consumer & gateway |
| 10:05:11          | Container exits (code=0)                                     | Docker             |
| 10:05:12          | Consumer uploads logs to S3, updates history(status=success) | consumer           |
| 10:05:13          | WebSocket stream ends (UI sees “Job completed successfully”) | UI                 |

---

### Step 10: Log Format and UI Rendering Details

* Preserve **ANSI escape codes** (e.g., `\u001b[31m` for red).
* The frontend terminal uses `<pre>` and converts ANSI → HTML:

  * Green text for success
  * Red for errors
  * Yellow for warnings

Example log segment as seen by user:

```
[INFO] Initializing job...
[INFO] Loading input files
[ERROR] Missing configuration file
[INFO] Retrying step...
[OK] Job completed successfully
```

Internally, these color codes remain:

```
\u001b[32m[OK] Job completed successfully\u001b[0m
```

so the browser renders colors accurately.

---

### Step 11: Terminal Patching (UI Log Live Sync)

Your goal was to “patch the terminal in logs and poll the logs stream into scheduler_history — preserving colors, spaces, formatting.”

Implementation conceptually:

1. **During run:** consumer pushes logs → Redis → WebSocket → UI.
2. **Every few seconds**, the consumer also batches the latest few KB of logs and **updates the DB field** `scheduler_history.log_text_partial`.

   * This ensures backend UI can reload partial logs if a user refreshes mid-run.
3. **After completion:** replace `log_text_partial` with the final `log_blob_id` (pointer to S3).
4. **Formatting preserved**: store as UTF-8 plain text with ANSI sequences intact.

   * Do **not strip** color codes at backend.
   * Let frontend handle rendering conversion.

---

### Step 12: Summary of Consumer Responsibilities

| Step | Action                                | Data Store    |
| ---- | ------------------------------------- | ------------- |
| 1    | Consume job message                   | Kafka         |
| 2    | Validate job (generation, active)     | MongoDB       |
| 3    | Create history entry (status=running) | MongoDB       |
| 4    | Pull & start Docker container         | Docker Engine |
| 5    | Stream stdout/stderr → Redis logs     | Redis         |
| 6    | WebSocket gateway streams logs to UI  | WebSocket     |
| 7    | On container exit, update status      | MongoDB       |
| 8    | Upload logs to S3, update history     | MongoDB       |
| 9    | Cleanup container                     | Docker        |
| 10   | Publish metrics (optional)            | Prometheus    |

---

### Step 13: Metrics & Observability (for production readiness)

Expose metrics like:

| Metric                           | Meaning                                    |
| -------------------------------- | ------------------------------------------ |
| `consumer_jobs_running`          | number of containers currently executing   |
| `consumer_jobs_success_total`    | total successful jobs                      |
| `consumer_jobs_failed_total`     | total failed jobs                          |
| `job_execution_duration_seconds` | histogram of run durations                 |
| `log_stream_latency_ms`          | time between log line and frontend display |
| `docker_pull_time_seconds`       | how long image pull took                   |

Send these to Prometheus / Grafana dashboards.

---

### Step 14: Example Failure Scenarios (and Recovery Behavior)

| Scenario                  | What Happens                                       | Result                             |
| ------------------------- | -------------------------------------------------- | ---------------------------------- |
| ECR image missing         | Container pull fails                               | History marked failed (pull error) |
| Command script typo       | Exit code >0                                       | Status=failed, log captured        |
| Redis down during logs    | Logs temporarily buffered, re-sent later           | Slight delay in UI                 |
| Worker crashes mid-job    | Container exits, job found incomplete on restart   | History finalized as failed        |
| UI disconnected           | Logs still flow to Redis & DB; UI can reload later | No data loss                       |
| Scheduler updated mid-run | Generation mismatch prevents future old runs       | Current one still executes         |

---

### Step 15: Extensibility & Future Enhancements

| Feature         | Description                                                            |
| --------------- | ---------------------------------------------------------------------- |
| Retry Policy    | In future, if container spawn fails (not code failure), retry N times. |
| Resource Limits | CPU/memory quotas for each job using Docker run options.               |
| Log Compression | Gzip long logs before S3 upload.                                       |
| Job Timeout     | Allow max runtime; auto-stop container if exceeded.                    |
| Parallel Logs   | Merge multiple containers (for distributed jobs) into one stream.      |
| Execution Tags  | Attach “env=prod” or “pipeline=batch” tags for filtering in UI.        |

---

### ✅ Final Summary

In plain terms:

* The **consumer** is your executor and reporter.
* It **runs containers**, **records everything**, and **streams logs** faithfully to both your backend and UI.
* It ensures **no loss of formatting**, **accurate timing**, **consistent history**, and **graceful recovery** even under partial failure.

Think of it as your “worker with eyes and mouth” —
it does the job *and* tells everyone exactly what’s happening, line by line, in real time.

---

Would you like me to now extend this section with **the WebSocket Gateway design** (how it subscribes to Redis logs, authenticates users, and forwards messages to specific clients securely)?
That’s the next logical layer between this consumer and your frontend.
