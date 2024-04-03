package main

import (
	"go-userm/controllers"
	"go-userm/initializers"
	"go-userm/middleware"

	//configurations "go-userm/startup"

	"github.com/gin-gonic/gin"
)

func init() {
	//initializers.LoadEnvVariables()
	configuration := configurations.NewConfigurations()
	//initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadUsers()
}

func main() {

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

	r.Run()

}
