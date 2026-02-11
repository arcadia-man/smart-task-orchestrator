package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/models"
)

type UserHandlers struct {
	db *mongo.Database
}

func NewUserHandlers(db *mongo.Database) *UserHandlers {
	return &UserHandlers{db: db}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	RoleID   string `json:"roleId" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   string `json:"roleId"`
	Active   *bool  `json:"active"`
}

type UserResponse struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	RoleID         string     `json:"roleId"`
	RoleName       string     `json:"roleName"`
	IsInitialLogin bool       `json:"isInitialLogin"`
	Active         bool       `json:"active"`
	LastLoginAt    *time.Time `json:"lastLoginAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (h *UserHandlers) GetUsers(c *gin.Context) {
	usersCollection := h.db.Collection("users")

	// Aggregation pipeline to join users with roles
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "roles",
				"localField":   "role_id",
				"foreignField": "_id",
				"as":           "role",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$role",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":              1,
				"username":         1,
				"email":            1,
				"role_id":          1,
				"is_initial_login": 1,
				"active":           1,
				"last_login_at":    1,
				"created_at":       1,
				"updated_at":       1,
				"role_name":        "$role.role_name",
			},
		},
		{
			"$sort": bson.M{"created_at": -1},
		},
	}

	cursor, err := usersCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(context.Background())

	var users []UserResponse
	for cursor.Next(context.Background()) {
		var result struct {
			ID             primitive.ObjectID `bson:"_id"`
			Username       string             `bson:"username"`
			Email          string             `bson:"email"`
			RoleID         primitive.ObjectID `bson:"role_id"`
			IsInitialLogin bool               `bson:"is_initial_login"`
			Active         bool               `bson:"active"`
			LastLoginAt    *time.Time         `bson:"last_login_at"`
			CreatedAt      time.Time          `bson:"created_at"`
			UpdatedAt      time.Time          `bson:"updated_at"`
			RoleName       string             `bson:"role_name"`
		}

		if err := cursor.Decode(&result); err != nil {
			continue
		}

		users = append(users, UserResponse{
			ID:             result.ID.Hex(),
			Username:       result.Username,
			Email:          result.Email,
			RoleID:         result.RoleID.Hex(),
			RoleName:       result.RoleName,
			IsInitialLogin: result.IsInitialLogin,
			Active:         result.Active,
			LastLoginAt:    result.LastLoginAt,
			CreatedAt:      result.CreatedAt,
			UpdatedAt:      result.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	usersCollection := h.db.Collection("users")
	var currentUser models.User
	if err := usersCollection.FindOne(context.Background(), bson.M{"_id": userCtx.UserID}).Decode(&currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user permissions"})
		return
	}

	rolesCollection := h.db.Collection("roles")
	var currentRole models.Role
	if err := rolesCollection.FindOne(context.Background(), bson.M{"_id": currentUser.RoleID}).Decode(&currentRole); err != nil || currentRole.RoleName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can create users"})
		return
	}

	// Validate role exists
	roleID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	var role models.Role
	if err = rolesCollection.FindOne(context.Background(), bson.M{"_id": roleID}).Decode(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	// Check if username already exists
	var existingUser models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	err = usersCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	now := time.Now()
	user := models.User{
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   string(passwordHash),
		RoleID:         roleID,
		IsInitialLogin: true, // Force password change on first login
		Active:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
		CreatedBy:      &userCtx.UserID,
		Metadata:       make(map[string]interface{}),
	}

	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, UserResponse{
		ID:             user.ID.Hex(),
		Username:       user.Username,
		Email:          user.Email,
		RoleID:         user.RoleID.Hex(),
		RoleName:       role.RoleName,
		IsInitialLogin: user.IsInitialLogin,
		Active:         user.Active,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	})
}

func (h *UserHandlers) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	usersCollection := h.db.Collection("users")

	// Build update document
	updateDoc := bson.M{
		"updated_at": time.Now(),
		"updated_by": userCtx.UserID,
	}

	if req.Username != "" {
		updateDoc["username"] = req.Username
	}
	if req.Email != "" {
		updateDoc["email"] = req.Email
	}
	if req.RoleID != "" {
		roleID, err := primitive.ObjectIDFromHex(req.RoleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}
		updateDoc["role_id"] = roleID
	}
	if req.Active != nil {
		updateDoc["active"] = *req.Active
	}

	result, err := usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateDoc},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	usersCollection := h.db.Collection("users")

	result, err := usersCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandlers) ResetPassword(c *gin.Context) {
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	type ResetPasswordRequest struct {
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	usersCollection := h.db.Collection("users")
	result, err := usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"password_hash":    string(passwordHash),
				"is_initial_login": true, // Force password change on next login
				"updated_at":       time.Now(),
				"updated_by":       userCtx.UserID,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
