package main

import (
	"context"
	"go-userm/controllers"
	"go-userm/initializers"
	"go-userm/middleware"
	configurations "go-userm/startup"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func init() {
	configuration := configurations.NewConfigurations()
	// initializers.LoadEnvVariables()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	// initializers.PreloadUsers()
}

func main() {

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
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	initializers.Neo4JDriver = driver
	// initializers.Ctx = ctx

	initializers.PreloadUsers()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.POST("/signup", controllers.SignUp)

	r.GET("/:id", controllers.GetById)

	// RequireAuth
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/follow/:username", middleware.RequireAuth, controllers.Follow)
	r.GET("/is-blocked/:id", middleware.RequireAuth, controllers.IsBlocked)
	r.GET("/does-follow/:followerId/:creatorId", middleware.RequireAuth, controllers.DoesFollow)
	//r.GET("/does-follow/:creatorId", middleware.RequireAuth, controllers.DoesFollow)

	// RequireAuth + CheckIfAdmin
	r.PUT("/ban/:username", middleware.RequireAuth, middleware.CheckIfAdmin, controllers.BlockUser)

	r.GET("/get-followed", middleware.RequireAuth, controllers.GetFollowed)

	// RequireAuth
	r.GET("/get-friends-recommendation", middleware.RequireAuth, controllers.GetFriendsRecommendation)

	r.Run()

}
