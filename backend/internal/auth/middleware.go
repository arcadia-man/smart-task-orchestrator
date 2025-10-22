package auth

import (
    "context"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserContext struct {
    UserID   primitive.ObjectID
    Username string
    RoleID   primitive.ObjectID
}

const UserContextKey = "user"

func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Check Bearer prefix
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        // Validate token
        claims, err := jwtManager.ValidateToken(tokenParts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Set user context
        userCtx := &UserContext{
            UserID:   claims.UserID,
            Username: claims.Username,
            RoleID:   claims.RoleID,
        }

        c.Set(UserContextKey, userCtx)
        c.Next()
    }
}

func GetUserFromContext(c *gin.Context) (*UserContext, bool) {
    user, exists := c.Get(UserContextKey)
    if !exists {
        return nil, false
    }

    userCtx, ok := user.(*UserContext)
    return userCtx, ok
}

func GetUserFromGoContext(ctx context.Context) (*UserContext, bool) {
    user := ctx.Value(UserContextKey)
    if user == nil {
        return nil, false
    }

    userCtx, ok := user.(*UserContext)
    return userCtx, ok
}

// Optional middleware for endpoints that don't require auth but can use it
func OptionalAuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.Next()
            return
        }

        claims, err := jwtManager.ValidateToken(tokenParts[1])
        if err != nil {
            c.Next()
            return
        }

        userCtx := &UserContext{
            UserID:   claims.UserID,
            Username: claims.Username,
            RoleID:   claims.RoleID,
        }

        c.Set(UserContextKey, userCtx)
        c.Next()
    }
}
