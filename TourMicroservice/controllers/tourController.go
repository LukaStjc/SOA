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
		Type        models.TourType   `json:"type"`
		Tags        string            `json:"tags"`
		Price       float64           `json:"price"`
		UserID      uint              `json:"userId"`
		KeyPoints   []models.KeyPoint `json:"keyPoints"` // Add this line
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
		Type:        body.Type,
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

	// Create keypoints if any are provided
	for _, kp := range body.KeyPoints {
		kp.TourID = int(tour.ID) // Ensure the foreign key is set correctly
		if kpResult := initializers.DB.Create(&kp); kpResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create keypoints",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": tour})
}
