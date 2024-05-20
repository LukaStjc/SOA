package main

import (
	"context"
	"example/gateway/config"
	auth "example/gateway/proto/auth/generatedFiles"
	user "example/gateway/proto/user/generatedFiles"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.GetConfig()

	authConn, err := grpc.DialContext(
		context.Background(),
		os.Getenv("AUTH_SERVICE_ADDRESS"),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial Auth Service:", err)
	}
	defer authConn.Close()

	userConn, err := grpc.DialContext(
		context.Background(),
		os.Getenv("USER_SERVICE_ADDRESS"),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial User Management Service:", err)
	}
	defer userConn.Close()

	gwmux := runtime.NewServeMux()
	authServiceClient := auth.NewAuthServiceClient(authConn)
	if err := auth.RegisterAuthServiceHandlerClient(context.Background(), gwmux, authServiceClient); err != nil {
		log.Fatalln("Failed to register Auth Service gateway:", err)
	}

	userServiceClient := user.NewUserServiceClient(userConn)
	if err := user.RegisterUserServiceHandlerClient(context.Background(), gwmux, userServiceClient); err != nil {
		log.Fatalln("Failed to register User Service gateway:", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3005"}, // or use * to allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	gwServer := &http.Server{
		Addr:    cfg.Address,
		Handler: c.Handler(gwmux),
	}

	// gwServer := &http.Server{
	// 	Addr:    cfg.Address,
	// 	Handler: gwmux,
	// }

	go func() {
		if err := gwServer.ListenAndServe(); err != nil {
			log.Fatal("Server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh

	if err = gwServer.Close(); err != nil {
		log.Fatalln("Error while stopping server:", err)
	}

}
