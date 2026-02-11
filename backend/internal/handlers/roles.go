package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/models"
)

type RoleHandlers struct {
	db *mongo.Database
}

func NewRoleHandlers(db *mongo.Database) *RoleHandlers {
	return &RoleHandlers{db: db}
}

type CreateRoleRequest struct {
	RoleName      string `json:"roleName" binding:"required"`
	Description   string `json:"description"`
	CanCreateTask bool   `json:"canCreateTask"`
}

type UpdateRoleRequest struct {
	RoleName      string `json:"roleName"`
	Description   string `json:"description"`
	CanCreateTask *bool  `json:"canCreateTask"`
	Active        *bool  `json:"active"`
}

type RoleResponse struct {
	ID            string    `json:"id"`
	RoleName      string    `json:"roleName"`
	Description   string    `json:"description"`
	Active        bool      `json:"active"`
	CanCreateTask bool      `json:"canCreateTask"`
	UserCount     int64     `json:"userCount"`
	IsSystem      bool      `json:"isSystem"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (h *RoleHandlers) GetRoles(c *gin.Context) {
	rolesCollection := h.db.Collection("roles")
	usersCollection := h.db.Collection("users")

	// Get all roles
	cursor, err := rolesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	defer cursor.Close(context.Background())

	var roles []RoleResponse
	for cursor.Next(context.Background()) {
		var role models.Role
		if err := cursor.Decode(&role); err != nil {
			continue
		}

		// Count users with this role
		userCount, _ := usersCollection.CountDocuments(context.Background(), bson.M{"role_id": role.ID})

		// Determine if it's a system role (created before a certain date or has specific names)
		isSystem := role.RoleName == "Administrator" || role.RoleName == "System Admin"

		roles = append(roles, RoleResponse{
			ID:            role.ID.Hex(),
			RoleName:      role.RoleName,
			Description:   role.Description,
			Active:        role.Active,
			CanCreateTask: role.CanCreateTask,
			UserCount:     userCount,
			IsSystem:      isSystem,
			CreatedAt:     role.CreatedAt,
			UpdatedAt:     role.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandlers) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
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

	// Check if role name already exists
	rolesCollection := h.db.Collection("roles")
	var existingRole models.Role
	err := rolesCollection.FindOne(context.Background(), bson.M{"role_name": req.RoleName}).Decode(&existingRole)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists"})
		return
	}

	// Create role
	now := time.Now()
	role := models.Role{
		RoleName:      req.RoleName,
		Description:   req.Description,
		Active:        true,
		CanCreateTask: req.CanCreateTask,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     &userCtx.UserID,
	}

	result, err := rolesCollection.InsertOne(context.Background(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	role.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, RoleResponse{
		ID:            role.ID.Hex(),
		RoleName:      role.RoleName,
		Description:   role.Description,
		Active:        role.Active,
		CanCreateTask: role.CanCreateTask,
		UserCount:     0,
		IsSystem:      false,
		CreatedAt:     role.CreatedAt,
		UpdatedAt:     role.UpdatedAt,
	})
}

func (h *RoleHandlers) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	rolesCollection := h.db.Collection("roles")

	// Build update document
	updateDoc := bson.M{
		"updated_at": time.Now(),
	}

	if req.RoleName != "" {
		updateDoc["role_name"] = req.RoleName
	}
	if req.Description != "" {
		updateDoc["description"] = req.Description
	}
	if req.CanCreateTask != nil {
		updateDoc["can_create_task"] = *req.CanCreateTask
	}
	if req.Active != nil {
		updateDoc["active"] = *req.Active
	}

	result, err := rolesCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateDoc},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func (h *RoleHandlers) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	rolesCollection := h.db.Collection("roles")
	usersCollection := h.db.Collection("users")

	// Check if any users are assigned to this role
	userCount, err := usersCollection.CountDocuments(context.Background(), bson.M{"role_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role usage"})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete role that is assigned to users"})
		return
	}

	result, err := rolesCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

func (h *RoleHandlers) GetPermissions(c *gin.Context) {
	// Return available permissions
	permissions := []map[string]interface{}{
		{
			"id":          "read",
			"name":        "Read",
			"description": "View schedulers and system information",
		},
		{
			"id":          "write",
			"name":        "Write",
			"description": "Create and modify schedulers",
		},
		{
			"id":          "delete",
			"name":        "Delete",
			"description": "Delete schedulers and data",
		},
		{
			"id":          "admin",
			"name":        "Admin",
			"description": "Full administrative access",
		},
	}

	c.JSON(http.StatusOK, permissions)
}
