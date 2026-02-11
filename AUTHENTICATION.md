# Authentication & Authorization System

## Overview

The Smart Task Orchestrator implements a secure JWT-based authentication and role-based authorization system.

## Key Features

### 1. Boot-Time Seeding
- **Single Admin Creation**: On first boot, the system creates a default admin user only if no admin exists
- **No Duplicate Seeding**: Subsequent boots skip seeding if admin user already exists
- **Default Credentials**:
  - Username: `admin`
  - Password: `12345`
  - Role: `admin`
  - Initial Login Flag: `true` (forces password change)

### 2. JWT Token System
- **Access Token**: 30-minute expiry
- **Refresh Token**: 7-day expiry
- **Token Claims Include**:
  - `user_id`: User's MongoDB ObjectID
  - `username`: User's username
  - `role_id`: User's role ObjectID
  - `is_initial_login`: Boolean flag indicating if password change is required

### 3. Password Change Flow

#### Initial Login (is_initial_login = true)
1. User logs in with default credentials
2. System returns `must_change_password: true` in login response
3. Frontend shows password change modal (cannot be closed)
4. User must change password before accessing the system
5. After successful password change:
   - `is_initial_login` is set to `false` in database
   - New JWT token is issued with updated flag
   - User is redirected to dashboard

#### Regular Password Change
1. User clicks "Change Password" in navbar
2. Modal opens (can be closed)
3. User enters old password and new password
4. System validates and updates password
5. New JWT token is issued
6. User continues working

#### UI Indicators
- **Red Badge in Navbar**: Shows when `is_initial_login = true`
  - Message: "Change password to secure your account"
  - Small, non-intrusive but visible
- **Password Change Button**: Always available in user menu

### 4. User Management (Admin Only)

#### Creating New Users
Only users with `role = admin` can create new users.

**Required Fields**:
- `username`: Unique username
- `password`: Minimum 8 characters
- `email`: Valid email address
- `roleId`: Valid role ObjectID

**Automatic Behavior**:
- All newly created users get `is_initial_login = true`
- Forces password change on first login
- User is created with `active = true`

**API Endpoint**:
```bash
POST /api/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "username": "newuser",
  "password": "temporary123",
  "email": "user@example.com",
  "roleId": "role_object_id"
}
```

### 5. Authorization Checks

#### Admin-Only Operations
The following operations require admin role:
- Create users (`POST /api/users`)
- Update users (`PUT /api/users/:id`)
- Delete users (`DELETE /api/users/:id`)
- Reset user passwords (`POST /api/users/:id/reset-password`)

#### Implementation
```go
// Check if user has admin role
var currentUser models.User
if err := usersCollection.FindOne(ctx, bson.M{"_id": userCtx.UserID}).Decode(&currentUser); err != nil {
    return http.StatusInternalServerError, "Failed to verify user permissions"
}

var currentRole models.Role
if err := rolesCollection.FindOne(ctx, bson.M{"_id": currentUser.RoleID}).Decode(&currentRole); err != nil || currentRole.RoleName != "admin" {
    return http.StatusForbidden, "Only administrators can perform this action"
}
```

## API Endpoints

### Authentication

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "12345"
}

Response:
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "must_change_password": true,
  "user": {
    "id": "...",
    "username": "admin",
    "email": "admin@orchestrator.local",
    "roleId": "...",
    "isInitialLogin": true
  }
}
```

#### Change Password
```bash
POST /api/auth/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "12345",
  "new_password": "newSecurePassword123"
}

Response:
{
  "message": "Password changed successfully",
  "success": true,
  "access_token": "eyJhbGc..."  // New token with is_initial_login = false
}
```

#### Refresh Token
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}

Response:
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

#### Get Current User
```bash
GET /api/me
Authorization: Bearer <token>

Response:
{
  "id": "...",
  "username": "admin",
  "email": "admin@orchestrator.local",
  "roleId": "...",
  "isInitialLogin": false
}
```

### User Management

#### List Users
```bash
GET /api/users
Authorization: Bearer <token>

Response:
[
  {
    "id": "...",
    "username": "admin",
    "email": "admin@orchestrator.local",
    "roleId": "...",
    "roleName": "admin",
    "isInitialLogin": false,
    "active": true,
    "lastLoginAt": "2025-11-16T03:14:07Z",
    "createdAt": "2025-11-16T03:14:07Z",
    "updatedAt": "2025-11-16T03:14:07Z"
  }
]
```

#### Create User (Admin Only)
```bash
POST /api/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "username": "newuser",
  "password": "temporary123",
  "email": "user@example.com",
  "roleId": "role_object_id"
}

Response:
{
  "id": "...",
  "username": "newuser",
  "email": "user@example.com",
  "roleId": "...",
  "roleName": "user",
  "isInitialLogin": true,
  "active": true,
  "createdAt": "2025-11-16T03:14:07Z",
  "updatedAt": "2025-11-16T03:14:07Z"
}
```

#### Reset User Password (Admin Only)
```bash
POST /api/users/:id/reset-password
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "newPassword": "temporary456"
}

Response:
{
  "message": "Password reset successfully"
}
```

## Security Best Practices

### Password Requirements
- Minimum 8 characters
- Must be different from current password
- Hashed using bcrypt with default cost (10)

### Token Security
- Tokens are signed with HS256 algorithm
- Secret key should be changed in production (set `JWT_SECRET` environment variable)
- Access tokens expire after 30 minutes
- Refresh tokens expire after 7 days

### Database Security
- Passwords are never stored in plain text
- Password hashes use bcrypt
- Unique indexes on username and email
- Active flag allows soft deletion

## Frontend Integration

### Login Flow
```typescript
const { login } = useAuth();

const handleSubmit = async (e: React.FormEvent) => {
  const result = await login(username, password);
  
  if (result.success) {
    if (result.user?.isInitialLogin) {
      // Show password change modal
      setShowChangePassword(true);
    } else {
      // Navigate to dashboard
      navigate('/');
    }
  }
};
```

### Logout Flow
```typescript
const { logout } = useAuth();

const handleLogout = () => {
  logout(); // Clears token and redirects to login
};
```

### Protected Routes
```typescript
<Route path="/" element={
  <ProtectedRoute>
    <Layout />
  </ProtectedRoute>
}>
  <Route index element={<Dashboard />} />
  {/* Other protected routes */}
</Route>
```

## Troubleshooting

### Issue: User can't login
- Verify credentials are correct (default: admin/12345)
- Check API logs for authentication errors
- Verify MongoDB is running and accessible

### Issue: Password change not working
- Ensure old password is correct
- Verify new password meets requirements (min 8 chars)
- Check that new password is different from old password

### Issue: Token expired
- Use refresh token to get new access token
- If refresh token is also expired, user must login again

### Issue: Admin user not created on boot
- Check MongoDB logs for connection errors
- Verify database indexes were created successfully
- Check API logs for seeding messages

## Database Schema

### Users Collection
```javascript
{
  _id: ObjectId,
  username: String (unique),
  email: String (unique, sparse),
  password_hash: String,
  role_id: ObjectId,
  is_initial_login: Boolean,
  active: Boolean,
  last_login_at: Date,
  created_at: Date,
  updated_at: Date,
  created_by: ObjectId,
  updated_by: ObjectId,
  metadata: Object
}
```

### Roles Collection
```javascript
{
  _id: ObjectId,
  role_name: String (unique),
  description: String,
  active: Boolean,
  can_create_task: Boolean,
  created_at: Date,
  updated_at: Date,
  created_by: ObjectId
}
```

## Environment Variables

```bash
# JWT Configuration
JWT_SECRET=your-jwt-secret-change-in-production

# Database
MONGO_URI=mongodb://mongodb:27017
DB_NAME=orchestrator

# Server
PORT=8080
```

## Testing

### Test Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "12345"}'
```

### Test Password Change
```bash
TOKEN="your_access_token"
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"old_password": "12345", "new_password": "newPassword123"}'
```

### Test User Creation (Admin)
```bash
TOKEN="admin_access_token"
ROLE_ID="role_object_id"
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "testuser",
    "password": "temporary123",
    "email": "test@example.com",
    "roleId": "'$ROLE_ID'"
  }'
```
