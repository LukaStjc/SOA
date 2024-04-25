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

	r.POST("/create-tour", controllers.CreateTour)
	r.GET("/guide/:id/tours", controllers.GetToursByUser)

	r.POST("/create-shoppingCart", controllers.CreateShoppingCart)
	r.PUT("/clear-shoppingCart/:id", controllers.ClearShoppingCart)
	r.PUT("/addToShoppingCart/:tourId/:shoppingCartId/:numOfPeople", controllers.AddToShoppingCart)
	r.PUT("/removeFromShoppingCart/:orderItemId/:shoppingCartId", controllers.RemoveFromShoppingCart)
	r.Run()

}
