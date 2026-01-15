# Quick Start Guide - Smart Task Orchestrator

## 🚀 Access the Application

### URLs
- **Frontend**: http://localhost:3001
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Default Login
- **Username**: `admin`
- **Password**: `12345`

⚠️ **Important**: You will be required to change the password on first login for security.

---

## 📋 First Login Steps

1. **Open the application**: http://localhost:3001

2. **Login with default credentials**:
   - Username: `admin`
   - Password: `12345`

3. **Change your password** (forced):
   - A modal will appear that cannot be closed
   - Enter current password: `12345`
   - Enter new password (min 8 characters)
   - Confirm new password
   - Click "Change Password"

4. **You're in!**
   - You'll be redirected to the dashboard
   - The red "Change password" badge will disappear

---

## 🔐 Password Requirements

- Minimum 8 characters
- Must be different from current password
- Cannot reuse old password

---

## 👥 Creating New Users (Admin Only)

1. Navigate to **User Management** in the sidebar

2. Click **Create User** button

3. Fill in the form:
   - Username (unique)
   - Email (unique)
   - Password (min 8 chars)
   - Role (select from dropdown)

4. Click **Create**

5. **Important**: All new users will be forced to change their password on first login

---

## 🔄 Managing Services

### Start Services
```bash
./run.sh start
```

### Stop Services
```bash
./run.sh stop
```

### Restart Services
```bash
./run.sh restart
```

### Check Status
```bash
./run.sh status
```

### View Logs
```bash
# All services
./run.sh logs

# Specific service
./run.sh logs api
./run.sh logs frontend
./run.sh logs mongodb
```

---

## 🧪 Testing the API

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

### Get Current User
```bash
TOKEN="your_access_token"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/me
```

### Change Password
```bash
TOKEN="your_access_token"
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_password"
  }'
```

---

## 📊 Dashboard Features

### Available Pages
- **Dashboard**: Overview and statistics
- **Schedulers**: Manage scheduled tasks
- **Create Scheduler**: Create new scheduled task
- **User Management**: Manage users (admin only)
- **Roles & Permissions**: Manage roles (admin only)
- **Docker Images**: Manage container images
- **System Logs**: View system logs
- **Monitoring**: System health monitoring

---

## 🔧 Troubleshooting

### Can't Login?
- Verify you're using the correct credentials
- Check if services are running: `./run.sh status`
- Check API logs: `./run.sh logs api`

### Password Change Not Working?
- Ensure old password is correct
- New password must be at least 8 characters
- New password must be different from old password

### Frontend Not Loading?
- Check if frontend is running: `./run.sh status`
- Verify port 3001 is not in use
- Check frontend logs: `./run.sh logs frontend`

### API Not Responding?
- Check if API is running: `./run.sh status`
- Verify port 8080 is not in use
- Check API logs: `./run.sh logs api`
- Test health endpoint: `curl http://localhost:8080/health`

---

## 📚 Documentation

- **AUTHENTICATION.md** - Complete authentication system documentation
- **FIXES_APPLIED.md** - Recent fixes and improvements
- **README.md** - Project overview and setup

---

## 🛡️ Security Notes

1. **Change Default Password**: Always change the default admin password immediately
2. **JWT Secret**: Change `JWT_SECRET` environment variable in production
3. **HTTPS**: Use HTTPS in production environments
4. **Database**: Secure MongoDB with authentication in production
5. **Firewall**: Restrict access to ports in production

---

## 🎯 Common Tasks

### Reset Admin Password (if forgotten)
```bash
# Stop services
./run.sh stop

# Remove MongoDB data volume
docker volume rm smart-task-orchestrator_mongodb_data

# Start services (will recreate admin with default password)
./run.sh start

# Login with admin/12345 and change password
```

### View All Users
```bash
TOKEN="admin_access_token"
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/users
```

### Create Scheduler
1. Navigate to **Create Scheduler** page
2. Fill in scheduler details
3. Set schedule (cron expression)
4. Select Docker image
5. Configure environment variables
6. Click **Create**

---

## 📞 Support

For issues or questions:
1. Check the logs: `./run.sh logs`
2. Review documentation in this repository
3. Check service status: `./run.sh status`

---

## ✅ System Health Check

Run this command to verify everything is working:

```bash
# Check all services
./run.sh status

# Test API health
curl http://localhost:8080/health

# Test frontend
curl -I http://localhost:3001

# Test login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "12345"}'
```

All should return successful responses.

---

## 🎉 You're Ready!

Your Smart Task Orchestrator is now fully configured and ready to use. Enjoy scheduling and managing your tasks!
