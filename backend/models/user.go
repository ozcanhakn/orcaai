package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Role definitions
const (
	RoleAdmin      = "admin"
	RoleUser       = "user"
	RoleEnterprise = "enterprise"
)

// Permission definitions
const (
	PermissionReadMetrics        = "read:metrics"
	PermissionWriteAPIKeys       = "write:api_keys"
	PermissionReadAPIKeys        = "read:api_keys"
	PermissionAdminUsers         = "admin:users"
	PermissionEnterpriseFeatures = "enterprise:features"
)

type TokenUsage struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}

type AIResponse struct {
	Content    string                 `json:"content"`
	Provider   string                 `json:"provider"`
	Model      string                 `json:"model"`
	TokensUsed TokenUsage             `json:"tokens_used"`
	Cost       float64                `json:"cost"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// APIKey represents an API key for a user
type APIKey struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Key         string    `json:"key" db:"key"`
	Permissions []string  `json:"permissions" db:"permissions"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RolePermissions maps roles to their permissions
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermissionReadMetrics,
		PermissionWriteAPIKeys,
		PermissionReadAPIKeys,
		PermissionAdminUsers,
		PermissionEnterpriseFeatures,
	},
	RoleUser: {
		PermissionReadMetrics,
		PermissionWriteAPIKeys,
		PermissionReadAPIKeys,
	},
	RoleEnterprise: {
		PermissionReadMetrics,
		PermissionWriteAPIKeys,
		PermissionReadAPIKeys,
		PermissionEnterpriseFeatures,
	},
}
