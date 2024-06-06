package controllers

import (
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTour(c *gin.Context) {
	var body struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Type        uint              `json:"type"`
		Tags        string            `json:"tags"`
		Price       float64           `json:"price"`
		AvgRate     float64           `json:"avgRate"`
		UserID      uint              `json:"userId"`
		KeyPoints   []models.KeyPoint `json:"keyPoints"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	tour := models.Tour{
		Name:        body.Name,
		Description: body.Description,
		Type:        models.TourType(body.Type),
		Tags:        body.Tags,
		Price:       body.Price,
		AvgRate:     body.AvgRate,
		UserID:      body.UserID,
	}

	result := initializers.DB.Create(&tour)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create tour",
		})
		return
	}

	for _, kp := range body.KeyPoints {
		kp.TourID = tour.ID
		if kpResult := initializers.DB.Create(&kp); kpResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create keypoints",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": tour})
}

func CheckCanUserLeaveRate(c *gin.Context, tour models.Tour, user models.User) bool {
	var shoppingCarts []models.ShoppingCart
	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&shoppingCarts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping carts"})
		return false
	}

	var isABuyer bool
	for _, cart := range shoppingCarts {
		var orderItems []models.OrderItem
		if err := initializers.DB.Where("cart_id = ?", cart.ID).Find(&orderItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return false
		}

		for _, item := range orderItems {
			// takodje, ne moze ostaviti ocenu ako je vec ostavio
			if item.TourID == tour.ID && item.TourRate == 0 {
				isABuyer = true
				break
			}
		}
		if isABuyer {
			break
		}
	}

	if !isABuyer {
		return false
	}

	var numberVisitedKeyPoints int
	var totalKeyPoints = len(tour.KeyPoints)

	for _, keyPoint := range tour.KeyPoints {
		var checkIns []models.TourCheckIn
		if err := initializers.DB.Where("key_point_id = ? AND user_id = ? AND NOT visiting_time = '0001-01-01 00:00:00'", keyPoint.ID, user.ID).Find(&checkIns).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch check-ins"})
			return false
		}
		if len(checkIns) > 0 {
			numberVisitedKeyPoints++
		}
	}

	return numberVisitedKeyPoints > (totalKeyPoints / 2)
}

func GetToursByUser(c *gin.Context) {
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User data is not valid"})
		return
	}

	var tours []models.Tour
	if err := initializers.DB.Preload("KeyPoints").Find(&tours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours"})
		return
	}

	// fmt.Println("")

	var listTourRateView []map[string]interface{}
	for _, tour := range tours {
		canRate := CheckCanUserLeaveRate(c, tour, user)
		listTourRateView = append(listTourRateView, map[string]interface{}{
			"tour": tour,
			"rate": canRate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"listTourRateView": listTourRateView})
}

func PostReview(c *gin.Context) {
	var body struct {
		TourID uint `json:"tourID"`
		Rate   uint `json:"rate"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User data is not valid"})
		return
	}

	// Find all shopping carts for the user
	var shoppingCarts []models.ShoppingCart
	if err := initializers.DB.Where("user_id = ?", user.ID).Find(&shoppingCarts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping carts"})
		return
	}

	var updated bool
	// Trebalo bi samo za JEDNU korpu! Odnosno samo da u jednoj korpi se nalazi ta tura
	// Iterate over each shopping cart
	for _, cart := range shoppingCarts {
		var orderItems []models.OrderItem
		// Fetch all order items within this cart
		if err := initializers.DB.Where("cart_id = ?", cart.ID).Find(&orderItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
			return
		}

		// Update the tour rate for relevant order items
		for _, item := range orderItems {
			if item.TourID == body.TourID {
				item.TourRate = body.Rate
				if err := initializers.DB.Save(&item).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tour rate"})
					return
				}
				updated = true
				break // Break if we update at least one item
			}
		}
		if updated {
			break // Break if any item was updated
		}
	}

	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order items found for the given tour"})
		return
	}

	var result struct {
		TotalRate     int
		NumberOfRates int
	}

	if err := initializers.DB.Model(&models.OrderItem{}).
		Where("tour_id = ? AND tour_rate > 0", body.TourID).
		Select("sum(tour_rate) as total_rate, count(*) as number_of_rates").
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average rate"})
		return
	}

	if result.NumberOfRates > 0 {
		var tour models.Tour
		if err := initializers.DB.First(&tour, body.TourID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tour not found"})
			return
		}
		tour.AvgRate = float64(result.TotalRate) / float64(result.NumberOfRates)
		if err := initializers.DB.Save(&tour).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update average rate"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully", "avgRate": tour.AvgRate})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Review updated but no average rate calculated due to lack of rates"})
		return
	}

}

// Nebitna
func Validate(c *gin.Context) {

	// ako hoces da dobavis neko polje usera
	// onda zapocinjes komandu na sledeci nacin
	// user.(models.User).

	userInterface, _ := c.Get("user")

	// NE SME *models.User!
	user, _ := userInterface.(models.User)

	if user.Role.String() == "Administrator" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Administrator Content.",
		})
	} else if user.Role.String() == "Guide" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Guide Content.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tourist Content.",
		})
	}
}
