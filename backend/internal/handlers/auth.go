package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/models"
)

type AuthHandlers struct {
	db         *mongo.Database
	jwtManager *auth.JWTManager
}

func NewAuthHandlers(db *mongo.Database, jwtManager *auth.JWTManager) *AuthHandlers {
	return &AuthHandlers{
		db:         db,
		jwtManager: jwtManager,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken        string   `json:"access_token"`
	RefreshToken       string   `json:"refresh_token"`
	MustChangePassword bool     `json:"must_change_password"`
	User               UserInfo `json:"user"`
}

type UserInfo struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	RoleID         string `json:"roleId"`
	IsInitialLogin bool   `json:"isInitialLogin"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Find user by username
	usersCollection := h.db.Collection("users")
	var user models.User
	err := usersCollection.FindOne(context.Background(), bson.M{
		"username": req.Username,
		"active":   true,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.RoleID, user.IsInitialLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Update last login
	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login_at": time.Now()}},
	)
	if err != nil {
		// Log error but don't fail the login
		// log.Printf("Failed to update last login: %v", err)
	}

	response := LoginResponse{
		AccessToken:        tokenPair.AccessToken,
		RefreshToken:       tokenPair.RefreshToken,
		MustChangePassword: user.IsInitialLogin,
		User: UserInfo{
			ID:             user.ID.Hex(),
			Username:       user.Username,
			Email:          user.Email,
			RoleID:         user.RoleID.Hex(),
			IsInitialLogin: user.IsInitialLogin,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandlers) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Find user in database
	usersCollection := h.db.Collection("users")
	var user models.User
	err := usersCollection.FindOne(context.Background(), bson.M{
		"_id": userCtx.UserID,
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Validate new password
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 8 characters long"})
		return
	}

	if req.OldPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be different from current password"})
		return
	}

	// Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password in database
	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": userCtx.UserID},
		bson.M{
			"$set": bson.M{
				"password_hash":    string(newPasswordHash),
				"is_initial_login": false,
				"updated_at":       time.Now(),
				"updated_by":       userCtx.UserID,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Generate new tokens with updated is_initial_login flag
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.RoleID, false)
	if err != nil {
		// Password was changed successfully, but token generation failed
		// Still return success
		c.JSON(http.StatusOK, gin.H{
			"message": "Password changed successfully",
			"success": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Password changed successfully",
		"success":      true,
		"access_token": tokenPair.AccessToken,
	})
}

func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Refresh the token
	tokenPair, err := h.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (h *AuthHandlers) Me(c *gin.Context) {
	// Get user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Find user in database
	usersCollection := h.db.Collection("users")
	var user models.User
	err := usersCollection.FindOne(context.Background(), bson.M{
		"_id": userCtx.UserID,
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userInfo := UserInfo{
		ID:             user.ID.Hex(),
		Username:       user.Username,
		Email:          user.Email,
		RoleID:         user.RoleID.Hex(),
		IsInitialLogin: user.IsInitialLogin,
	}

	c.JSON(http.StatusOK, userInfo)
}
