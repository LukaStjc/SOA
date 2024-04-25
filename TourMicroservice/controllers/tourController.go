package controllers

import (
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTour(c *gin.Context) {
	var body struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Type        uint              `json:"type"`
		Tags        string            `json:"tags"`
		Price       float64           `json:"price"`
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
		kp.TourID = int(tour.ID)
		if kpResult := initializers.DB.Create(&kp); kpResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create keypoints",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": tour})
}

func GetToursByUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var tours []models.Tour
	result := initializers.DB.Preload("KeyPoints").Where("user_id = ?", userID).Find(&tours)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tours": tours})
}
