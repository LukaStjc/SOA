package main

import (
	"go-tourm/initializers"
	"go-tourm/middleware"
	configurations "go-tourm/startup"

	"github.com/gin-gonic/gin"
)

func init() {
	configuration := configurations.NewConfigurations()
	// initializers.LoadEnvVariables()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadTours()
}

func main() {

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	//r.POST("/signup", controllers.SignUp)

	//r.GET("/:id", controllers.GetById)

	// RequireAuth
	// r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	// r.POST("/follow/:username", middleware.RequireAuth, controllers.Follow)
	// r.GET("/is-blocked/:id", middleware.RequireAuth, controllers.IsBlocked)
	// r.GET("/does-follow/:followerId/:creatorId", middleware.RequireAuth, controllers.DoesFollow)
	// //r.GET("/does-follow/:creatorId", middleware.RequireAuth, controllers.DoesFollow)

	// RequireAuth + CheckIfAdmin
	// r.PUT("/ban/:username", middleware.RequireAuth, middleware.CheckIfAdmin, controllers.BlockUser)

	// r.GET("/get-followed", middleware.RequireAuth, controllers.GetFollowed)

	r.Run()

}
