package initializers

import (
	"go-jwt/models"
)

func PreloadUsers() {
	var users = []models.User{
		{
			ID:       1,
			Username: "darko123",
			// password: "pass"
			Password: "$2a$10$H6NBU8YNijGecCR6.iI1VucLMW76/2a40MAazaZrQpdEdg9YzCrUq",
			Role:     2,
		},
		{
			ID:       2,
			Username: "marko123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Role:     2,
		},
		{
			ID:       3,
			Username: "pavle123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Role:     3,
		},
		{
			ID:       4,
			Username: "admin123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Role:     3,
		},
	}

	for _, u := range users {
		DB.Create(&u)
	}
}
