package controllers

// import (
// 	"fmt"
// 	"go-jwt/initializers"
// 	"go-jwt/models"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// 	"golang.org/x/crypto/bcrypt"
// )

import (
	"context"
	"encoding/json"
	"fmt"
	"go-jwt/initializers"
	"go-jwt/models"
	auth "go-jwt/proto/generatedFiles"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	NATSClient *nats.Conn
}

type SignUpSuccessEvent struct {
	UserID uint32 `json:"userId"`
}

type SignUpFailedEvent struct {
	UserID uint32 `json:"userId"`
	Error  string `json:"error"`
}

// TODO: Maybe it is unnecessary, but I can cover transaction management and in auth ms.
// => Potpuno nepotrebno. Transakcija nema ovde posla, jer moze da se ode u dva pravca i to je sve.
func PublishSignUpSuccess(natsClient *nats.Conn, userID uint32) {
	successEvent := SignUpSuccessEvent{
		UserID: userID,
	}
	fmt.Println("Usao u publish sign up success iz auth ms")
	eventData, err := json.Marshal(successEvent)
	if err != nil {
		// TODO: What to do with transactions? And generally what to do in this case?
		// => This wouldn't happen. This is a server side logic.
		log.Printf("Error marshalling SignUpSuccess event: %v", err)
		return
	}
	if err := natsClient.Publish("SignUpSuccess", eventData); err != nil {
		// TODO: What to do with transactions? And generally what to do in this case?
		// => This wouldn't happen. This is a server side logic.
		log.Printf("Failed to publish SignUpSuccess event: %v", err)
	}
	fmt.Println("Satro je uspeo da odradio publish signup success dogadjaja iz auth ms")
}

func PublishSignUpFailed(natsClient *nats.Conn, userID uint32, errorMsg string) {
	failedEvent := SignUpFailedEvent{
		UserID: userID,
		Error:  errorMsg,
	}
	eventData, err := json.Marshal(failedEvent)
	if err != nil {
		// TODO: What to do with transactions? And generally what to do in this case?
		// => This wouldn't happen. This is a server side logic.
		log.Printf("Error marshalling SignUpFailed event: %v", err)
		return
	}
	if err := natsClient.Publish("SignUpFailed", eventData); err != nil {
		// TODO: What to do with transactions? And generally what to do in this case?
		// => This wouldn't happen. This is a server side logic.
		log.Printf("Failed to publish SignUpFailed event: %v", err)
	}
}

func (s *AuthHandler) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.SignUpResponse, error) {

	newUser := models.User{
		ID:       uint(req.Id),
		Username: req.Username,
		Password: req.Password,
		Role:     uint8(req.Role),
	}

	fmt.Println("Ispis iz auth signupa ", newUser)

	result := initializers.DB.Create(&newUser)
	if result.Error != nil {
		PublishSignUpFailed(s.NATSClient, uint32(newUser.ID), "Sign up process failed in authentication microservice!")
		return &auth.SignUpResponse{Status: "Failed"}, result.Error
	}

	fmt.Println("Stiglo do pred objavljivanje dogadjana o uspesnosti iz auth ms")

	PublishSignUpSuccess(s.NATSClient, uint32(newUser.ID))

	return &auth.SignUpResponse{Status: "Success"}, nil
}

// func SignUp(c *gin.Context) {
// 	var body s    uint8
// 	}truct {
// 		ID       uint
// 		Username string
// 		Password string
// 		Role

// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to read body",
// 		})

// 		return
// 	}

// 	user := models.User{ID: body.ID, Username: body.Username, Password: body.Password, Role: uint8(body.Role)}
// 	result := initializers.DB.Create(&user)

// 	if result.Error != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create user",
// 		})

// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{})
// }

func (s *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user models.User
	result := initializers.DB.Where("username = ?", req.Username).First(&user)

	if result.Error != nil || user.ID == 0 {
		return nil, fmt.Errorf("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to create token")
	}

	return &auth.LoginResponse{
		Username: user.Username,
		Token:    tokenString,
		Id:       uint32(user.ID),
		Role:     uint32(user.Role),
	}, nil
}

// func Login(c *gin.Context) {
// 	var body struct {
// 		Username string
// 		Password string
// 	}

// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to read body",
// 		})

// 		return
// 	}

// 	var user models.User
// 	initializers.DB.First(&user, "username = ?", body.Username)

// 	if user.ID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid username or password",
// 		})

// 		return
// 	}

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid username or password",
// 		})

// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub":  user.ID, // subject
// 		"role": user.Role,
// 		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
// 	})

// 	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create token",
// 		})

// 		return
// 	}

// 	// Send it back
// 	c.JSON(http.StatusOK, gin.H{
// 		"username": user.Username,
// 		"token":    tokenString,
// 		"id":       user.ID,
// 		"role":     user.Role,
// 	})
// }

func (s *AuthHandler) Authenticate(ctx context.Context, req *emptypb.Empty) (*auth.AuthenticateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to extract metadata")
	}

	tokens, ok := md["authorization"]
	if !ok || len(tokens) < 1 {
		return nil, status.Error(codes.Unauthenticated, "Authorization token not found")
	}

	tokenString := strings.TrimPrefix(tokens[0], "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("SECRET")
		if secret == "" {
			return nil, fmt.Errorf("secret key not set")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["sub"].(float64)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Invalid token claims")
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", uint(id))
		if result.Error != nil || user.ID == 0 {
			return nil, status.Error(codes.NotFound, "User not found")
		}

		if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
			return nil, status.Error(codes.Unauthenticated, "Token is expired")
		}

		return &auth.AuthenticateResponse{
			Id:       uint32(user.ID),
			Username: user.Username,
			Password: user.Password,
			Role:     uint32(user.Role),
		}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Invalid token")
}

// func Authenticate(c *gin.Context) {
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
// 		return
// 	}

// 	tokenString = tokenString[len("Bearer "):]

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		secret := os.Getenv("SECRET")
// 		if secret == "" {
// 			return nil, fmt.Errorf("secret key not set")
// 		}
// 		return []byte(secret), nil
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		sub, ok := claims["sub"].(float64)
// 		fmt.Println("token", tokenString)
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 			return
// 		}

// 		var user models.User
// 		result := initializers.DB.First(&user, "id = ?", sub)
// 		if result.Error != nil || user.ID == 0 {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			return
// 		}

// 		if exp, ok := claims["exp"].(float64); ok {
// 			if time.Now().Unix() > int64(exp) {
// 				c.JSON(http.StatusUnauthorized, gin.H{
// 					"error": "Token is expired"})
// 				return
// 			}
// 		}

// 		c.JSON(http.StatusOK, user)
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 	}
// }
