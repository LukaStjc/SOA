package main

import (
	"context"
	"encoding/json"
	"go-jwt/controllers"
	"go-jwt/initializers"
	auth "go-jwt/proto/generatedFiles"
	configurations "go-jwt/startup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var natsClient *nats.Conn

func initNATSClient() *nats.Conn {
	conn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	return conn
}

func init() {
	//initializers.LoadEnvVariables()
	configuration := configurations.NewConfigurations()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadUsers()
	// initNATSClient()
}

func main() {
	// NATS
	natsClient := initNATSClient()
	defer natsClient.Close()

	_, err := natsClient.Subscribe("UserCreated", func(m *nats.Msg) {
		var userCreatedEvent map[string]interface{}
		err := json.Unmarshal(m.Data, &userCreatedEvent)
		if err != nil {
			log.Printf("Failed to unmarshal user created event: %v", err)
			return
		}

		// Call SignUp in AuthService
		ctx := context.Background()
		signUpRequest := &auth.SignUpRequest{
			Id:       uint32(userCreatedEvent["Id"].(float64)),
			Username: userCreatedEvent["Username"].(string),
			Password: userCreatedEvent["Password"].(string),
			Role:     uint32(userCreatedEvent["Role"].(float64)),
		}
		authHandler := &controllers.AuthHandler{
			NATSClient: natsClient,
		}

		_, err = authHandler.SignUp(ctx, signUpRequest)
		if err != nil {
			log.Printf("Failed to sign up user in auth service: %v", err)
			return
		}

		log.Printf("Successfully signed up user %s in auth service", signUpRequest.Username)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to UserCreated: %v", err)
	}

	// r := gin.Default()

	// r.Use(middleware.CORSMiddleware())

	// // r.POST("/signup", controllers.SignUp)
	// // r.POST("/login", controllers.Login)
	// r.POST("/authenticate", controllers.Authenticate)

	// r.Run()

	// // Load environment variables and connect to the database
	// configuration := configurations.NewConfigurations()
	// // initializers.LoadEnvVariables() Nakon ove linije proradio, nije mogao ucitati varijable i to se bunio
	// initializers.ConnectToDb(configuration)
	// initializers.SyncDatabase()
	// initializers.PreloadUsers()

	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, &controllers.AuthHandler{})
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	grpcServer.GracefulStop()
	lis.Close()
	log.Println("Shutting down gRPC server...")

}
