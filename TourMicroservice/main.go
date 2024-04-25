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

	r.POST("/create-shoppingCart" /*middleware.RequireAuth, middleware.CheckIfTourist,*/, controllers.CreateShoppingCart)
	r.PUT("/clear-shoppingCart/:id" /* middleware.RequireAuth, middleware.CheckIfTourist,*/, controllers.ClearShoppingCart)
	r.PUT("/addToShoppingCart/:tourId/:shoppingCartId/:numOfPeople" /*middleware.RequireAuth, middleware.CheckIfTourist,*/, controllers.AddToShoppingCart)
	r.PUT("/removeFromShoppingCart/:orderItemId/:shoppingCartId" /*middleware.RequireAuth, middleware.CheckIfTourist,*/, controllers.RemoveFromShoppingCart)
	r.Run()

}
