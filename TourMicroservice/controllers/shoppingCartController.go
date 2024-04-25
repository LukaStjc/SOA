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

func RemoveFromShoppingCart(c *gin.Context) {
	// Extracting the orderItem id from the path
	orderItemId := c.Param("orderItemId")

	// Extracting the ShoppingCart id from the path
	shoppingCartId := c.Param("shoppingCartId")

	// Find the tour by id
	var orderItem models.OrderItem
	resultOrderItem := initializers.DB.Where("id = ?", orderItem).First(&orderItemId)

	// Find the shoppingCart by id
	var shoppingCart models.ShoppingCart
	resultShoppingCart := initializers.DB.Where("id = ?", shoppingCartId).First(&shoppingCart)

	if resultOrderItem.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OrderItem not found"})
		return
	}

	if resultShoppingCart.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShoppingCart not found"})
		return
	}

	//brisemo orderItem iz shoppingCart-a
	indexToDelete := -1
	for i, item := range shoppingCart.OrderItems {
		if item == &orderItem {
			indexToDelete = i
			break
		}
	}

	if indexToDelete != -1 {
		shoppingCart.OrderItems = append(shoppingCart.OrderItems[:indexToDelete], shoppingCart.OrderItems[indexToDelete+1:]...)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "OrderItem not found in the ShoppingCart"})
	}

	//smanjujemo cenu u shoppingCart-u
	shoppingCart.Price -= orderItem.TourPrice

	c.JSON(http.StatusOK, gin.H{"shoppingCart": shoppingCart})
}
