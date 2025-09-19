package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"orcaai/backend/database"
	"orcaai/backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch user from database
	var user models.User
	query := "SELECT id, email, name, role, is_active, created_at, updated_at FROM users WHERE id = $1"
	row := database.DB.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch user profile",
			"details": err.Error(),
		})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, gin.H{
		"user": models.UserProfile{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	})
}

// UpdateProfile updates the authenticated user's profile
func UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Prepare update query
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	// Add name to update if provided
	if req.Name != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}

	// Add email to update if provided
	if req.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}

	// Add password to update if provided
	if req.Password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		updates = append(updates, fmt.Sprintf("password_hash = $%d", argIndex))
		args = append(args, string(hashedPassword))
		argIndex++
	}

	// Always update the updated_at timestamp
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add user ID as the last argument
	args = append(args, userID)

	// Construct the query
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		argIndex,
	)

	// Execute the update
	_, err := database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user profile",
			"details": err.Error(),
		})
		return
	}

	// Fetch updated user
	var updatedUser models.User
	query = "SELECT id, email, name, role, is_active, created_at, updated_at FROM users WHERE id = $1"
	row2 := database.DB.QueryRow(query, userID)
	err = row2.Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.Name, &updatedUser.Role, &updatedUser.IsActive, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch updated user profile",
			"details": err.Error(),
		})
		return
	}

	// Return updated user profile
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": models.UserProfile{
			ID:        updatedUser.ID,
			Email:     updatedUser.Email,
			Name:      updatedUser.Name,
			Role:      updatedUser.Role,
			IsActive:  updatedUser.IsActive,
			CreatedAt: updatedUser.CreatedAt,
		},
	})
}
