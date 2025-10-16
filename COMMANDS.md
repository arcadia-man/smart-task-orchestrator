# Smart Task Orchestrator - Commands Reference

## 🚀 **MAIN COMMANDS**

### Start the Application
```bash
cd smart-task-orchestrator
./run.sh start
```

### Stop the Application  
```bash
./run.sh stop
```

### Check Status
```bash
./run.sh status
```

### Restart Everything
```bash
./run.sh restart
```

### Test the API
```bash
./test.sh
```

---

## 📊 **WHAT GETS STARTED**

When you run `./run.sh start`:

1. **MongoDB** (Docker container on port 27017)
2. **API Server** (Go process on port 8080)
3. **Worker Process** (Go process, consumes from Kafka)
4. **Scheduler Process** (Go process, handles cron jobs)
5. **Frontend** (React dev server on port 3000)

---

## 🔍 **TROUBLESHOOTING**

### If services won't start:
```bash
# Check what's using the ports
lsof -ti:8080 | xargs kill -9  # Kill API
lsof -ti:3000 | xargs kill -9  # Kill Frontend

# Force stop everything
./run.sh stop
pkill -f "go run"
pkill -f "npm run dev"
```

### View logs:
```bash
tail -f logs/api.log      # API server logs
tail -f logs/worker.log   # Worker process logs
tail -f logs/scheduler.log # Scheduler logs
tail -f logs/frontend.log # Frontend logs
```

### Manual API test:
```bash
# Create a job
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Job", "type": "immediate", "payload": {"test": true}, "maxRetries": 3}'

# Get all jobs
curl http://localhost:8080/api/jobs

# Check if frontend is accessible
curl http://localhost:3000
```

---

## 🎯 **QUICK WORKFLOW**

1. **First time setup:**
   ```bash
   cd smart-task-orchestrator/frontend
   npm install
   cd ..
   ```

2. **Daily usage:**
   ```bash
   ./run.sh start    # Start everything
   ./test.sh         # Test the API
   # Use the app at http://localhost:3000
   ./run.sh stop     # Stop everything when done
   ```

3. **If something goes wrong:**
   ```bash
   ./run.sh stop     # Stop everything
   ./run.sh start    # Start fresh
   ```

---

## ✅ **SUCCESS INDICATORS**

After running `./run.sh start`, you should see:
- ✅ API Server: http://localhost:8080
- ✅ Frontend: http://localhost:3000  
- ✅ MongoDB: localhost:27017
- ✅ Kafka: localhost:9092

After running `./test.sh`, you should see:
- ✅ API is running
- ✅ Job created successfully
- ✅ Found X jobs
- ✅ Job status: queued/completed
- ✅ Cron job created successfully

---

## 🚫 **WHAT TO AVOID**

❌ Don't run `docker-compose up` (we use local processes now)  
❌ Don't manually start Go processes (use the run script)  
❌ Don't forget to stop services when done (prevents port conflicts)  
❌ Don't start if Kafka isn't running (prerequisite check will fail)

---

## 📱 **USING THE APPLICATION**

1. **Web Dashboard**: http://localhost:3000
   - View all jobs
   - Create new jobs
   - Monitor job status in real-time
   - Retry failed jobs

2. **API Endpoints**: http://localhost:8080/api
   - `GET /jobs` - List all jobs
   - `POST /jobs` - Create new job
   - `GET /jobs/:id` - Get job details
   - `POST /jobs/:id/retry` - Retry job

3. **Job Types**:
   - **Immediate**: Executes right away
   - **Cron**: Executes on schedule (e.g., "*/5 * * * *" = every 5 minutes)

That's it! The application is now properly organized and easy to use. 🎉