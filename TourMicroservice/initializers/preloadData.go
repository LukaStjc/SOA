package initializers

import (
	"fmt"
	"go-tourm/models"
	"time"
)

func PreloadTours() {
	var tours = []models.Tour{
		{
			ID:          1,
			Name:        "Poseta Baru i Susnju",
			Description: "prva tura, veoma zanimljiva",
			Type:        models.TourType(1),
			Tags:        "susanj;bar",
			Price:       2500.55,
			AvgRate:     0,
			UserID:      3,
		},
		{
			ID:          2,
			Name:        "Poseta Zlatiboru i Cajetini",
			Description: "druga tura, vrlo interesantna",
			Type:        models.TourType(2),
			Tags:        "zlatibor;cajetina",
			Price:       3500.35,
			AvgRate:     0,
			UserID:      3,
		},
	}

	// Create tours
	for _, t := range tours {
		DB.Create(&t)
	}

	// Preload key points for each tour
	PreloadKeyPoints()

	PreloadRest()

}

func PreloadKeyPoints() {
	var keyPoints = []models.KeyPoint{
		{
			ID:        1,
			Longitude: 42.1,
			Latitude:  19.1,
			TourID:    1,
		},
		{
			ID:        2,
			Longitude: 42.3,
			Latitude:  19.8,
			TourID:    1,
		},
		{
			ID:        3,
			Longitude: 42.2,
			Latitude:  19.9,
			TourID:    1,
		},
		{
			ID:        4,
			Longitude: 43.734520,
			Latitude:  19.64,
			TourID:    2,
		},
		{
			ID:        5,
			Longitude: 43.74520,
			Latitude:  19.73140,
			TourID:    2,
		},
		{
			ID:        6,
			Longitude: 43.714520,
			Latitude:  19.623140,
			TourID:    2,
		},
	}

	// Create key points
	for _, kp := range keyPoints {
		DB.Create(&kp)
	}
}

func PreloadRest() {
	// Seed ShoppingCart and OrderItems
	shoppingCart := []models.ShoppingCart{
		{ID: 1, UserID: 1, Price: 10000.90},
		{ID: 2, UserID: 2, Price: 7000.45},
	}
	DB.Create(&shoppingCart)

	orderItems := []models.OrderItem{
		{ID: 1, TourID: 1, TourRate: 0, CartID: shoppingCart[0].ID, TourName: "Poseta Baru i Susnju", TourPrice: 2500.55, NumberOfPeople: 2},
		{ID: 2, TourID: 2, TourRate: 0, CartID: shoppingCart[1].ID, TourName: "Poseta Zlatiboru i Cajetini", TourPrice: 3500.35, NumberOfPeople: 3},
		{ID: 3, TourID: 1, TourRate: 0, CartID: shoppingCart[1].ID, TourName: "Poseta Baru i Susnju", TourPrice: 6000.35, NumberOfPeople: 4},
	}
	DB.Create(&orderItems)

	// Seed TourCheckIns
	tourCheckIns := []models.TourCheckIn{
		// UserID 1
		{ID: 1, KeyPointID: 1, UserID: 1, VisitingTime: time.Now()},
		{ID: 2, KeyPointID: 2, UserID: 1, VisitingTime: time.Now()},
		{ID: 3, KeyPointID: 3, UserID: 1, VisitingTime: time.Now()},

		// UserID 2, druga tura
		{ID: 4, KeyPointID: 4, UserID: 2, VisitingTime: time.Now()},
		{ID: 5, KeyPointID: 5, UserID: 2, VisitingTime: time.Now()},
		{ID: 6, KeyPointID: 6, UserID: 2, VisitingTime: time.Now()},
		// UserID 2, prva tura
		{ID: 7, KeyPointID: 1, UserID: 2, VisitingTime: time.Now()},
		// visiting time kasnije dodajemo kao primer
		{ID: 8, KeyPointID: 2, UserID: 2, VisitingTime: time.Time{}}, // garbage time value, nesto kao nil, proverava se sa isZero
		{ID: 9, KeyPointID: 3, UserID: 2, VisitingTime: time.Time{}}, // garbage time value, nesto kao nil, proverava se sa isZero
	}

	if err := DB.Create(&tourCheckIns).Error; err != nil {
		fmt.Printf("Error creating tour check-ins: %v\n", err)
	}
}
