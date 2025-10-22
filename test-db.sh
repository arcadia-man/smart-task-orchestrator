#!/bin/bash

# Test database connection and seeding

echo "🔍 Testing Database Connection and Seeding"
echo "=========================================="

# Test MongoDB connection
echo ""
echo "📊 Testing MongoDB connection..."
if mongosh --quiet --eval "db.runCommand('ping')" orchestrator > /dev/null 2>&1; then
    echo "✅ MongoDB connection successful"
else
    echo "❌ MongoDB connection failed"
    exit 1
fi

# Check if admin user was created
echo ""
echo "👤 Checking if admin user was seeded..."
USER_COUNT=$(mongosh --quiet --eval "db.users.countDocuments()" orchestrator)
echo "   Users in database: $USER_COUNT"

if [ "$USER_COUNT" -gt 0 ]; then
    echo "✅ Admin user seeding successful"
    
    # Show admin user details (without password)
    echo ""
    echo "📋 Admin user details:"
    mongosh --quiet --eval "db.users.findOne({username: 'admin'}, {password_hash: 0})" orchestrator
else
    echo "❌ No users found - seeding may have failed"
fi

# Check roles
echo ""
echo "🔐 Checking roles..."
ROLE_COUNT=$(mongosh --quiet --eval "db.roles.countDocuments()" orchestrator)
echo "   Roles in database: $ROLE_COUNT"

if [ "$ROLE_COUNT" -gt 0 ]; then
    echo "✅ Role seeding successful"
    echo ""
    echo "📋 Available roles:"
    mongosh --quiet --eval "db.roles.find({}, {role_name: 1, description: 1, can_create_task: 1})" orchestrator
fi

# Check images
echo ""
echo "🖼️  Checking default images..."
IMAGE_COUNT=$(mongosh --quiet --eval "db.images.countDocuments()" orchestrator)
echo "   Images in database: $IMAGE_COUNT"

if [ "$IMAGE_COUNT" -gt 0 ]; then
    echo "✅ Image seeding successful"
    echo ""
    echo "📋 Available images:"
    mongosh --quiet --eval "db.images.find({}, {name: 1, image: 1, version: 1})" orchestrator
fi

echo ""
echo "🎉 Database testing complete!"