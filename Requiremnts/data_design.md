# 1) Schema (data model) — complete and concrete

Notes before the schemas:

* Primary choice: you stated **MongoDB** as DB. I’ll present schemas optimized for Mongo, with suggested indexes and some SQL-style DDL equivalents where helpful.
* Use UUIDs (v4) as IDs for cross-service stability. In Mongo store them as strings (e.g., `"6f1e3c2b-..."`) or as `Binary(UUID)` if you prefer compact binary form.
* Timestamps stored in UTC ISO format.
* Keep `generation` on scheduler definitions (integer) for invalidation logic.
* Add `created_by`, `updated_by` to audit changes (string user id).
* Use unique constraints for things like `username`, `(scheduler_id, run_at, generation)` etc.

---

## A. `users` collection (Mongo)

Purpose: store accounts, initial-login flow, and metadata.

Fields:

```text
users:
{
  _id: UUID,                // user id (primary)
  username: string,         // unique, lowercased
  email: string,            // optional, unique if present
  password_hash: string,    // bcrypt hash
  role_id: UUID,            // FK -> roles._id
  is_initial_login: bool,   // default true -> forces password change
  active: bool,             // default true
  last_login_at: datetime,
  created_at: datetime,
  updated_at: datetime,
  created_by: UUID|null,
  updated_by: UUID|null,
  metadata: { ... }         // optional free-form user meta
}
```

Indexes & constraints:

* Unique index on `username` (case-insensitive ideally).
* Unique index on `email` if used.
* Index on `role_id` for RBAC lookups.

Sample document:

```json
{
  "_id": "uuid-1234",
  "username": "admin",
  "email": "ops@example.com",
  "password_hash": "$2b$12$...",
  "role_id": "role-0001",
  "is_initial_login": true,
  "active": true,
  "last_login_at": null,
  "created_at": "2025-10-22T10:00:00Z"
}
```

Security note:

* Never store plaintext passwords. Use bcrypt with cost >= 12 (tune for your CPU).
* Do not return `password_hash` in any API response.

---

## B. `roles` collection

Purpose: role definitions and canned permission flags for UI defaults.

Fields:

```text
roles:
{
  _id: UUID,
  role_name: string,      // unique, e.g., "admin", "developer"
  description: string,
  active: bool,           // default true
  created_at: datetime,
  updated_at: datetime,
  created_by: UUID|null,
  can_create_tak: bool
}
```

Indexes:

* Unique index on `role_name`.

Notes:

* `can_create_tak` here are coarse convenience flags. The fine-grained enforcement will rely on the `can_create_tak` collection (below). this data will be true in dev for all role but in production for admin only it will true.
* Allow dynamic role creation (UI + API create role) — per your requirements.

---

## C. `permissions` collection (role-scoped scheduler permissions)

Purpose: map which roles can do what on which scheduler_definition.

Fields:

```text
permissions:
{
  _id: UUID,
  role_id: UUID,            // FK -> roles._id
  scheduler_id: UUID,       // FK -> scheduler_definition._id
  role_see: bool,           // view permission
  role_execute: bool,       // can run / trigger
  role_alter_config: bool,  // update scheduler settings
  created_at: datetime,
  created_by: UUID|null
}
```

Indexes:

* Compound unique index on `(role_id, scheduler_id)` to avoid duplicates.
* Index on `scheduler_id` and `role_id` for lookups.

Notes:

* For UI: fetch permissions by scheduler to show which roles can do what.
* For middleware: query quickly by role_id + scheduler_id to check permissions.

---

## D. `images` collection

Purpose: list known container images (ECR URIs), versions, and meta.

Fields:

```text
images:
{
  _id: UUID,
  registry_url: string,    // e.g., "123456789012.dkr.ecr.us-east-1.amazonaws.com"
  image: string,           // full image name e.g., "repo/name:tag"
  name: string,            // friendly name
  description: string,
  version: string,
  is_default: bool,
  created_at: datetime,
  created_by: UUID|null
}
```

Indexes:

* Unique index on `registry_url + image` or `image` depending on usage.

---

## E. `scheduler_definition` collection

Purpose: immutable (config-only) job definitions. DO NOT store computed run times here.

Fields:

```text
scheduler_definition:
{
  _id: UUID,
  name: string,
  description: string,
  image: string,           // ECR path or reference to images collection
  job_type: string,        // "immediate", "cron", "interval"
  cron_expr: string|null,  // cron string if job_type == cron
  interval_seconds: int|null,
  command: string,         // multi-line script
  status: string,          // "active", "disable" etc.
  timezone: string,        // e.g., "UTC" or "Asia/Kolkata"
  generation: int,         // increment on each update
  created_at: datetime,
  updated_at: datetime,
  created_by: UUID|null,
  updated_by: UUID|null
}
```

Indexes:

* Index on `status` (for precompute to find active jobs).
* Optional index on `name` (unique if desired).
* `generation` is an integer used for invalidation — include in updates.

Important:

* **Do not** add next_run_at here (you explicitly want no computed fields).
* Use `generation` to tag precomputed rows.

Sample doc:

```json
{
  "_id": "sched-101",
  "name": "Nightly ETL",
  "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/etl:v3",
  "job_type": "cron",
  "cron_expr": "0 2 * * *", // daily at 02:00
  "timezone": "UTC",
  "command": "bash /app/run_etl.sh",
  "generation": 3,
  "status": "active"
}
```

---

## F. `scheduler_precompute` collection (durable queue backup)

Purpose: durable storage of precomputed future run items (source-of-truth for recovery if Redis is lost).

Fields:

```text
scheduler_precompute:
{
  _id: UUID,              // precompute_id
  scheduler_id: UUID,
  run_at: datetime,       // when job should run (UTC)
  generation: int,        // generation at compute time
  status: string,         // 'pending', 'dispatched', 'canceled', 'failed', 'discarded' (for abrublty cancelling)
  created_at: datetime
}
```

Indexes & constraints:

* Unique compound index on `(scheduler_id, run_at, generation)` to avoid duplicates.
* Index on `run_at` for range queries (`run_at <= now`).
* Index on `(status, run_at)` to find pending runs before a time.

Notes:

* Precompute service writes these rows and also pushes members into Redis ZSETs with score `run_at_ms`.

---

## G. `scheduler_history` collection (execution logs & status)

Purpose: store execution metadata, logs (or reference to logs), exit code, and actor.

Fields:

```text
scheduler_history:
{
  _id: UUID,              // run id
  scheduler_id: UUID,
  precompute_id: UUID|null,
  run_id: UUID,           // same as _id or duplicate id
  executed_by: UUID,      // user id if invoked manually, or system id
  status: string,         // 'pending', 'running', 'success', 'failed', 'discarded'
  start_time: datetime,
  end_time: datetime|null,
  command: sting // sript ran at that time for autdit
  exit_code: int|null,
  process_id: string,     // container id
  log_blob_id: string|null, // pointer to full log (S3 or DB blob for now)
  log_text: string|null,  // short text or null (if storing in S3)
  created_at: datetime,
  updated_at: datetime,
  error_message: string|null
}
```

Indexes:

* Index on `scheduler_id` + `start_time` for history queries.
* Index on `run_id` unique.

Logging note:

* For large logs, store the full logs in S3 or long-term store and keep a pointer here. Keep a short `log_text` tail or summary for fast UI display.

---


## I. Redis data model (fast queue / streaming)

Keys & structure:

* **Sharded sorted sets** for scheduled items:

  * `sched:zset:<shard>`

    * Member: `precompute:<precompute_id>` or `scheduler:<scheduler_id>|<run_at_ms>|<generation>` (compact string)
    * Score: `run_at_ms` (Unix milliseconds)

* **Meta hash (optional)**:

  * `sched:meta:<precompute_id>` => small JSON/hash containing `{scheduler_id, run_at_ms, generation}` for quick access without DB roundtrip

* **Pending publish list**:

  * `sched:pending_publish` (list) — stores serialized messages to try publishing again.

* **Per-scheduler index (optional)**:

  * `sched:ids:<scheduler_id>` => set or list of precompute ids for easier invalidation.

* **Log streaming**:

  * Redis stream `logs:<run_id>` or a Message broker for WebSocket gateway.

Notes:

* Use sharding by `hash(scheduler_id) % N` to split load across multiple ZSET keys.
* For now there is no actual shar you can think there is only one shard on redis.

---

# 2) Authorization & Authentication — implementation plan

This section is the full, step-by-step description of how to implement authentication, RBAC, first-login password change, session/token lifecycle, authorization middleware, audit logging, and security hardening. All written as plain instructions and examples.

---

## A. Goals & properties

* **Secure authentication** (passwords stored hashed).
* **Initial admin seeding** on first startup (username `admin`, password `admin`) but requiring password change on first login.
* **Role-Based Access Control (RBAC)** enforced at API and UI level.
* **Fine-grained scheduler permissions** per role per scheduler.
* **Audit logging** for critical actions.
* **Usable tokens** (JWT) for stateless API authentication with secure refresh flow.
* **Brute-force / lockout protection** and monitoring.

---

## B. Authentication flow (step-by-step)

### 1. Password storage & verification

* Use **bcrypt** (or Argon2 if available in Go lib) for password hashing.

  * bcrypt cost factor e.g., 12 (tune for CPU).
* On user creation: compute `password_hash = bcrypt(password)`, store `password_hash` in `users` doc.
* On login:

  * Lookup `users` by `username` (case-sensitive).
  * If user not found or `active=false` → return 401.
  * Compare supplied password with stored `password_hash` using bcrypt `Compare`.
  * If match, proceed; else increment failed-login counter and return 401.

### 2. Initial admin user & first-login

* On app startup (DB init), check `users` collection. If zero users, seed:

  * username: `admin`, password: `admin`, `is_initial_login = true`, role `admin`.
* On login, if `is_initial_login == true`:

  * Allow login but **force** immediate redirect to "change password" flow.
  * After password change: set `is_initial_login = false`.
* Rationale: prevents a default admin password staying in use.

### 3. Token-based sessions (JWT)

* After successful password verification, produce a **JWT access token** and a **refresh token**.

  * Access token short TTL (e.g., 15m). Refresh token longer (e.g., 7 days).
  * JWT payload minimal: `{ sub: user_id, username, role_id, iat, exp, jti }`. keep expiry 30min
* Store refresh tokens in DB (or Redis) keyed by `refresh_token_id` and `user_id` for revocation support.
* Return tokens to client:

  * Access token in response body and optionally cookie (HttpOnly) if using browser cookies.
  * Refresh token stored in HttpOnly secure cookie (recommended) or persisted with careful client handling.

### 4. Token refresh

* Endpoint: `POST /api/auth/refresh` with refresh token.
* Validate refresh token exists and not revoked.
* Issue new access token (and optionally new refresh token).
* Save rotation: rotate refresh tokens to prevent reuse (invalidate old when issuing new).

### 5. Logout

* Invalidate refresh token in DB/Redis so it can’t be used again.
* Revoke active access tokens if you keep token blacklist (optional).

---

## C. Authorization (RBAC) — conceptual design

We apply **two levels** of permission checks:

1. **Coarse role-level checks** (based on `roles.permissions`): used to gate access to admin screens (e.g., only roles with `can_manage_users` can access `/api/users` endpoints).

2. **Fine-grained per-scheduler checks** (based on `permissions` collection): used to permit actions on a per-scheduler basis (view, execute, alter_config).

### Typical operations and checks:

* `GET /api/schedulers` → roles with `can_see_all` OR filter results to only schedulers where the current role has `role_see == true`.
* `POST /api/schedulers` → require role `can_create_task` (coarse).
* `PATCH /api/schedulers/:id` → require `permissions.role_alter_config` for that scheduler or be `admin`.
* `POST /api/schedulers/:id/run` → require `permissions.role_execute` for that scheduler or be `admin`.

---

## D. Middleware & enforcement pattern (how to implement)

### 1. Authentication middleware (per request)

* Extract token (Authorization header `Bearer <token>` or cookie).
* Validate signature and expiry.
* Attach `user_id`, `role_id` to request context so downstream handlers can use it.

### 2. Authorization helper function

* Implement helper: `func RequirePermission(ctx, permissionName, schedulerId) error`

  * If permissionName is coarse (application-level), check `roles.permissions[permissionName]`. If true, pass.
  * Otherwise (scheduler-level):

    * Lookup `permissions` collection with `(role_id, scheduler_id)`.
    * If found and the specific `roleExecute`/`roleSee`/`roleAlterConfig` is true → allow.
    * Else deny 403.

Example usage:

* In handler for `/api/schedulers/:id/run`:

  * call `RequirePermission(ctx, "execute", schedulerId)`.

### 3. Caching & performance

* To avoid DB call for every permission check, cache role-permission maps in memory for short TTL (e.g., 30s) or use Redis cache keyed by `role_id:scheduler_id`. Invalidation on role/permission change is required.

### 4. UI-level awareness

* When the UI loads scheduler detail, include computed permission flags for the current user in the API response:

  ```json
  {
    "scheduler": { ... },
    "current_user_permissions": {
      "canSee": true,
      "canExecute": false,
      "canAlterConfig": false
    }
  }
  ```
* This reduces the need for the frontend to call a separate endpoint to detect permissions.

---

## E. First-login password change flow (detailed)

1. User `admin` logs in with default password.
2. Backend returns a short-lived access token but also a flag `must_change_password: true`.
3. Frontend redirects to a dedicated change-password page that blocks other actions.
4. User submits new password (validate strength).
5. Backend verifies token, updates `password_hash` with new bcrypt hash, sets `is_initial_login = false`. Optionally revoke other sessions.
6. On success, proceed to normal app flow.

Security:

* Enforce password policy (min length, character classes, no common passwords).
* Optionally enforce MFA for admin accounts.

---

## F. Secure account & brute-force protection

1. **Failed login counters**

   * Keep per-username and per-IP failed attempt counters in Redis with TTL (e.g., block after 5 failed attempts for 15 minutes).
2. **CAPTCHA**

   * After N failed attempts, require CAPTCHA to continue to prevent scripting.
3. **Rate limiting**

   * Limit auth endpoints per IP address globally (reverse-proxy or middleware).
4. **Lockout**

   * For suspicious activity, lock account and notify admin.

---

## G. Role creation & auto-population (your requirement)

You wanted roles to be auto-populated and created when typed in UI:

Implementation:

* When creating/updating a scheduler via API, accept `role_name` strings in the `permissions` payload.
* For each `role_name` not found in `roles` collection:

  * Create a new role record with default permissions (e.g., `can_create_task=false`, `can_execute=false`) and `active=true`.
  * Return role id to the caller.
* Record this action in `audit_logs` and notify admin if you want approvals for new roles.

Security:

* Consider an admin approval step for new roles to avoid privilege creep. But per your requirement, auto-creation is allowed.

---

## H. Manual run & audit trail

* `POST /api/schedulers/:id/run` should:

  * Check `role_execute` permission.
  * Create an entry in `scheduler_history` with `status=pending`, `executed_by=user_id`.
  * Publish a job_execution event to Kafka with `manual: true` in metadata.
  * Record an `audit_logs` entry: `{action: "run-manual", entity_id: scheduler_id, performed_by: user_id, details: {...}}`.

---

## I. Token revocation & session invalidation

* Store active refresh tokens in DB with `jti` and expiry, enabling revocation.
* For critical actions (role permission change, password reset):

  * Revoke existing refresh tokens for that user (and optionally invalidate access tokens via short TTL or a blacklist).
* For admins, consider immediate session revocation on role demotion.

---

## J. Permission evaluation examples (plain English)

1. **API: Edit scheduler**

   * Check: user role is `admin` OR `permissions.role_alter_config` for that scheduler is true.
   * If yes: allow update. On update: increment `scheduler_definition.generation`, cancel pending precompute rows (mark canceled) and notify precompute worker to re-compute.

2. **API: Trigger run**

   * Check: user role `admin` OR `permissions.role_execute` for that scheduler is true.
   * If allowed: create `scheduler_history` run, publish run event to Kafka.

3. **API: View list of schedulers**

   * If user role is `admin` or `can_see_all` flag: return full list.
   * Else: filter to only schedulers where `permissions` contains `role_see = true` for that user’s role. (This can be done efficiently by querying `permissions` collection first to get scheduler IDs.)

---

## K. Logging & audit (must do)

Whenever a security-sensitive action occurs, write an `audit_logs` entry:

* Actions: `login_success`, `login_failed`, `password_change`, `create_scheduler`, `update_scheduler`, `delete_scheduler`, `permission_change`, `manual_run`, `role_create`.
* Each entry should record `performed_by`, `performed_at`, `entity_type`, `entity_id`, and `details` (e.g., diff of change).

Keep audit retention policy (e.g., 90 days in DB then archive to S3 for 7 years depending on compliance).

---

## L. UI & API contract quick notes for auth

* Login endpoint: `POST /api/auth/login` → body `{username, password}` → returns `{ access_token, refresh_token, must_change_password, role_id, user_id }`
* Change password endpoint: `POST /api/auth/change-password` → requires access token; body `{ old_password, new_password }` or for forced-change include `token` & new password.
* Refresh endpoint: `POST /api/auth/refresh` → accepts refresh token → returns new access token (+ optional new refresh token).
* Logout endpoint: `POST /api/auth/logout` → invalidates refresh token.
* A `GET /api/me` endpoint to return `{user_id, username, role_id, role_permissions}` for UI use.

---

## N. Example real-world sequences (to verify implementation)

### Sequence 1: New user & first-login

1. Admin seeds default admin if none present.
2. Admin logs in with `admin/admin`.
3. API returns token + `must_change_password = true`.
4. Admin navigates to change-password page, sets new strong password.
5. API updates user (`is_initial_login = false`) and rotates sessions.

### Sequence 2: Edit scheduler & invalidation

1. User edits scheduler A; API validates `role_alter_config`.
2. API increments `generation` for scheduler A.
3. API marks pending precompute rows with older generation as `canceled`.
4. Precompute service is triggered to recalc items with new generation.
5. UI sees updated next-run preview and permissions unaffected.

### Sequence 3: Manual run with permission

1. User clicks “Run Now”.
2. UI calls `POST /api/schedulers/:id/run`.
3. API checks permission; creates `scheduler_history` run entry; publishes event to Kafka.
4. Worker consumes and streams logs; UI subscribes to WebSocket for logs.

---