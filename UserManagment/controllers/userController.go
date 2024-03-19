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
