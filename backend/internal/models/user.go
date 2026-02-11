package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID             primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
    Username       string                 `json:"username" bson:"username"`
    Email          string                 `json:"email" bson:"email"`
    PasswordHash   string                 `json:"-" bson:"password_hash"`
    RoleID         primitive.ObjectID     `json:"roleId" bson:"role_id"`
    IsInitialLogin bool                   `json:"isInitialLogin" bson:"is_initial_login"`
    Active         bool                   `json:"active" bson:"active"`
    LastLoginAt    *time.Time             `json:"lastLoginAt" bson:"last_login_at"`
    CreatedAt      time.Time              `json:"createdAt" bson:"created_at"`
    UpdatedAt      time.Time              `json:"updatedAt" bson:"updated_at"`
    CreatedBy      *primitive.ObjectID    `json:"createdBy" bson:"created_by"`
    UpdatedBy      *primitive.ObjectID    `json:"updatedBy" bson:"updated_by"`
    Metadata       map[string]interface{} `json:"metadata" bson:"metadata"`
}

type Role struct {
    ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    RoleName      string              `json:"roleName" bson:"role_name"`
    Description   string              `json:"description" bson:"description"`
    Active        bool                `json:"active" bson:"active"`
    CanCreateTask bool                `json:"canCreateTask" bson:"can_create_task"`
    CreatedAt     time.Time           `json:"createdAt" bson:"created_at"`
    UpdatedAt     time.Time           `json:"updatedAt" bson:"updated_at"`
    CreatedBy     *primitive.ObjectID `json:"createdBy" bson:"created_by"`
}

type Permission struct {
    ID              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    RoleID          primitive.ObjectID  `json:"roleId" bson:"role_id"`
    SchedulerID     primitive.ObjectID  `json:"schedulerId" bson:"scheduler_id"`
    RoleSee         bool                `json:"roleSee" bson:"role_see"`
    RoleExecute     bool                `json:"roleExecute" bson:"role_execute"`
    RoleAlterConfig bool                `json:"roleAlterConfig" bson:"role_alter_config"`
    CreatedAt       time.Time           `json:"createdAt" bson:"created_at"`
    CreatedBy       *primitive.ObjectID `json:"createdBy" bson:"created_by"`
}

type Image struct {
    ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    RegistryURL string              `json:"registryUrl" bson:"registry_url"`
    Image       string              `json:"image" bson:"image"`
    Name        string              `json:"name" bson:"name"`
    Description string              `json:"description" bson:"description"`
    Version     string              `json:"version" bson:"version"`
    IsDefault   bool                `json:"isDefault" bson:"is_default"`
    CreatedAt   time.Time           `json:"createdAt" bson:"created_at"`
    CreatedBy   *primitive.ObjectID `json:"createdBy" bson:"created_by"`
}
