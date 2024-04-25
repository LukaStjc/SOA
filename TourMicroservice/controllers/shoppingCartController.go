package controllers

import (
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateShoppingCart(c *gin.Context) {

	var body struct {
		UserID uint `json:"userId"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//na pocetku nema orderItems
	shoppingCart := models.ShoppingCart{
		UserID: body.UserID,
		Price:  0, //cena je na pocetku 0
	}

	result := initializers.DB.Create(&shoppingCart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create shoppingCart",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": shoppingCart})
}

func ClearShoppingCart(c *gin.Context) {
	// Immediately return if the previous middleware aborted the request
	if c.IsAborted() {
		return
	}
	// Extracting the shoppingCart id from the path

	shoppingCartId := c.Param("id")

	// Find the shoppingCart by id
	var shoppingCart models.ShoppingCart
	result := initializers.DB.Where("id = ?", shoppingCartId).First(&shoppingCart)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShoppingCart not found"})
		return
	}

	//isprazniti sve orderItems iz shoppingCart-a
	shoppingCart.OrderItems = nil

	// When orderItems are cleared, send OK status
	c.JSON(http.StatusOK, gin.H{"message": "OrderItems are cleared."})
}

func AddToShoppingCart(c *gin.Context) {
	//bice prosledjena tourId i shoppingCartId

	// Extracting the Tour id from the path
	tourId := c.Param("tourId")

	// Extracting the ShoppingCart id from the path
	shoppingCartId := c.Param("shoppingCartId")

	// Find the tour by id
	var tour models.Tour
	resultTour := initializers.DB.Where("id = ?", tourId).First(&tour)

	// Find the shoppingCart by id
	var shoppingCart models.ShoppingCart
	resultShoppingCart := initializers.DB.Where("id = ?", shoppingCartId).First(&shoppingCart)

	if resultTour.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
		return
	}

	if resultShoppingCart.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShoppingCart not found"})
		return
	}

	numOfPeople, _ := strconv.Atoi(c.Param("numOfPeople"))

	//pravimo orderItem od te ture
	orderItem := models.OrderItem{
		TourId:         tour.ID,
		TourName:       tour.Name,
		TourPrice:      tour.Price,
		NumberOfPeople: numOfPeople,
	}

	//dodajemo orderItem u shoppingCart
	shoppingCart.OrderItems = append(shoppingCart.OrderItems, &orderItem)

	//povecavamo cenu u shoppingCart-u
	shoppingCart.Price += orderItem.TourPrice

	c.JSON(http.StatusOK, gin.H{"shoppingCart": shoppingCart})
}
