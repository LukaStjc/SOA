package initializers

import (
	"context"
	"go-userm/graphdb"
	"go-userm/models"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

			ID:       1,
			Username: "darko123",
			// password: "pass"
			Password: "$2a$10$H6NBU8YNijGecCR6.iI1VucLMW76/2a40MAazaZrQpdEdg9YzCrUq",
			Email:    "darko@gmail.com",
			Role:     models.UserRole(1),
		},
		{
			Model: gorm.Model{
				CreatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				UpdatedAt: time.Date(2024, 3, 19, 16, 58, 7, 81778000, time.FixedZone("CET", 1*3600)),
				DeletedAt: gorm.DeletedAt{Valid: false},
			},
			ID:       2,
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
			ID:       3,
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
			ID:       4,
			Username: "admin123",
			// password: "pass"
			Password: "$2a$10$uDfTToyGNwFY6HZsbSnveeR7FW7UhRvk7IHP8KGsvvbHbRlXjZmca",
			Email:    "admin@gmail.com",
			Role:     models.UserRole(3),
		},
	}

	ctx := context.Background()
	session := Neo4JDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	neo4j_tx, err := session.BeginTransaction(ctx)
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		DB.Create(&u)

		graphUser := models.GraphDBUser{
			// ID: u.ID,
			ID:       int64(u.ID),
			Username: u.Username,
		}

		// log.Printf("username: %s", u.Username)
		// log.Println("Neo4j driver prvi ulaz:  ", Neo4JDriver)

		if err := graphdb.WriteUser(&graphUser, ctx, neo4j_tx); err != nil {
			log.Printf("Error creating user node for %s in Neo4j database: %v", u.Username, err)
		} else {
			neo4j_tx.Commit(ctx)
		}
	}
}
