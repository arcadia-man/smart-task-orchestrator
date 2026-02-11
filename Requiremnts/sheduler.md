

# 🧭 SMART TASK ORCHESTRATOR — FUNCTIONAL IMPLEMENTATION EXPLANATION (Scheduler Library + Frontend)

---

## PART 1 — Scheduler Library (Backend Core Logic)

### 🧩 What this component does

The **scheduler library** is the *brain* of your orchestration system.
It decides **when** a job is ready to run, **publishes** that execution event, and ensures **no duplicate** or **outdated** jobs are ever dispatched.
It does **not** execute jobs — it only signals that a job is due.

---

### 🧱 Where it fits

It sits inside the backend system (written in Go), and interacts with:

1. **MongoDB** — to read scheduler definitions and update their precomputed runs.
2. **Redis** — to maintain lightweight “ready-to-run” queues.
3. **Kafka** — to publish messages to worker services when jobs are ready to execute.

---

### 🔄 High-Level Workflow

1. **Scheduler Definition**

   * A configuration object (like a job blueprint).
   * It defines *what* to run, *how often*, and *using which Docker image*.

2. **Precompute Service**

   * Runs every few minutes (say every 5 minutes).
   * Looks at each scheduler and calculates the next few execution times — usually for the next 15 minutes.
   * These upcoming times are stored in a table called `scheduler_precompute`, and also pushed into Redis as time-sorted entries.

3. **Scheduler Library (this component)**

   * Every second, it checks Redis for jobs whose time has arrived (or passed slightly).
   * When it finds one, it verifies that:

     * The scheduler is still active.
     * The definition has not been modified since it was precomputed.
   * If all is good, it marks the entry as “dispatched” and sends a message to Kafka — telling the worker system to run the job.

4. **Worker System (not part of this library)**

   * Listens to Kafka and runs the job in Docker.
   * Streams logs and status updates back into the system.

---

### 🧠 How the scheduler knows “when” to run

Each job has:

* A **cron expression** (like `"*/5 * * * *"` = every 5 minutes), or
* An **interval** value (like every 10 seconds).

The scheduler library doesn’t evaluate cron expressions in real time.
Instead, the **precompute service** already calculated future timestamps for it.

So the library just needs to:

* Compare the *current time* with the *scores in Redis* (which are stored as timestamps).
* If the score ≤ now, that means the job’s time has arrived.
* The entry is popped out and processed.

Example:
If Redis has

```
scheduler:zset:3 → 
  ("jobA:1698039600000")  // = 2025-10-22T10:00:00Z
```

and now = 10:00:01,
the library sees it’s time to run “jobA”.

---

### 🧰 How it prevents duplicate or outdated jobs

Each scheduler definition has a **generation number** (an integer that increases every time the job configuration is updated).

Example:

* Job A initially has generation `1`.
* You edit it in the UI to change frequency.
* The backend increments generation to `2`.

All precomputed runs created earlier are still tagged with generation `1`.
When the scheduler library reads them and sees generation `1` ≠ latest generation `2`,
it **skips** them because they’re outdated.
This prevents old schedules from running after a job’s configuration changed.

---

### 🧩 How dispatch works (in plain words)

When a job is ready:

1. The library pops it from Redis.
2. It asks MongoDB:

   * “Hey, what’s the current generation for this job?”
   * If generation matches → okay to run.
3. It marks the precomputed entry in the database as *dispatched* (to avoid future repeats).
4. It creates a **run ID** (a UUID) and publishes a JSON message to Kafka, like:

   ```json
   {
     "run_id": "abcd-1234",
     "scheduler_id": "jobA",
     "generation": 2,
     "run_at": "2025-10-22T10:00:00Z"
   }
   ```
5. Kafka sends that message to a worker service, which will then run the Docker container.

---

### 🧩 Error handling (human explanation)

* **If Redis fails** → the library skips that tick and tries again in the next second.
  Jobs are never lost because precompute data still exists in MongoDB.

* **If Kafka publish fails** → the library temporarily saves the message to a small Redis list named `pending_publish`,
  and a retry routine keeps trying every few seconds.

* **If the system restarts** → the library doesn’t lose track because all scheduled data is in the DB or Redis.
  On startup, it resumes where it left off.

---

### ⚙️ Example: Full lifecycle of a single job

Let’s take Job A (runs every 5 minutes).

1. **Scheduler Definition**

   ```
   id: jobA
   image: aws.ecr/my-app:v1
   command: ./run_task.sh
   cron: */5 * * * *
   generation: 1
   ```

2. **Precompute service** calculates:

   ```
   next runs: 10:00, 10:05, 10:10, 10:15
   ```

   and stores them both in DB and Redis.

3. **Scheduler library** at 10:00:01 pops “jobA@10:00”.

   * Sees generation = 1 matches → OK.
   * Marks precompute as dispatched.
   * Sends message to Kafka.

4. **Worker** runs the job at 10:00:02, logs appear in UI.

5. You edit jobA at 10:02 → generation increments to 2.

   * Old precomputed times (10:05, 10:10…) are invalid now.
   * Library automatically skips those since generation mismatched.
   * Precompute service soon recalculates new schedule entries with generation 2.

---

### 📏 How scalability is handled

* Redis stores upcoming runs using **sorted sets** (ZSET).
* To handle many jobs efficiently, there can be multiple Redis “shards” (like multiple small buckets).
* Each library instance “owns” one or more shards.
  Example:

  * Instance A handles jobs 0–3.
  * Instance B handles jobs 4–7.

This keeps the system horizontally scalable without central coordination.

---

### 📈 How we measure performance

The scheduler library exposes metrics like:

* How many jobs popped per second.
* How many published to Kafka.
* Average delay between scheduled time and actual dispatch (should be <1s ideally).

If metrics show lag increasing, you can add more shards or instances to balance load.

---

### ✅ In short

| Function           | Explanation                               |
| ------------------ | ----------------------------------------- |
| Decide timing      | Pops entries whose time ≤ now             |
| Avoid duplicates   | Uses generation check                     |
| Ensure reliability | Uses Redis + MongoDB consistency          |
| Scale              | Sharded Redis + distributed locks         |
| Communicate        | Publishes JSON events to Kafka            |
| Recover            | Retries failed publishes via Redis buffer |

---

## PART 2 — Frontend (React-Based UI)

### 🖥 Purpose

The frontend is the **control center** for everything.
It allows users to:

* Create and configure new jobs.
* View running and completed jobs.
* Watch live logs as containers run.
* Manage permissions and users.

---

### 🧭 Main Areas of the UI

| Page                        | Description                                      |
| --------------------------- | ------------------------------------------------ |
| Login                       | User login using username/password               |
| Dashboard                   | List of all scheduled jobs                       |
| Scheduler Create/Edit       | Create new or modify existing jobs               |
| Scheduler Detail            | Shows job info, execution history, and live logs |
| Users / Roles / Permissions | Admin-only screens for managing access           |

---

### 🧩 Example User Flow

1. A user logs in with username/password.

   * The backend returns a token and role.
   * The frontend remembers the role and hides unauthorized actions.

2. On the dashboard, the user sees all schedulers they can access.

   * Active jobs show a green badge; inactive ones are gray.
   * Each row has “Edit”, “Delete”, and “View Logs”.

3. When creating a job:

   * The user selects:

     * Docker image (from ECR list)
     * Job type (cron/interval/immediate)
     * Command (multi-line)
   * For cron type, a small helper shows next 5 execution times.
   * User sets permission for which roles can see or execute this job.

4. When the job runs:

   * The UI displays live terminal output — exactly as it appears in Docker.
   * ANSI colors (like green for success, red for error) are preserved.

5. After completion:

   * Logs are saved and visible in the history table.
   * User can filter history by success/failure, time, or who executed.

---

### 🧠 How real-time logs appear

When a job starts executing, the worker streams logs through WebSocket.

Example message sequence:

```json
{ "type": "log", "line": "Starting job..." }
{ "type": "log", "line": "Processing step 1..." }
{ "type": "log", "line": "\u001b[32mDone!\u001b[0m" }
{ "type": "end", "status": "success" }
```

Frontend behavior:

* Each incoming `log` message appends to the terminal view.
* The `ansi-to-html` utility converts colors.
* Auto-scroll follows the newest log line.
* When `type: "end"` arrives, UI shows job as “Finished”.

---

### 🧾 Scheduler Creation Form (explained)

#### Section 1: Scheduler Definition

* Text fields: job name, description
* Dropdown: Docker image (fetched from backend)
* Radio buttons: Job type (`immediate`, `cron`, `interval`)
* For “cron”, show input like `*/5 * * * *`
* For “interval”, show numeric seconds input
* Textarea: command (multi-line)

Validation:

* Command is required.
* Image is required.
* Cron expression must be valid.

#### Section 2: Permission Assignment

* Role dropdown (from backend)
* Three checkboxes: view, execute, edit
* When you select “developer” role and tick “execute”, it means developers can trigger this job manually.

#### Section 3: Preview

* Shows a summary:

  ```
  Image: aws.ecr/myetl:v2
  Command: ./start.sh
  Next 5 runs: 10:00, 10:05, 10:10, 10:15, 10:20
  ```
* “Save” button → sends POST to backend.

---

### 📜 Scheduler Detail View

This page combines three main panels:

1. **Overview**

   * Shows image, command, frequency, and owner.
   * Has buttons:

     * “Run Now” (manual trigger)
     * “Pause” or “Resume” (toggle active status)

2. **Execution History**

   * Table showing:

     * Start time, end time, duration, status, executed_by.
   * Each row has “View Logs” → opens log viewer.

3. **Live Logs**

   * Displays streaming log output when a job is running.
   * Uses `<pre>` tag for terminal-like view.
   * Converts ANSI colors to HTML.

---

### 🧍 User Management

* Admin can add new users with roles.
* Admin can modify roles and permissions directly from UI.
* Role permissions are automatically synced to backend collections.

---

### 💬 Example Interactions

* **Creating a job:**

  * Fill the form → click “Save” → see confirmation toast → job appears in dashboard.
* **Monitoring logs:**

  * Click “Logs” → terminal opens → see colored output live.
* **Editing a job:**

  * Change command → save → backend increments generation → old runs skipped automatically.
* **User without permission:**

  * Buttons like “Edit” or “Run Now” are hidden or disabled.

---

### 🧠 How the frontend communicates with backend

* All backend APIs follow REST style.
* Example requests:

  | Action             | Endpoint                           | Method    |
  | ------------------ | ---------------------------------- | --------- |
  | Get all schedulers | `/api/schedulers`                  | GET       |
  | Create scheduler   | `/api/schedulers`                  | POST      |
  | View history       | `/api/schedulers/:id/history`      | GET       |
  | Stream logs        | `/ws/logs?scheduler_id=x&run_id=y` | WebSocket |
  | Manage users       | `/api/users`                       | CRUD      |
  | Manage permissions | `/api/permissions`                 | CRUD      |

---

### 🧭 Overall UI Design Philosophy

| Principle                | Explanation                                         |
| ------------------------ | --------------------------------------------------- |
| **Clarity first**        | Avoid clutter. Tables and modals over nested menus. |
| **Real-time insight**    | Focus on live logs and job health.                  |
| **Permission awareness** | Show only what the current role can do.             |
| **Developer feel**       | Simple, dark-mode compatible interface.             |
| **Extensible**           | Easy to add metrics, alerts, and dashboards later.  |

---

### 🧪 Example: Real Use Scenario

**User:** DevOps engineer
**Task:** Monitor data pipeline job.

1. Opens Dashboard → sees “Nightly ETL” job active.
2. Clicks “Logs” → sees live terminal:

   ```
   [10:00:00] Starting job
   [10:00:05] Extracting source files
   [10:00:10] Transforming data...
   [10:00:15] ✅ Completed successfully
   ```
3. Log colors are preserved.
4. Checks “History” → sees all past runs with duration chart.
5. Admin edits job to change interval → new schedule starts next minute automatically.

---

## 🧩 Together: End-to-End Behavior Summary

1. User creates a job in the UI.
2. Backend stores it in DB → Precompute service calculates next 15 minutes of runs.
3. Scheduler library pops each due entry, publishes Kafka event.
4. Worker executes container and streams logs back.
5. UI shows live logs, updates status automatically.
6. When the run ends, the job’s status is stored in history and visible in table.

---

### 🧭 Summary Table

| Layer                  | What it does             | Key Concept       |
| ---------------------- | ------------------------ | ----------------- |
| **UI**                 | Create/edit/monitor jobs | Human interaction |
| **Backend API**        | CRUD + RBAC              | HTTP REST         |
| **Scheduler Library**  | Decide when to execute   | Redis + Mongo     |
| **Precompute Service** | Calculate future runs    | Time window       |
| **Worker**             | Execute containers       | Docker            |
| **Redis/Kafka**        | Transport & queue        | Decoupling layer  |

---

### ✅ Final Understanding (Natural Summary)

* The **scheduler library** acts like a conductor — it doesn’t play the music, but ensures every instrument (worker) starts at the right second.
* The **frontend** acts like a control panel — it lets humans see, adjust, and enjoy the performance in real time.
* Together, they create a reliable, observable orchestration system — one that’s developer-friendly, scalable, and future-ready.

---

Would you like me to now turn this explanation into a **downloadable `.md` file** (so you can open it in your editor or share it)?
