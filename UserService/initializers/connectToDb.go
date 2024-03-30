package initializers

import (
	"fmt"
	configurations "go-userm/startup"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb(config *configurations.Configurations) {
	var err error
	// dsn := os.Getenv("")
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	connectionParams := fmt.Sprintf("user=postgres password=ftn dbname=SOA host=%s port=%s sslmode=disable", config.UserDBHost, config.UserDBPort)
	DB, err = gorm.Open(postgres.Open(connectionParams), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}
}
