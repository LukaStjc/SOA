package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-userm/graphdb"
	"go-userm/initializers"
	"go-userm/models"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
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

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email cannot be empty",
		})
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	if body.Role.String() == "Unknown" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Role is not correct or can not be empty",
		})
		return
	}

	if body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password cannot be empty",
		})
		return
	}

	if body.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username cannot be empty",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	tx := initializers.DB.Begin()

	user := models.User{
		Username: body.Username,
		Password: string(hash),
		Email:    body.Email,
		Role:     body.Role,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user. The log is: " + err.Error(),
		})
		return
	}

	graphUser := models.GraphDBUser{
		ID:       int64(user.ID),
		Username: user.Username,
	}

	if err := graphdb.WriteUser(&graphUser, initializers.Neo4JDriver); err != nil {
		log.Printf("Error creating user node for %s in Neo4j database: %v", user.Username, err)
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user. The log is: " + err.Error(),
		})
		return
	}

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

	resp, err := http.Post("http://auth-service:3001/signup", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to send request",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create user. Status code: %d", resp.StatusCode),
		})
		fmt.Println(resp)
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {

	// ako hoces da dobavis neko polje usera
	// onda zapocinjes komandu na sledeci nacin
	// user.(models.User).

	userInterface, _ := c.Get("user")

	// NE SME *models.User!
	user, _ := userInterface.(models.User)

	if user.Role.String() == "Administrator" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Administrator Content.",
		})
	} else if user.Role.String() == "Guide" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Guide Content.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tourist Content.",
		})
	}
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not founddddd"})
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
	followeeUsername := c.Param("username")

	var followeeUser models.User
	result := initializers.DB.Where("username = ?", followeeUsername).First(&followeeUser)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if followeeUser.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Cannot follow an admin!"})
		return
	}

	if user.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are an admin, you cant follow anybody!"})
		return
	}

	if user.ID == followeeUser.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Cannot follow yourself, there will be a profile when frontend is done!"})
		return
	}

	// Vasilije:
	// Ovde se vise ne radi automatski unfollow ako korisnik ponovo pokusa da zaprati osobu koju vec prati.
	// To se sada radi u drugom kontroleru.
	// Takodje, ta odluka ne smeta frontu, jer se na frontu dugme follow disable-uje
	// ako se korisnik vec prati, a to se sazna preko endpointa GetFollowed.

	err := graphdb.FollowUser(user.Username, followeeUsername, initializers.Neo4JDriver)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

func IsBlocked(c *gin.Context) {
	fmt.Printf("Usao u is blocked")
	// Immediately return if the previous middleware aborted the request
	if c.IsAborted() {
		return
	}
	// Extracting the user id from the path

	userId := c.Param("id")
	fmt.Printf("Zabo u atoiu")
	userIdAsNum, _ := strconv.Atoi(userId)

	fmt.Printf("Id usera je %d", userIdAsNum)

	// Find the user by id
	var user models.User
	result := initializers.DB.Where("id = ?", userId).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not foundddd"})
		return
	}

	// Check if user is blocked
	if user.Blocked {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is blocked"})
		return
	}

	// If user is not blocked send ok status
	c.JSON(http.StatusOK, gin.H{"message": "User is not blocked"})

}

func DoesFollow(c *gin.Context) {
	if c.IsAborted() {
		return
	}

	followerId := c.Param("followerId")
	creatorId := c.Param("creatorId")

	follower, err := graphdb.FindUserByID(followerId, initializers.Neo4JDriver)
	if err != nil || follower == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User (follower) not found"})
		return
	}
	log.Println(follower)

	creator, err := graphdb.FindUserByID(creatorId, initializers.Neo4JDriver)
	if err != nil || creator == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User (creator) not found"})
		return
	}
	log.Println(creator)

	doesFollow, err := graphdb.DoesFollow(followerId, creatorId, initializers.Neo4JDriver)

	log.Println("does follow")
	log.Println(doesFollow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Database exception"})
		return
	}
	if doesFollow {
		c.JSON(http.StatusOK, gin.H{"follows": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"follows": false})
}

func GetFriendsRecommendation(c *gin.Context) {
	authUser, _ := c.Get("user")
	user := authUser.(models.User)

	if user.Role.String() == "Administrator" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are the admin. You can't have friends."})
		return
	}

	recommendations, err := graphdb.RecommendFriends(user.ID, initializers.Neo4JDriver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get friend recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

func GetById(c *gin.Context) {

	// Find the user by id
	id := c.Param("id")
	var user models.User
	result := initializers.DB.First(&user, id)
	if result.Error != nil {
		// If user not found, return a 404 Not Found response
		c.JSON(http.StatusNotFound, gin.H{"error": "User not founddasdd"})
		return
	}
	// If user found, return user data in JSON format
	c.JSON(http.StatusOK, user)

}

func GetFollowed(c *gin.Context) {
	if c.IsAborted() {
		return
	}
	auth_user, _ := c.Get("user")

	// NE SME *models.User!
	user, _ := auth_user.(models.User)

	followedUserIDs, err := graphdb.GetFollowees(user.ID, initializers.Neo4JDriver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve followed users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followed_user_ids": followedUserIDs})
}
