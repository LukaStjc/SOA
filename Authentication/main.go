package main

import (
	"go-jwt/controllers"
	"go-jwt/initializers"
	auth "go-jwt/proto/generatedFiles"
	configurations "go-jwt/startup"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	//initializers.LoadEnvVariables()
	configuration := configurations.NewConfigurations()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadUsers()

}

func main() {

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
