package middleware

import (
	"net/http"
	"strings"

	"orcaai/models"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware checks if the user has the required permissions
func RBACMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by auth middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		role := userRole.(string)

		// Get permissions for the role
		rolePermissions, exists := models.RolePermissions[role]
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
			c.Abort()
			return
		}

		// Check if user has all required permissions
		if !hasPermissions(rolePermissions, requiredPermissions) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermissions checks if the user has all required permissions
func hasPermissions(userPermissions, requiredPermissions []string) bool {
	// If no specific permissions required, allow access
	if len(requiredPermissions) == 0 {
		return true
	}

	// Create a map for faster lookup
	permissionMap := make(map[string]bool)
	for _, perm := range userPermissions {
		permissionMap[perm] = true
	}

	// Check if user has all required permissions
	for _, requiredPerm := range requiredPermissions {
		if !permissionMap[requiredPerm] {
			return false
		}
	}

	return true
}

// RoleMiddleware checks if the user has one of the allowed roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by auth middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		role := userRole.(string)

		// Check if user's role is in the allowed roles
		if !isRoleAllowed(role, allowedRoles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isRoleAllowed checks if the user's role is in the allowed roles
func isRoleAllowed(userRole string, allowedRoles []string) bool {
	// If no specific roles required, allow access
	if len(allowedRoles) == 0 {
		return true
	}

	// Check if user's role is in the allowed roles
	for _, allowedRole := range allowedRoles {
		if strings.EqualFold(userRole, allowedRole) {
			return true
		}
	}

	return false
}
