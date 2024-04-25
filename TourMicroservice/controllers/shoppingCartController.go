package controllers

import (
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"

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
