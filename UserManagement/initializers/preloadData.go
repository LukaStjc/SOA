package initializers

import (
	"go-userm/models"
	"time"

	"gorm.io/gorm"
)

func PreloadUsers() {
	var users = []models.User{
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 57, 29, 351794000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 57, 29, 351794000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Username: "darko123",
			// password: "pass"
			Password: "$2a$10$H6NBU8YNijGecCR6.iI1VucLMW76/2a40MAazaZrQpdEdg9YzCrUq",
			Email:    "darko@gmail.com",
			Role:     models.UserRole(2),
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Username: "marko123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Email:    "marko@gmail.com",
			Role:     models.UserRole(2),
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Username: "pavle123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Email:    "pavle@gmail.com",
			Role:     models.UserRole(1),
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			Username: "admin123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Email:    "admin@gmail.com",
			Role:     models.UserRole(3),
		},
	}

	for _, u := range users {
		DB.Create(&u)
	}
}
