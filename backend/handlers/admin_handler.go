package handlers

import (
	"net/http"
    "os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
    "orcaai/backend/database"
    "orcaai/backend/utils"
)

// GetAllUsers returns a list of all users (admin only)
func GetAllUsers(c *gin.Context) {
	// In a real implementation, this would fetch users from the database
	users := []map[string]interface{}{
		{
			"id":    uuid.New(),
			"email": "admin@example.com",
			"name":  "Admin User",
			"role":  "admin",
		},
		{
			"id":    uuid.New(),
			"email": "user@example.com",
			"name":  "Regular User",
			"role":  "user",
		},
		{
			"id":    uuid.New(),
			"email": "enterprise@example.com",
			"name":  "Enterprise User",
			"role":  "enterprise",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// UpdateUser updates a user's information (admin only)
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	// In a real implementation, this would update the user in the database
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user_id": userID,
	})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// In a real implementation, this would delete the user from the database
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

// GetDetailedMetrics returns detailed metrics (admin only)
func GetDetailedMetrics(c *gin.Context) {
	// In a real implementation, this would fetch detailed metrics
	c.JSON(http.StatusOK, gin.H{
		"message": "Detailed metrics endpoint",
		"data":    "Detailed metrics data would be returned here",
	})
}

// SetProviderKey encrypts and stores provider API key (admin only)
func SetProviderKey(c *gin.Context) {
    var req struct{ Provider string `json:"provider" binding:"required"`; ApiKey string `json:"api_key" binding:"required"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    secret := os.Getenv("PROVIDER_SECRET_KEY")
    if secret == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "secret key not configured"})
        return
    }
    enc, err := utils.EncryptAESGCM([]byte(req.ApiKey), secret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
        return
    }
    _, err = database.DB.Exec("UPDATE ai_providers SET api_key_encrypted = $1 WHERE name = $2", enc, req.Provider)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "provider key updated"})
}

// GetProviderKeyStatus returns if provider key is configured (no secret returned)
func GetProviderKeyStatus(c *gin.Context) {
    provider := c.Query("provider")
    if provider == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
        return
    }
    var exists bool
    err := database.DB.QueryRow("SELECT api_key_encrypted IS NOT NULL FROM ai_providers WHERE name = $1", provider).Scan(&exists)
    if err != nil { exists = false }
    c.JSON(http.StatusOK, gin.H{"configured": exists})
}
