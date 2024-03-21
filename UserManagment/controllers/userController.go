package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-userm/initializers"
	"go-userm/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// TODO: validations
func SignUp(c *gin.Context) {
	// Get the username/pass/email/role of req body
	var body struct {
		Username string
		Password string
		Email    string
		Role     models.UserRole
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Begin a transaction
	tx := initializers.DB.Begin()

	// Create the user
	user := models.User{
		Username: body.Username,
		Password: string(hash),
		Email:    body.Email,
		Role:     body.Role,
	}

	if err := tx.Create(&user).Error; err != nil {
		// If an error occurs, rollback the transaction
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"ID":       user.ID,
		"Username": body.Username,
		"Password": string(hash),
		"Role":     body.Role,
	})
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to marshal request body",
		})
		return
	}
	fmt.Println("JSON Spreman: ", body.Role.String())

	// Send the request
	resp, err := http.Post("http://localhost:3001/signup", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to send request",
		})
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create user. Status code: %d", resp.StatusCode),
		})
		fmt.Println(resp)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	// ako hoces da dobavis neko polje usera
	// onda zapocinjes komandu na sledeci nacin
	// user.(models.User).

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func BlockUser(c *gin.Context) {
	// Immediately return if the previous middleware aborted the request
	if c.IsAborted() {
		return
	}
	// Extracting the username from the path

	username := c.Param("username")

	// Find the user by username
	var user models.User
	result := initializers.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if user.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Can't block an admin"})
		return
	}

	// Check if user is already blocked
	if user.Blocked {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is already blocked"})
		return
	}

	// Update the Blocked field
	user.Blocked = true
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})

}

func Follow(c *gin.Context) {
	authUser, _ := c.Get("user")

	user := authUser.(models.User)
	username := c.Param("username")

	// Find the user that I want to follow by username
	var newUser models.User
	result := initializers.DB.Where("username = ?", username).First(&newUser)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if newUser.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Cannot follow an admin!"})
		return
	}

	if user.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are an admin, you cant follow anybody!"})
		return
	}

	if user.ID == newUser.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Cannot follow yourself, there will be a profile when frontend is done!"})
		return
	}

	for _, u := range user.Follows {
		if u.ID == newUser.ID {
			initializers.DB.Model(&user).Association("follow_id").Delete(&newUser)
			c.JSON(http.StatusOK, gin.H{"message": "User unfollowed"})
			return
		}
	}

	user.Follows = append(user.Follows, &newUser)
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})

}
