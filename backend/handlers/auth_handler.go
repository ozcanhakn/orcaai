package handlers

import (
	"net/http"

	"orcaai/database"
	"orcaai/middleware"
	"orcaai/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account
func Register(c *gin.Context) {
	var registration models.UserRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	query := "SELECT id FROM users WHERE email = $1"
	row := database.DB.QueryRow(query, registration.Email)
	err := row.Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user (default to 'user' role)
	userID := uuid.New()
	user := models.User{
		ID:           userID,
		Email:        registration.Email,
		PasswordHash: string(hashedPassword),
		Name:         registration.Name,
		Role:         models.RoleUser, // Default role
		IsActive:     true,
	}

	// Insert user into database
	query = `
		INSERT INTO users (id, email, password_hash, name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, err = database.DB.Exec(query, user.ID, user.Email, user.PasswordHash, user.Name, user.Role, user.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, user.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return success response with token
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": models.UserProfile{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		"token": token,
	})
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	var login models.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	query := "SELECT id, email, password_hash, name, role, is_active, created_at, updated_at FROM users WHERE email = $1"
	row := database.DB.QueryRow(query, login.Email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is deactivated"})
		return
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, user.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return success response with token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": models.UserProfile{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		"token": token,
	})
}

// RefreshToken generates a new JWT token for an authenticated user
func RefreshToken(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEmail, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	isActive, exists := c.Get("is_active")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate new JWT token
	token, err := middleware.RefreshToken(
		userID.(uuid.UUID),
		userEmail.(string),
		userRole.(string),
		isActive.(bool),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return success response with new token
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token":   token,
	})
}
