package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"smart-task-orchestrator/internal/models"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}

	// Initialize database
	if err := db.Initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Printf("✅ Connected to MongoDB: %s", dbName)
	return db, nil
}

func (db *MongoDB) Initialize(ctx context.Context) error {
	// Create indexes
	if err := db.createIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Seed initial data
	if err := db.seedInitialData(ctx); err != nil {
		return fmt.Errorf("failed to seed initial data: %w", err)
	}

	return nil
}

func (db *MongoDB) createIndexes(ctx context.Context) error {
	indexes := []struct {
		collection string
		index      mongo.IndexModel
	}{
		// Users
		{
			collection: "users",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			collection: "users",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetSparse(true),
			},
		},
		// Roles
		{
			collection: "roles",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "role_name", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// Permissions
		{
			collection: "permissions",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "role_id", Value: 1}, {Key: "scheduler_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// Scheduler Precompute - Critical for performance
		{
			collection: "scheduler_precompute",
			index: mongo.IndexModel{
				Keys: bson.D{{Key: "run_at", Value: 1}, {Key: "status", Value: 1}},
			},
		},
		{
			collection: "scheduler_precompute",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "scheduler_id", Value: 1}, {Key: "run_at", Value: 1}, {Key: "generation", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// Scheduler History
		{
			collection: "scheduler_history",
			index: mongo.IndexModel{
				Keys: bson.D{{Key: "scheduler_id", Value: 1}, {Key: "start_time", Value: -1}},
			},
		},
		{
			collection: "scheduler_history",
			index: mongo.IndexModel{
				Keys:    bson.D{{Key: "run_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	}

	for _, idx := range indexes {
		collection := db.Database.Collection(idx.collection)
		_, err := collection.Indexes().CreateOne(ctx, idx.index)
		if err != nil {
			return fmt.Errorf("failed to create index on %s: %w", idx.collection, err)
		}
	}

	log.Println("✅ Database indexes created")
	return nil
}

func (db *MongoDB) seedInitialData(ctx context.Context) error {
	usersCollection := db.Database.Collection("users")
	rolesCollection := db.Database.Collection("roles")

	// Check if admin user exists
	var existingAdmin models.User
	err := usersCollection.FindOne(ctx, bson.M{"username": "admin"}).Decode(&existingAdmin)

	if err == nil {
		// Admin user already exists, skip seeding
		log.Println("✅ Admin user already exists, skipping seed")
		return nil
	}

	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to check for admin user: %w", err)
	}

	// Admin user doesn't exist, create it
	log.Println("🌱 Creating initial admin user...")

	// Create or get admin role
	var adminRole models.Role
	err = rolesCollection.FindOne(ctx, bson.M{"role_name": "admin"}).Decode(&adminRole)

	if err == mongo.ErrNoDocuments {
		// Create admin role
		adminRole = models.Role{
			RoleName:      "admin",
			Description:   "System Administrator",
			Active:        true,
			CanCreateTask: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		result, err := rolesCollection.InsertOne(ctx, adminRole)
		if err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}
		adminRole.ID = result.InsertedID.(primitive.ObjectID)
	} else if err != nil {
		return fmt.Errorf("failed to check for admin role: %w", err)
	}

	// Create admin user with password "12345"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	adminUser := models.User{
		Username:       "admin",
		Email:          "admin@orchestrator.local",
		PasswordHash:   string(passwordHash),
		RoleID:         adminRole.ID,
		IsInitialLogin: true,
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	_, err = usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("✅ Initial admin user created - Username: admin, Password: 12345")
	return nil
}

func (db *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Client.Disconnect(ctx)
}
