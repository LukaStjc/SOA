package main

import (
	"go-userm/controllers"
	"go-userm/initializers"
	"go-userm/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.PreloadUsers()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()

}
