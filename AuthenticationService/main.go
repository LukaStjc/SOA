package main

import (
	"go-jwt/controllers"
	"go-jwt/initializers"

	configurations "go-jwt/startup"

	"github.com/gin-gonic/gin"
)

func init() {
	//initializers.LoadEnvVariables()
	configuration := configurations.NewConfigurations()
	initializers.ConnectToDb(configuration)
	initializers.SyncDatabase()
	initializers.PreloadUsers()

}

func main() {

	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Run()
}
