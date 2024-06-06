package main

import (
	"go-tourm/controllers"
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

	r.POST("/create-tour", middleware.RequireAuth, controllers.CreateTour)
	r.GET("/guide/:id/tours", middleware.RequireAuth, controllers.GetToursByUser)

	r.POST("/postReview", middleware.RequireAuth, controllers.PostReview)

	// nebitna
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()

}
