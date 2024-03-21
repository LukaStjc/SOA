package middleware

import (
	"fmt"
	"go-userm/initializers"
	"go-userm/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckIfAdmin(c *gin.Context) {
	// Call RequireAuth middleware to ensure authentication is handled.
	// However, this should ideally be setup as a middleware for the route,
	// not called directly within this function. This example retains it for consistency with your approach.
	//RequireAuth(c)
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	fmt.Println(tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is not found"})
		}

		// Attach to req
		c.Set("user", user)

	} else {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims aren't okay"})
	}

	// Attempt to retrieve the user from the context.
	userInterface, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - User not found in context"})
		return
	}

	// Assert the type of the user to models.User.
	user, ok := userInterface.(models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error - User type assertion failed"})
		return
	}

	// Check if the user's role is "admin".

	if user.Role.String() != "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - User is not an admin"})
		return
	}

	// If the user is an admin, continue processing the request.

	c.Next()
}
