package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/models"
)

type ImageHandlers struct {
	db *mongo.Database
}

func NewImageHandlers(db *mongo.Database) *ImageHandlers {
	return &ImageHandlers{db: db}
}

type CreateImageRequest struct {
	RegistryURL string `json:"registryUrl"`
	Image       string `json:"image" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Version     string `json:"version"`
	IsDefault   bool   `json:"isDefault"`
}

type UpdateImageRequest struct {
	RegistryURL string `json:"registryUrl"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	IsDefault   *bool  `json:"isDefault"`
}

type ImageResponse struct {
	ID          string    `json:"id"`
	RegistryURL string    `json:"registryUrl"`
	Image       string    `json:"image"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	IsDefault   bool      `json:"isDefault"`
	FullName    string    `json:"fullName"`
	UsageCount  int64     `json:"usageCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (h *ImageHandlers) GetImages(c *gin.Context) {
	imagesCollection := h.db.Collection("images")
	schedulersCollection := h.db.Collection("scheduler_definitions")

	// Get all images
	cursor, err := imagesCollection.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}
	defer cursor.Close(context.Background())

	var images []ImageResponse
	for cursor.Next(context.Background()) {
		var image models.Image
		if err := cursor.Decode(&image); err != nil {
			continue
		}

		// Count usage in schedulers
		usageCount, _ := schedulersCollection.CountDocuments(context.Background(), bson.M{"image": image.Image})

		// Build full name
		fullName := image.Image
		if image.RegistryURL != "" {
			fullName = image.RegistryURL + "/" + image.Image
		}
		if image.Version != "" {
			fullName += ":" + image.Version
		}

		images = append(images, ImageResponse{
			ID:          image.ID.Hex(),
			RegistryURL: image.RegistryURL,
			Image:       image.Image,
			Name:        image.Name,
			Description: image.Description,
			Version:     image.Version,
			IsDefault:   image.IsDefault,
			FullName:    fullName,
			UsageCount:  usageCount,
			CreatedAt:   image.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, images)
}

func (h *ImageHandlers) CreateImage(c *gin.Context) {
	var req CreateImageRequest
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

	// Check if image already exists
	imagesCollection := h.db.Collection("images")
	var existingImage models.Image
	err := imagesCollection.FindOne(context.Background(), bson.M{
		"image":   req.Image,
		"version": req.Version,
	}).Decode(&existingImage)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Image with this version already exists"})
		return
	}

	// If this is set as default, unset other defaults
	if req.IsDefault {
		imagesCollection.UpdateMany(
			context.Background(),
			bson.M{"is_default": true},
			bson.M{"$set": bson.M{"is_default": false}},
		)
	}

	// Create image
	now := time.Now()
	image := models.Image{
		RegistryURL: req.RegistryURL,
		Image:       req.Image,
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		IsDefault:   req.IsDefault,
		CreatedAt:   now,
		CreatedBy:   &userCtx.UserID,
	}

	result, err := imagesCollection.InsertOne(context.Background(), image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image"})
		return
	}

	image.ID = result.InsertedID.(primitive.ObjectID)

	// Build full name
	fullName := image.Image
	if image.RegistryURL != "" {
		fullName = image.RegistryURL + "/" + image.Image
	}
	if image.Version != "" {
		fullName += ":" + image.Version
	}

	c.JSON(http.StatusCreated, ImageResponse{
		ID:          image.ID.Hex(),
		RegistryURL: image.RegistryURL,
		Image:       image.Image,
		Name:        image.Name,
		Description: image.Description,
		Version:     image.Version,
		IsDefault:   image.IsDefault,
		FullName:    fullName,
		UsageCount:  0,
		CreatedAt:   image.CreatedAt,
	})
}

func (h *ImageHandlers) UpdateImage(c *gin.Context) {
	imageID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	var req UpdateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	imagesCollection := h.db.Collection("images")

	// Build update document
	updateDoc := bson.M{}

	if req.RegistryURL != "" {
		updateDoc["registry_url"] = req.RegistryURL
	}
	if req.Image != "" {
		updateDoc["image"] = req.Image
	}
	if req.Name != "" {
		updateDoc["name"] = req.Name
	}
	if req.Description != "" {
		updateDoc["description"] = req.Description
	}
	if req.Version != "" {
		updateDoc["version"] = req.Version
	}
	if req.IsDefault != nil {
		updateDoc["is_default"] = *req.IsDefault

		// If setting as default, unset other defaults
		if *req.IsDefault {
			imagesCollection.UpdateMany(
				context.Background(),
				bson.M{"_id": bson.M{"$ne": objectID}},
				bson.M{"$set": bson.M{"is_default": false}},
			)
		}
	}

	if len(updateDoc) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	result, err := imagesCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateDoc},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully"})
}

func (h *ImageHandlers) DeleteImage(c *gin.Context) {
	imageID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	imagesCollection := h.db.Collection("images")
	schedulersCollection := h.db.Collection("scheduler_definitions")

	// Get the image to check its name
	var image models.Image
	err = imagesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}

	// Check if image is being used by any schedulers
	usageCount, err := schedulersCollection.CountDocuments(context.Background(), bson.M{"image": image.Image})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check image usage"})
		return
	}

	if usageCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete image that is being used by schedulers"})
		return
	}

	result, err := imagesCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
