package middleware

import (
	"encoding/json"
	"fmt"
	"go-userm/initializers"
	"go-userm/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	port := os.Getenv("AUTHENTICATION_SERVICE_PORT")
	url := fmt.Sprintf("http://auth-service:%s/authenticate", port)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating authentication request"})
		return
	}

	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error sending authentication request"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respError map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respError); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error decoding authentication service response"})
		}

		if errorMsg, ok := respError["error"].(string); ok && errorMsg != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		}

		return
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error decoding authentication response"})
		return
	}

	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not foundddd"})
		return
	}

	// iz auth mikroservisa imam samo atributa, ovde inicijalizujem sve ostale
	var fullyInitializedUser models.User
	result := initializers.DB.First(&fullyInitializedUser, "id = ?", user.ID)
	if result.Error != nil || fullyInitializedUser.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found in database"})
		return
	}

	c.Set("user", fullyInitializedUser)

	fmt.Println(fullyInitializedUser.Email)

	c.Next()
}

// func RequireAuth(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")

// 	if authHeader == "" {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
// 		return
// 	}

// 	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

// 	fmt.Println(tokenString)

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(os.Getenv("SECRET")), nil
// 	})

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// check the exp
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
// 		}

// 		// Find the user with token sub
// 		var user models.User
// 		initializers.DB.First(&user, claims["sub"])

// 		if user.ID == 0 {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is not found"})
// 		}

// 		// Attach to req
// 		c.Set("user", user)

// 		// Continue
// 		c.Next()

// 	} else {
// 		fmt.Println(err)
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims aren't okay"})
// 	}

// 	c.Next()
// }
