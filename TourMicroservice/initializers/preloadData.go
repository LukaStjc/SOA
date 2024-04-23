package initializers

import (
	"go-tourm/models"
	"time"

	"gorm.io/gorm"
)

func PreloadTours() {
	var tours = []models.Tour{
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 57, 29, 351794000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 57, 29, 351794000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Name:        "tura1",
			Description: "prva tura, veoma zanimljiva oca mi",
			Type:        models.TourType(1),
			Tags:        "susanj;bar",
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Name:        "tura2",
			Description: "druga tura, veoma zanimljiva matere mi, ali malo teza",
			Type:        models.TourType(2),
			Tags:        "zlatibor;cajetina",
		},
	}

	for _, t := range tours {
		DB.Create(&t)
	}
}
