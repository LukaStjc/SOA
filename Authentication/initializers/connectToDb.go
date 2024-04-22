package initializers

import (
	"fmt"
	configurations "go-jwt/startup"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func ConnectToDb() {
// 	var err error
// 	dsn := os.Getenv("DB")
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic("Failed to connect to db")
// 	}
// }

func ConnectToDb(config *configurations.Configurations) {
	var err error
	// dsn := os.Getenv("DB")

	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	connectionParams := fmt.Sprintf("user=postgres password=ftn dbname=SOA_auth host=%s port=%s sslmode=disable", config.AuthenticationDBHost, config.AuthenticationDBPort)
	DB, err = gorm.Open(postgres.Open(connectionParams), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}
}
