package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-userm/graphdb"
	"go-userm/initializers"
	"go-userm/models"
	user "go-userm/proto/user/generatedFiles"
	"go-userm/transactionmanager"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	AuthServiceClient  user.AuthServiceClient
	NATSClient         *nats.Conn // I've added it, buy why?
	TransactionManager *transactionmanager.Manager
	Neo4jDriver        neo4j.DriverWithContext
	Neo4jSession       neo4j.SessionWithContext
}

type SignUpSuccessEvent struct {
	UserID uint32 `json:"userId"`
}

type SignUpFailedEvent struct {
	UserID uint32 `json:"userId"`
	Error  string `json:"error"`
}

func (h *UserHandler) HandleSignUpSuccess() {
	fmt.Println("usli u sign up sucess u user management ms")
	_, err := h.NATSClient.Subscribe("SignUpSuccess", func(m *nats.Msg) {
		fmt.Println("Okinuo se sign up sucess u user management ms")
		var event SignUpSuccessEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Error unmarshalling SignUpSuccess event: %v", err)
			return
		}

		log.Printf("Received SignUpSuccess for user ID %d", event.UserID)
		// Here, ideally, you would update the status of the user to "active" or similar
		// Example: updateUserStatus(event.UserID, "active")

		h.TransactionManager.Commit(event.UserID, context.Background())
	})

	if err != nil {
		// TODO: What about transactions handling?
		log.Fatalf("Failed to subscribe to SignUpSuccess events: %v", err)
	}
}

func (h *UserHandler) HandleSignUpFailed() {
	_, err := h.NATSClient.Subscribe("SignUpFailed", func(m *nats.Msg) {
		fmt.Println("Okinuo se sign up failed u user management ms")
		var event SignUpFailedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Printf("Error unmarshalling SignUpFailed event: %v", err)
			h.TransactionManager.Rollback(event.UserID, context.Background())
			return
		}

		log.Printf("Received SignUpFailed for user ID %d, Error: %s", event.UserID, event.Error)
		// Here, you would handle the rollback of the user creation
		// Example: rollbackUserCreation(event.UserID)

		h.TransactionManager.Rollback(event.UserID, context.Background())
	})

	if err != nil {
		// TODO: What about transactions handling?
		log.Fatalf("Failed to subscribe to SignUpFailed events: %v", err)
	}

}

func (h *UserHandler) UserSignUp(ctx context.Context, req *user.UserSignUpRequest) (*user.UserSignUpResponse, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email, username, and password cannot be empty")
	}

	if req.Role == 0 || req.Role > 3 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value for role")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email format")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}

	tx := initializers.DB.Begin()

	userr := models.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
		Role:     models.UserRole(req.Role),
	}

	if err := tx.Create(&userr).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	graphUser := models.GraphDBUser{
		ID:       int64(userr.ID),
		Username: userr.Username,
	}

	// neo4j_tx, err := session.BeginTransaction(ctx)
	neo4j_tx, err := h.Neo4jSession.BeginTransaction(ctx)
	if err != nil {
		panic(err)
	}

	// err = graphdb.WriteUser(&graphUser, initializers.Neo4JDriver)
	// if err != nil {
	// 	log.Printf("Error creating user node for %s in Neo4j database: %v", userr.Username, err)
	// 	tx.Rollback()
	// 	neo4j_tx.Rollback(ctx)
	// 	return nil, status.Errorf(codes.Internal, "Failed to create user in the neo4j: %v", err)
	// }
	// err = graphdb.WriteUser(&graphUser, session, ctx)
	err = graphdb.WriteUser(&graphUser, ctx, neo4j_tx)
	if err != nil {
		tx.Rollback()

		if err.Error() == "user already exists" {
			return nil, status.Errorf(codes.Internal, "Failed to create user in Neo4j: ", err)
		}

		// Desio se error nakon pokusaja cuvanja korisnika u grafskoj bazi.
		return nil, status.Errorf(codes.Internal, "Failed to create user in Neo4j: ", err)
	}
	// neo4j_tx.Commit(ctx)
	// neo4j_tx.Close(ctx)

	// if err != nil {
	// 	log.Printf("Username already exists: %v", err)
	// 	// log.Printf("Error in graph database: %v", err)
	// 	tx.Rollback()
	// 	neo4j_tx.Rollback(ctx)
	// 	return nil, status.Errorf(codes.Internal, "Failed to create user in Neo4j: ", err)
	// }

	// signUpResp, err := s.AuthServiceClient.SignUp(ctx, &user.SignUpRequest{
	// 	Id:       uint32(userr.ID),
	// 	Username: req.Username,
	// 	Password: string(hash),
	// 	Role:     req.Role,
	// })

	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "Failed to register user in auth service: %v", err)
	// }

	// TODO: Should I delete the transaction afterwards?

	// Publishing event through NATS
	userCreatedEvent := map[string]interface{}{
		"Id":       uint32(userr.ID),
		"Username": req.Username,
		"Password": string(hash),
		"Role":     req.Role,
	}

	eventData, err := json.Marshal(userCreatedEvent)
	if err != nil {
		tx.Rollback()
		neo4j_tx.Rollback(ctx)
		return nil, status.Errorf(codes.Internal, "Failed to marshal user created event: %v", err)
	}

	h.TransactionManager.SavePendingTransaction(uint32(userr.ID), tx, neo4j_tx)

	// log.Fatalf("Transakcija save pending, sacuvana u recniku: duzine je %v", len(s.TransactionManager.Transactions))

	fmt.Println("Prosao prethodni ispis u vei transakcije save pending")

	err = h.NATSClient.Publish("UserCreated", eventData)
	fmt.Println("Prosao publish dogadjaja usercreated")
	if err != nil {
		// tx.Rollback()
		h.TransactionManager.Rollback(uint32(userr.ID), ctx)
		fmt.Println("Usao u roll back zbog loseg publisha")
		return nil, status.Errorf(codes.Internal, "Failed to publish user created event: %v", err)
	}

	// if err := tx.Commit().Error; err != nil {
	// 	return nil, status.Errorf(codes.Internal, "Failed to commit transaction %v", err)
	// }

	return &user.UserSignUpResponse{
		Status: "STATUS: Pending. \nAccount needs to be confirmed by the administrator. Try to log in in a few minutes.", // I cant forward the answer later, since I wouldnt receive a response, I can make another consumer and another event to listen for actions result
	}, nil

	// return &user.UserSignUpResponse{
	// 	Status: signUpResp.Status,
	// }, nil
}

// func SignUp(c *gin.Context) {
// 	var body struct {
// 		Username string
// 		Password string
// 		Email    string
// 		Role     models.UserRole
// 	}

// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to read body",
// 		})
// 		return
// 	}

// 	if body.Email == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Email cannot be empty",
// 		})
// 		return
// 	}

// 	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
// 	if !emailRegex.MatchString(body.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid email format",
// 		})
// 		return
// 	}

// 	if body.Role.String() == "Unknown" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Role is not correct or can not be empty",
// 		})
// 		return
// 	}

// 	if body.Password == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Password cannot be empty",
// 		})
// 		return
// 	}

// 	if body.Username == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Username cannot be empty",
// 		})
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to hash password",
// 		})
// 		return
// 	}

// 	tx := initializers.DB.Begin()

// 	user := models.User{
// 		Username: body.Username,
// 		Password: string(hash),
// 		Email:    body.Email,
// 		Role:     body.Role,
// 	}

// 	if err := tx.Create(&user).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create user. The log is: " + err.Error(),
// 		})
// 		return
// 	}

// 	graphUser := models.GraphDBUser{
// 		ID:       int64(user.ID),
// 		Username: user.Username,
// 	}

// 	if err := graphdb.WriteUser(&graphUser, initializers.Neo4JDriver); err != nil {
// 		log.Printf("Error creating user node for %s in Neo4j database: %v", user.Username, err)
// 		tx.Rollback()
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create user. The log is: " + err.Error(),
// 		})
// 		return
// 	}

// 	requestBody, err := json.Marshal(map[string]interface{}{
// 		"ID":       user.ID,
// 		"Username": body.Username,
// 		"Password": string(hash),
// 		"Role":     body.Role,
// 	})
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to marshal request body",
// 		})
// 		return
// 	}

// 	resp, err := http.Post("http://auth-service:3001/signup", "application/json", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to send request",
// 		})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		tx.Rollback()
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": fmt.Sprintf("Failed to create user. Status code: %d", resp.StatusCode),
// 		})
// 		fmt.Println(resp)
// 		return
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to commit transaction",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{})
// }

func (h *UserHandler) Validate(ctx context.Context, req *emptypb.Empty) (*user.ValidateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve metadata")
	}

	roles, ok := md["role"]
	if !ok || len(roles) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "User role not found")
	}
	role := roles[0]

	var responseContent string
	switch role {
	case "Administrator":
		responseContent = "Administrator Content."
	case "Guide":
		responseContent = "Guide Content."
	default:
		responseContent = "Tourist Content."
	}

	return &user.ValidateResponse{
		ResponseContent: responseContent,
	}, nil
}

// func Validate(c *gin.Context) {

// 	// ako hoces da dobavis neko polje usera
// 	// onda zapocinjes komandu na sledeci nacin
// 	// user.(models.User).

// 	userInterface, _ := c.Get("user")

// 	// NE SME *models.User!
// 	user, _ := userInterface.(models.User)

// 	if user.Role.String() == "Administrator" {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Administrator Content.",
// 		})
// 	} else if user.Role.String() == "Guide" {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Guide Content.",
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Tourist Content.",
// 		})
// 	}
// }

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
