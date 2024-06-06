package controllers

import (
	"errors"
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

/*
func CreateTour(c *gin.Context) {
	// Start a new trace
	_, span := otel.Tracer(serviceName).Start(c.Request.Context(), "weapon-get")
	defer span.End()

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

	// Create a new subtrace for communicating with the database
	span.AddEvent("Establishing connection to the database")
	db := initializers.DB
	if db == nil {
		return
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
}*/

func CreateTour(c *gin.Context) {
	// Start a new trace
	traceContext, span := otel.Tracer(serviceName).Start(c.Request.Context(), "CreateTour")
	defer func() { span.End() }()

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

	// Pass the traceContext to the database operation
	tour, err := createTour(traceContext, body, span)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create tour",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tour})
}

func createTour(ctx context.Context, body struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        uint              `json:"type"`
	Tags        string            `json:"tags"`
	Price       float64           `json:"price"`
	UserID      uint              `json:"userId"`
	KeyPoints   []models.KeyPoint `json:"keyPoints"`
}, trace.Span span) (*models.Tour, error) {

	span.AddEvent("Establishing connection to the database...")
	// Retrieve database instance
	db := initializers.DB
	if db == nil {
		return nil, errors.New("database instance is nil")
	}

	// Create a new tour instance
	tour := models.Tour{
		Name:        body.Name,
		Description: body.Description,
		Type:        models.TourType(body.Type),
		Tags:        body.Tags,
		Price:       body.Price,
		UserID:      body.UserID,
	}

	// Pass the context to the Create method
	result := db.WithContext(ctx).Create(&tour)
	if result.Error != nil {
		return nil, result.Error
	}

	// Create key points for the tour
	for _, kp := range body.KeyPoints {
		kp.TourID = int(tour.ID)
		if kpResult := db.Create(&kp); kpResult.Error != nil {
			return nil, kpResult.Error
		}
	}

	return &tour, nil
}

func GetToursByUser(c *gin.Context) {
	traceContext, span := otel.Tracer(serviceName).Start(c.Request.Context(), "GetToursByUser")
	defer func() { span.End() }()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create a new subtrace for communicating with the database
	span.AddEvent("Establishing connection to the database")
	db := initializers.DB
	if db == nil {
		return
	}

	/*var tours []models.Tour
	result := initializers.DB.Preload("KeyPoints").Where("user_id = ?", userID).Find(&tours)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours"})
		return
	}*/
	tours, err := getToursByUserID(traceContext, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tours": tours})
}

func getToursByUserID(ctx context.Context, userID uint64) ([]models.Tour, error) {
	// Retrieve database instance
	db := initializers.DB
	if db == nil {
		return nil, errors.New("database instance is nil")
	}

	// Fetch tours using the provided userID and the provided context
	var tours []models.Tour
	result := db.WithContext(ctx).Preload("KeyPoints").Where("user_id = ?", userID).Find(&tours)
	if result.Error != nil {
		return nil, result.Error
	}

	return tours, nil
}
