package initializers

import (
	"go-tourm/models"
)

func PreloadTours() {
	var tours = []models.Tour{
		{
			Name:        "tura1",
			Description: "prva tura, veoma zanimljiva oca mi",
			Type:        models.TourType(1),
			Tags:        "susanj;bar",
			Price:       2500.55,
			UserID:      1,
		},
		{
			Name:        "tura2",
			Description: "druga tura, veoma zanimljiva matere mi, ali malo teza",
			Type:        models.TourType(2),
			Tags:        "zlatibor;cajetina",
			Price:       3500.35,
			UserID:      1,
		},
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Tour{}, &models.KeyPoint{})

	// Create tours
	for _, t := range tours {
		DB.Create(&t)
	}

	// Preload key points for each tour
	PreloadKeyPoints()

}

func PreloadKeyPoints() {
	var keyPoints = []models.KeyPoint{
		{
			Name:      "Pocetna",
			Longitude: 42.1,
			Latitude:  19.1,
			TourID:    1,
		},
		{
			Name:      "Poslednja",
			Longitude: 42.1,
			Latitude:  19.9,
			TourID:    1,
		},
		{
			Name:      "Pocetna",
			Longitude: 43.734520,
			Latitude:  19.64,
			TourID:    2,
		},
		{
			Name:      "Poslednja",
			Longitude: 43.734520,
			Latitude:  19.693140,
			TourID:    2,
		},
	}

	// Create key points
	for _, kp := range keyPoints {
		DB.Create(&kp)
	}
}
