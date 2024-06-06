package main

import (
	"context"
	"go-userm/controllers"
	"go-userm/initializers"
	"go-userm/interceptors"
	user "go-userm/proto/user/generatedFiles"
	configurations "go-userm/startup"
	"go-userm/transactionmanager"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var authServiceClient user.AuthServiceClient

var natsClient *nats.Conn

func initAuthServiceClient() {
	conn, err := grpc.Dial("auth-service:3001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to auth service: %v", err)
	}
	authServiceClient = user.NewAuthServiceClient(conn)
}

func initNATSClient() *nats.Conn {
	conn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	// controllers.HandleSignUpFailed(conn)
	// controllers.HandleSignUpSuccess(conn)

	return conn
}

func init() {
	configuration := configurations.NewConfigurations()
	// initializers.LoadEnvVariables()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	// initializers.PreloadUsers()

}

func main() {
	// NATS

	natsClient := initNATSClient()
	defer natsClient.Close()
	txManager := transactionmanager.NewManager()

	ctx := context.Background()
	// dbUri := fmt.Sprintf("bolt://%s:%s", os.Getenv("USER_GRAPH_DB_HOST"), os.Getenv("USER_GRAPH_DB_PORT"))
	// dbUser := os.Getenv("USER_GRAPH_DB_USERNAME")
	// dbPassword := os.Getenv("USER_GRAPH_DB_PASS")
	// driver, err := neo4j.NewDriverWithContext(
	// 	dbUri,
	// 	neo4j.BasicAuth(dbUser, dbPassword, ""))
	// if err != nil {
	// 	panic(err)
	// }
	// defer driver.Close(ctx)

	// err = driver.VerifyConnectivity(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	dbUri := "bolt://neo4j:7687"
	dbUser := "neo4j"
	dbPassword := "nekaSifra"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	initializers.Neo4JDriver = driver
	// initializers.Ctx = ctx

	initializers.PreloadUsers()

	// r := gin.Default()

	// r.Use(middleware.CORSMiddleware())

	// r.POST("/signup", controllers.SignUp)

	// r.GET("/:id", controllers.GetById)

	// // RequireAuth
	// r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	// r.POST("/follow/:username", middleware.RequireAuth, controllers.Follow)
	// r.GET("/is-blocked/:id", middleware.RequireAuth, controllers.IsBlocked)
	// r.GET("/does-follow/:followerId/:creatorId", middleware.RequireAuth, controllers.DoesFollow)
	// //r.GET("/does-follow/:creatorId", middleware.RequireAuth, controllers.DoesFollow)

	// // RequireAuth + CheckIfAdmin
	// r.PUT("/ban/:username", middleware.RequireAuth, middleware.CheckIfAdmin, controllers.BlockUser)

	// r.GET("/get-followed", middleware.RequireAuth, controllers.GetFollowed)

	// // RequireAuth
	// r.GET("/get-friends-recommendation", middleware.RequireAuth, controllers.GetFriendsRecommendation)

	// r.Run()

	// Konekcija ka auth ms
	initAuthServiceClient()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor(authServiceClient)))
	userHandler := &controllers.UserHandler{
		AuthServiceClient:  authServiceClient,
		NATSClient:         natsClient,
		TransactionManager: txManager,
		Neo4jDriver:        initializers.Neo4JDriver,
		Neo4jSession:       session,
	}
	user.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	// SAMO NATS U NAREDNE DVE LINIJE
	userHandler.HandleSignUpFailed()
	userHandler.HandleSignUpSuccess()
	// DOVDE

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

	// // N

}
