package middleware

import (
	"go-tourm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckIfAdmin(c *gin.Context) {

	// Attempt to retrieve the user from the context, which was presumably set by RequireAuth.
	userInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - User not found in context"})
		return
	}

	// Assert the type of the user to your user model.
	user, ok := userInterface.(models.User) // Ensure this type assertion aligns with how you're setting the user in RequireAuth.
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error - User type assertion failed"})
		return
	}

	// Check if the user's role is "admin".
	if user.Role.String() != "Administrator" { // Adjust the field and value based on your user model and role values.
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - User is not an admin"})
		return
	}

	// If the user is an admin, continue processing the request.
	c.Next()
}
